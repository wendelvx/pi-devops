package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConectDb() *gorm.DB {
	dns := "host=localhost user=dev password=1234 dbname=dev port=5432 sslmode=disable TimeZone=America/Sao_Paulo"

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})

	if err != nil {
		panic("Problem connect to the database;")
	}

	//db.AutoMigrate()

	return db
}
