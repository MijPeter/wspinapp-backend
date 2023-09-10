package common

import (
	"example/wspinapp-backend/internal/common/schema"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func InitDbWithConfig(cfg *gorm.Config) *gorm.DB {
	db := ConnectDb(cfg)

	db.AutoMigrate(
		&schema.Wall{},
		&schema.Route{},
		&schema.Hold{},
	)

	return db
}

func ConnectDb(cfg *gorm.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		EnvDBHost(), EnvDBPort(), EnvDBUser(), EnvDBPassword(), EnvDBName())

	for {
		db, err := gorm.Open(postgres.Open(dsn), cfg)
		if err == nil {
			log.Println("Successfully connected to db.")
			return db
		} else {
			time.Sleep(1 * time.Second)
			log.Printf("Couldn't connect to db, retrying in a moment %s\n", err)
		}
	}
}
