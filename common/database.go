package common

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

// Init Opening a database and save the reference to `Database` struct.
func Init() *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       os.Getenv("DB_NAME") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(db:" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_TABLE") + "?" + os.Getenv("DB_OPTION"), // data source name
		DefaultStringSize:         256,                                                                                                                                                     // default size for string fields
		DisableDatetimePrecision:  true,                                                                                                                                                    // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,                                                                                                                                                    // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                                                                                                                                                    // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                                                                                                                                                   // auto configure based on currently MySQL version
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		log.Fatal("DB not init", err)
	}
	DB = db
	return DB
}

// Using this function to get a connection, you can create your connection pool here.
func GetDB() *gorm.DB {
	return DB
}
