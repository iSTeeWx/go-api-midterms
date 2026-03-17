package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

const (
	driver   = "mysql"
	user     = "admin"
	password = "admin"
	host     = "localhost"
	port     = 3306
	database = "examblanc"
)

var Instance *sql.DB

func Connect() *sql.DB {

	var info = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", user, password, host, port, database)

	conn, err := sql.Open(driver, info)

	if err != nil {
		panic(err)
	}

	err = conn.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfuly connected to database")


	return conn
}
