package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLDB() (*sql.DB, error) {
	db, err := sql.Open("mysql","root:123456@tcp(localhost:3306)/nutrition")
	if err != nil {
		return nil, err
	}
	err=db.Ping()


	if err!=nil{
		return nil,err
	}

	return db, nil
}
