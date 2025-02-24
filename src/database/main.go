package database

import (
	"fmt"
	"nokib/campwiz/consts"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDB() (db *gorm.DB, close func()) {
	dsn := consts.Config.Database.Main.DSN
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
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
	db.AutoMigrate(&Campaign{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&CampaignRound{})
	db.AutoMigrate(&Batch{})
	db.AutoMigrate(&Evaluation{})
	db.AutoMigrate(&Jury{})
	db.AutoMigrate(&Participant{})
	db.AutoMigrate(&Submission{})
	fmt.Println((db))
	db.Commit()
}
