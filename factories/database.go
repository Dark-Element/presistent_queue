package factories

import (
	"database/sql"
	"log"
	"strconv"
)

func DbConn(host string, user string, password string, port int, dbName string, dbPool int) *sql.DB {
	db, err := sql.Open("mysql", generateDSN(host, user, password, port, dbName))
	if err != nil {
		log.Panic(err)
	}
	db.SetMaxOpenConns(dbPool)
	return db
}

func generateDSN(host string, user string, password string, port int, dbName string) string {
	return user + ":"+ password +"@tcp(" + host + ":"+ strconv.Itoa(port)+")/"+ dbName
}
