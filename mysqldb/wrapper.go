package mysqldb

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-sql-driver/mysql"
	reflexer "github.com/sanksons/go-reflexer"
)

//This does not creates any connection. It just creates an empty pool based on the supplied config.
//The connection is opened when Prepare statement is called.
func Initiate(config MySqlConfig) (*MySqlPool, error) {
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(config.MaxIdleConnections)
	db.SetMaxOpenConns(config.MaxOpenConnections)
	return &MySqlPool{db: db}, nil
}

//Define custom errors
var ErrNoRows = sql.ErrNoRows
var ErrToBeImpl = fmt.Errorf("To be Implemented")

type MySqlConfig struct {
	User               string
	Passwd             string
	Addr               string
	DBName             string
	MaxOpenConnections int
	MaxIdleConnections int
}

func (this *MySqlConfig) FormatDSN() string {
	c := mysql.Config{
		User:   this.User,
		Passwd: this.Passwd,
		Net:    "tcp",
		Addr:   this.Addr,
		DBName: this.DBName,
	}
	return c.FormatDSN()
}

type MySqlPool struct {
	db *sql.DB
}

//Ping checks if we can still access the database.
func (this *MySqlPool) Ping() error {
	return this.db.Ping()
}

//GetConnection returns a fresh MySqlConnection object which can be further used
// to perform queries.
func (this *MySqlPool) GetConnection() *MySqlConnection {
	connection := &MySqlConnection{db: this.db}
	return connection
}

func (this *MySqlPool) GetRawConnection() *sql.DB {
	return this.db
}

//Close the DBpool
func (this *MySqlPool) Close() error {
	return this.db.Close()
}

type MySqlConnection struct {
	db   *sql.DB
	tx   *sql.Tx
	stmt *sql.Stmt
}

//A dummy function whihc pretends to close the MySqlConnection
//but actually MySqlConnection is a virtual entity that does not make any connection
//thus does not needs to be closed.
// Its actually the stmt and tx that needs to be closed.
// Closing of stmt and tx is internally handled by this wrapper.
// SO, its safe if the user does not call this close method.
func (this *MySqlConnection) Close() error {
	return nil
}

func (this *MySqlConnection) PrepareStatement(query string) error {
	var stmt *sql.Stmt
	var err error
	this.stmt = nil
	if this.IsInTransaction() {
		stmt, err = this.tx.Prepare(query)
	} else {
		stmt, err = this.db.Prepare(query)
	}
	if err != nil {
		return err
	}
	this.stmt = stmt
	return nil
}

func (this *MySqlConnection) prepareError(err error) error {

	if err == sql.ErrNoRows {
		err = ErrNoRows
	}
	return err
}

// Fetches a particular row based on the criteria supplied
// Usage:
//  conn := pool.GetConnection()
//  conn.PrepareStatement(Query)
//  conn.FetchRow(value holder, params)
func (this *MySqlConnection) FetchRow(holder interface{}, args ...interface{}) error {
	defer this.stmt.Close()
	err := this.stmt.QueryRow(args...).Scan(holder)
	return this.prepareError(err)
}

// Fetches a particular row based on the query and criteria supplied
// Usage:
//  conn := pool.GetConnection()
//  conn.FetchRowByQuery(query, valueholder, params)
func (this *MySqlConnection) FetchRowByQuery(query string, holder interface{}, args ...interface{}) error {
	return this.FetchRowsByQuery(query, holder, args...)
}

func (this *MySqlConnection) FetchRowsByQuery(query string, holder interface{}, args ...interface{}) error {
	fmt.Println("fetch rows by query")
	fmt.Printf("%v", args)
	var rows *sql.Rows
	var err error
	if this.IsInTransaction() {
		rows, err = this.tx.Query(query, args...)
	} else {
		rows, err = this.db.Query(query, args...)
	}
	if err != nil {
		return err
	}

	mysqlRows := MySqlRows{rows: rows}
	fmt.Println("Going for scan")
	return mysqlRows.Scan(holder)
}

type MySqlRows struct {
	rows *sql.Rows
}

func (this *MySqlRows) Scan(holder interface{}) error {

	defer this.rows.Close()
	//check if holder is a pointer to struct i.e *struct, if not
	//check if holder is a pointer to slice of structs i.e *[]structs, if not
	//Err: Not a valid type supplied
	reflectObj := reflexer.ReflectObj{}
	reflectObj.Initiate(holder)
	if !reflectObj.CheckIfPtr() { //since we expect a pointer here, check for it.
		return fmt.Errorf("Expected a pointer but supplied, [%v]", reflectObj.Kind)
	}
	if !reflectObj.HasChild() {
		return fmt.Errorf("The supplied pointer points to blackhole")
	}
	child := reflectObj.GetChild()
	var structInfo map[string]int
	var err error

	var childStruct *reflexer.ReflectObj
	var isMulti bool
	if child.CheckIfSlice() {
		//Its probably a slice of structs. Drill down to get to struct.
		isMulti = true
		if !child.HasChild() {
			return fmt.Errorf("Expected slice of structs but didn't got it.")
		}
		childStruct = child.GetChild()
	} else if child.CheckIfStruct() {
		//Its  a struct itself.
		childStruct = child
		isMulti = false
	} else {
		return fmt.Errorf("Its neither a struct nor slice of structs")
	}

	//Get Column info
	columns, err := this.GetColumns()
	if err != nil {
		return fmt.Errorf("Could not get columns Info: %s", err.Error())
	}
	//Get info about struct
	structInfo, err = reflexer.GetInfoAboutFieldsofStruct(*childStruct)
	if err != nil {
		return fmt.Errorf("Scan Failed: %s", err.Error())
	}
	var iteration int
	var structList []reflect.Value
	for this.rows.Next() {
		//break out of loop incase we only need to fetch single row.
		iteration++
		if !isMulti && iteration > 1 {
			break
		}
		var rowStruct reflect.Value
		if isMulti {
			rowStruct = reflect.New(childStruct.T).Elem()

		} else {
			rowStruct = childStruct.V
		}
		var final []interface{}
		for _, col := range columns {
			col = strings.ToLower(col)
			index, ok := structInfo[col]
			if !ok {
				var skipVal string = ""
				pointerSkipval := &skipVal
				final = append(final, &pointerSkipval)
				continue //skip columns not found in struct
			}
			final = append(final, rowStruct.FieldByIndex([]int{index}).Addr().Interface())
		}
		err = this.rows.Scan(final...)
		if err != nil {
			return err
		}
		if isMulti {
			structList = append(structList, rowStruct)
		}
	}
	if isMulti {
		//!!IMPORTANT set the data in slice.
		tmp := reflect.Append(child.V, structList...)
		child.V.Set(tmp)
	}
	return nil
}

func (this *MySqlRows) GetColumns() (columns []string, err error) {
	columns, err = this.rows.Columns()
	return
}

//Transaction related functions below

func (this *MySqlConnection) IsInTransaction() bool {
	if this.tx != nil {
		return true
	}
	return false
}

func (this *MySqlConnection) resetTx() {
	this.tx = nil
}

func (this *MySqlConnection) StartTransaction() error {
	//First, close previous transaction if any open on this connection.
	this.RollBack()

	tx, err := this.db.Begin()
	if err != nil {
		return err
	}
	this.tx = tx
	return nil
}

//NOTE: commit and rollback automatically closes tx. SO no need to close it explicitely.
func (this *MySqlConnection) Commit() error {
	if !this.IsInTransaction() {
		return nil
	}
	defer this.resetTx()
	err := this.tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (this *MySqlConnection) RollBack() error {
	if !this.IsInTransaction() {
		return nil
	}
	defer this.resetTx()
	err := this.tx.Rollback()
	if err != nil {
		return err
	}
	return nil
}
