package db

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"sports-book.com/pkg/gorm/query"
)

var dbConn *gorm.DB

// Connect opens a connection to the database, so that other functions may query against it
func Connect() (*gorm.DB, error) {
	url := fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(
		mysql.Open(url),
	)
	query.SetDefault(db)
	dbConn = db
	return db, err
}

func getConn() *gorm.DB {
	if dbConn == nil {
		Connect()
	}
	return dbConn
}
