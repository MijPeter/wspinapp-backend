package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var DB *sql.DB

const (
	host     = "0.0.0.0"
	dbPort   = 5432
	user     = "wspinapp"
	password = "sprayitwhileyoucanmyfriend"
	dbname   = "wspinapp"
)

func ConnectDb() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, dbPort, user, password, dbname)
	var err error
	DB, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		DB.Close()
		panic(err)
	}
}

func DisconnectDb() {
	err := DB.Close()

	if err != nil {
		panic(err)
	}

}
