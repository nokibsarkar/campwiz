package cache

import (
	"log"
	"nokib/campwiz/consts"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetCacheDB() (db *gorm.DB, close func()) {
	dsn := consts.Config.Database.Cache.DSN
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatal("failed to connect cache database")
	}
	return db, func() {
		raw_db, err := db.DB()
		if err != nil {
			log.Fatal("failed to get cache database on close")
		}
		raw_db.Close()
	}
}
func GetTestCacheDB() (db *gorm.DB, close func()) {
	dsn := consts.Config.Database.Cache.TestDSN
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatal("failed to connect cache database")
	}
	return db, func() {
		raw_db, err := db.DB()
		if err != nil {
			log.Fatal("failed to get cache database on close")
		}
		raw_db.Close()
	}
}
func InitCacheDB(testing bool) {
	db, close := GetCacheDB()
	if testing {
		db, close = GetTestCacheDB()
	}
	defer close()
	db.AutoMigrate(&Session{})
}
