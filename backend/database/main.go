package database

import (
	"fmt"
	"nokib/campwiz/consts"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type CommonFilter struct {
	Limit         int    `form:"limit"`
	ContinueToken string `form:"next"`
	PreviousToken string `form:"prev"`
}

func GetDB() (db *gorm.DB, close func()) {
	dsn := consts.Config.Database.Main.DSN
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}
	return db, func() {
		raw_db, err := db.DB()
		if err != nil {
			panic("failed to connect database")
		}
		raw_db.Close()
	}
}
func GetDbWithoutDefaultTransaction() (db *gorm.DB, close func()) {
	dsn := consts.Config.Database.Main.DSN
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Warn),
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic("failed to connect database")
	}
	return db, func() {
		raw_db, err := db.DB()
		if err != nil {
			panic("failed to connect database")
		}
		raw_db.Close()
	}
}
func GetTestDB() (db *gorm.DB, close func()) {
	dsn := consts.Config.Database.Main.TestDSN
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		panic("failed to connect database")
	}
	return db, func() {
		raw_db, err := db.DB()
		if err != nil {
			panic("failed to connect database")
		}
		raw_db.Close()
	}
}
func InitDB(testing bool) {
	conn, close := GetDB()
	if testing {
		conn, close = GetTestDB()
		conn.Exec("DROP DATABASE IF EXISTS campwiz_test;")
		conn.Exec("CREATE DATABASE campwiz_test;")
		conn.Exec("USE campwiz_test;")
	}
	defer close()

	db := conn.Begin()
	// set character set to utf8mb4
	db.Exec("ALTER DATABASE campwiz CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci;")
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Campaign{})
	db.AutoMigrate(&Round{})
	db.AutoMigrate(&Task{})
	db.AutoMigrate(&Role{})
	db.AutoMigrate(&Submission{})
	db.AutoMigrate(&Evaluation{})
	fmt.Println((db))
	db.Commit()
}
