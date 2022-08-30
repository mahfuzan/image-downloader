package config

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Connect() {
	d, err := sqlx.Open("mysql", "root:root@tcp(127.0.0.1:3306)/training?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println("Conection failed")
	} else {
		log.Println("Connection established")
	}

	err = d.Ping()
	if err != nil {
		log.Println(err)
	}
	db = d
}

func GetDb() *sqlx.DB {
	return db
}
