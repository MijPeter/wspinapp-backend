package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var DB *sql.DB

const (
	host     = "db"
	dbPort   = 5432
	user     = "wspinapp"
	password = "sprayitwhileyoucanmyfriend"
	dbname   = "db"
)

func ConnectDb() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, dbPort, user, password, dbname)
	var err error

	for {
		DB, err = sql.Open("postgres", psqlInfo)

		if err != nil {
			panic(err)
		}

		err = DB.Ping()
		if err == nil {
			log.Println("Successfully connected to db.")
			return
		} else {
			time.Sleep(1 * time.Second)
			log.Println("Couldn't connect to db, retrying in a moment")
		}
	}
}

func DisconnectDb() {
	err := DB.Close()

	if err != nil {
		panic(err)
	}

}
