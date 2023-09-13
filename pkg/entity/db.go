package entity

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"sports-book.com/pkg/db_query"
)

var dbConn *gorm.DB

// ConnectDB opens a connection to the database, so that other functions may query against it
func ConnectDB() error {
	gormDb, err := gorm.Open(mysql.Open("root:password@tcp(127.0.0.1:3306)/sports-book?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		fmt.Println(err)
		return err
	}
	db_query.SetDefault(gormDb)
	dbConn = gormDb
	return nil
}

func getConn() *gorm.DB {
	if dbConn == nil {
		ConnectDB()
	}
	return dbConn
}
