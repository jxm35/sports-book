package db

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

// Connect opens a connection to the database, so that other functions may query against it
func Connect() (*gorm.DB, error) {
	return gorm.Open(
		mysql.Open(
			fmt.Sprintf(
				"%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				os.Getenv("DB_USER"),
				os.Getenv("DB_PASSWORD"),
				os.Getenv("DB_HOST"),
				os.Getenv("DB_NAME"),
			),
		),
	)
}

func getConn() *gorm.DB {
	if dbConn == nil {
		Connect()
	}
	return dbConn
}
