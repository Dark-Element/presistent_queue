package initializers

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func database() *sql.DB {
	db, err := sql.Open("mysql", "root:getalife@tcp(192.168.239.129:3306)/queue") //todo: ad config file, viper?
	if err != nil{
		log.Panic(err)
	}
	db.SetMaxOpenConns(100)
	return db
}