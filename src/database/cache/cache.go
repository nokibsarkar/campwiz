package cache

import (
	"fmt"
	"nokib/campwiz/consts"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetCacheDB() (db *gorm.DB, close func()) {
	dsn := consts.Config.Database.Cache.DSN
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
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
func InitCacheDB() {
	db, close := GetCacheDB()
	defer close()
	db.AutoMigrate(&Session{})
	// db.AutoMigrate(&Evaluation{})
	// db.AutoMigrate(&Participant{})
	// db.AutoMigrate(&Submission{})
	fmt.Println((db))
}
