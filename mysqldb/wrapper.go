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

//Takesup the configuration for mysql connection.
type MySqlConfig struct {
	User               string
	Passwd             string
	Addr               string
	DBName             string
	MaxOpenConnections int
	MaxIdleConnections int
}

//converts the configuration to the format understood by go sql driver.
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

// A pool maintains a set of connections.
// Bydefault no connection is created. The connection is created only when query is fired.
type MySqlPool struct {
	db *sql.DB
}

//Ping checks if we can still access the database.
func (this *MySqlPool) Ping() error {
	return this.db.Ping()
}

//GetConnection returns a fresh *MySqlConnection object which can be further used to perform queries.
func (this *MySqlPool) GetConnection() *MySqlConnection {
	connection := &MySqlConnection{db: this.db}
	return connection
}

//Close the DBpool
func (this *MySqlPool) Close() error {
	return this.db.Close()
}

//On a broader level this can be seen as a Mysql connection.
type MySqlConnection struct {
	db   *sql.DB
	tx   *sql.Tx
	stmt *sql.Stmt
}

//A dummy function which pretends to close the MySqlConnection
//but actually MySqlConnection is a virtual entity that does not make any connection, thus does not needs to be closed.
// Its actually the stmt and tx that needs to be closed. Closing of stmt and tx is internally handled by this wrapper.
// SO, its safe if the user does not call this close method.But for clarity purpose user should call this method.
func (this *MySqlConnection) Close() error {
	return nil
}

//Access to underlying tx object.
func (this *MySqlConnection) GetRawTx() *sql.Tx {
	return this.tx
}

//Access to underlying db object.
func (this *MySqlConnection) GetRawConnection() *sql.DB {
	return this.db
}

//Wrapper for Prepare() sql method.
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

//map custom errors to sql driver errors.
func (this *MySqlConnection) prepareError(err error) error {

	if err == sql.ErrNoRows {
		err = ErrNoRows
	}
	return err
}

// Fetches a particular row based on the query and criteria supplied
// You need to supply pointer to struct(*struct) as holder for row values.
//
// Usage:
//  conn := pool.GetConnection()
//  holder := User{}
//  conn.FetchRowByQuery(query, &holder, params)
func (this *MySqlConnection) FetchRowByQuery(query string, holder interface{}, args ...interface{}) error {
	return this.FetchRowsByQuery(query, holder, args...)
}

// Fetches one or more rows based on the supplied query.
// You need to supply pointer to slice of struct (*[]struct) as holder.
//
// Usage:
//  conn := pool.GetConnection()
//  holder := []User{}
//  conn.FetchRowByQuery(query, &holder, params)
func (this *MySqlConnection) FetchRowsByQuery(query string, holder interface{}, args ...interface{}) error {
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
	return mysqlRows.scan(holder)
}

//Transaction related functions below

//Start a Transaction.
func (this *MySqlConnection) StartTransaction() error {
	//Before starting a new transaction on this connection
	//First, close previous transaction if any open on this connection.
	this.RollBack()

	tx, err := this.db.Begin()
	if err != nil {
		return err
	}
	this.tx = tx
	return nil
}

//Commit the existing transaction, if any
//It automatically closes the Tx object, SO you not need to do it explicitely.
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

//It automatically closes the Tx object, SO you not need to do it explicitely.
func (this *MySqlConnection) RollBack() error {
	//First check if we are in a transaction
	//If so, rollback the transaction and reset every thing.
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

func (this *MySqlConnection) IsInTransaction() bool {
	if this.tx != nil {
		return true
	}
	return false
}

func (this *MySqlConnection) resetTx() {
	this.tx = nil
}

//Contains rows object returned from db
type MySqlRows struct {
	rows *sql.Rows
}

//Scans the data from sql.Rows into the holder provided
//
// Holder can either be:
// Pointer to struct (*struct)
//      or
// Pointer to slice of structs (*[]struct).
func (this *MySqlRows) scan(holder interface{}) error {

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
	columns, err := this.getColumns()
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

//Get the columns returned by query.
func (this *MySqlRows) getColumns() (columns []string, err error) {
	columns, err = this.rows.Columns()
	return
}
