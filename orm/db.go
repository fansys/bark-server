package orm

import (
	"github.com/mritd/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
	"os"
	"path/filepath"
	"sync"
)

var dbOnce sync.Once
var db *gorm.DB

func New(dataDir string) {
	dbOnce.Do(func() {
		dbName := filepath.Join(dataDir, "bark.db")
		logger.Infof("init database [%s]...", dataDir)
		if _, err := os.Stat(dataDir); os.IsNotExist(err) {
			if err = os.MkdirAll(dataDir, 0755); err != nil {
				logger.Fatalf("failed to create database storage dir(%s): %v", dataDir, err)
			}
		} else if err != nil {
			logger.Fatalf("failed to open database storage dir(%s): %v", dataDir, err)
		}
		_db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
			Logger: gorm_logger.Default.LogMode(gorm_logger.Info),
		})
		if err != nil {
			logger.Fatalf("failed to create database file(%s): %v", dbName, err)
		}
		db = _db

		// Migrate the schema
		db.AutoMigrate(&Device{}, &User{}, &PushMessage{})
	})
}
