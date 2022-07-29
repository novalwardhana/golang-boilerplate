package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBMasterConnection() *gorm.DB {
	uri := "postgres://noval:noval@localhost:5432/golang_boilerplate_db_master"
	db := CreateConnection(uri)
	return db
}

func CreateConnection(uri string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(uri))
	if err != nil {
		fmt.Println("Check db connection")
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("Check db connection")
		return nil
	}

	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(2)

	return db
}
