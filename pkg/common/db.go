package common

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func ConnectDb() *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		EnvDBHost(), EnvDBPort(), EnvDBUser(), EnvDBPassword(), EnvDBName())

	for {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Successfully connected to db.")
			return db
		} else {
			time.Sleep(1 * time.Second)
			log.Println("Couldn't connect to db, retrying in a moment")
		}
	}
}
