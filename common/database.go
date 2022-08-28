package common

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

// Opening a database and save the reference to `Database` struct.
func Init() *gorm.DB {

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "testuser:testpassword@tcp(db:3306)/testing?parseTime=true&loc=Local&charset=utf8mb4", // data source name
		DefaultStringSize:         256,                                                                                   // default size for string fields
		DisableDatetimePrecision:  true,                                                                                  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,                                                                                  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                                                                                  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                                                                                 // auto configure based on currently MySQL version
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	db.Exec("SET NAMES 'utf8mb4'; SET CHARACTER SET utf8mb4;")
	db.Exec("ALTER TABLE transactions MODIFY operation TEXT CHARACTER SET utf8mb4;")
	fmt.Println("db err: (Init) ", db.Error)
	if err != nil {
		fmt.Println("db err: (Init) ", err)
	}
	//db.LogMode(true)
	DB = db
	return DB
}

// Using this function to get a connection, you can create your connection pool here.
func GetDB() *gorm.DB {
	return DB
}
