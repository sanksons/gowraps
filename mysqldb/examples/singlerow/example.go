package main

import (
	"fmt"

	"github.com/sanksons/gowraps/mysqldb"
)

func main() {

	type User struct {
		Name       string
		Data       string
		Occupation *string
	}
	config := mysqldb.MySqlConfig{
		User:               "root",
		Passwd:             "123456",
		Addr:               "ubuntuvm:3306",
		DBName:             "tradeanalysis",
		MaxOpenConnections: 10,
		MaxIdleConnections: 2,
	}

	pool, err := mysqldb.Initiate(config)
	if err != nil {
		panic(err.Error())
	}
	defer pool.Close()

	conn := pool.GetConnection()
	defer conn.Close()

	query := "SELECT name,occ as occupation,data from user"
	user := User{}

	err = conn.FetchRowByQuery(query, &user)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("%v", user)
}
