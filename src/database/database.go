package database

import (
	"database/sql"
	"fmt"
)

type DatabaseConf struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

func Connect(dbConf DatabaseConf) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
