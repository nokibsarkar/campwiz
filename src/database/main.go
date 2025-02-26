package database

import (
	"fmt"
	"nokib/campwiz/consts"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type IDType string

func (i *IDType) String() string {
	return string(*i)
}
func (i *IDType) GormDataType() string {
	return "varchar(255)"
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
func InitDB() {
	conn, close := GetDB()
	defer close()
	db := conn.Begin()
	// set character set to utf8mb4
	db.Exec("ALTER DATABASE campwiz CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci;")
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Campaign{})
	db.AutoMigrate(&Round{})
	db.AutoMigrate(&Task{})
	db.AutoMigrate(&Jury{})
	db.AutoMigrate(&Submission{})
	db.AutoMigrate(&Evaluation{})
	fmt.Println((db))
	db.Commit()
}
