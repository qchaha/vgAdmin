package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func dbConnect() *sql.DB {
	dbUser := "root"
	dbPassword := "root123"
	dbHost := "127.0.0.1"
	dbPort := "3306"
	dbName := "fordealOps"
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", dbUser, dbPassword, dbHost, dbPort, dbName))
	if err != nil {
		fmt.Println("open mysql error...")
		fmt.Println(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("msyql is not alive...")
		fmt.Println(err)
	}
	return db
}

func dbQuery(sql string) *sql.Rows {
	db := dbConnect()
	defer db.Close()
	row, err := db.Query(sql)

	if err != nil {
		fmt.Println("query database error...")
		fmt.Println(err)
	}
	return row
}

func dbModify(sql string, args ...interface{}) bool {
	var rc bool
	db := dbConnect()
	defer db.Close()
	ins, err := db.Prepare(sql)
	if err != nil {
		fmt.Println(sql)
		fmt.Println("db prepare error...")
		return false
	}
	defer ins.Close()
	r, err := ins.Exec(args...)
	if err != nil {
		fmt.Println(sql)
		fmt.Println("exec query error...")
		return false
	}
	affectRows, err := r.RowsAffected()
	if err != nil {
		fmt.Println(sql)
		fmt.Println("return affected rows error...")
		return false
	}
	if affectRows > 0 {
		rc = true
	} else {
		rc = false
	}
	return rc
}
