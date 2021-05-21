package orm

import (
	"bark-server/model"
	"github.com/mritd/logger"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
)

type DbSqlite struct {
	Path string
}
type DbMySQL struct {
	Host     string
	Port     int
	User     string
	Pass     string
	Database string
	Params   string
}

type Database struct {
	DbType string
	Sqlite DbSqlite
	Mysql  DbMySQL
}

func parseDb(c *cli.Context) Database {
	database := Database{
		DbType: c.String("db-type"),
		Sqlite: DbSqlite{
			Path: c.String("sqlite-path"),
		},
		Mysql: DbMySQL{
			Host:     c.String("mysql-host"),
			Port:     c.Int("mysql-port"),
			User:     c.String("mysql-user"),
			Pass:     c.String("mysql-pass"),
			Database: c.String("mysql-database"),
			Params:   c.String("mysql-params"),
		},
	}
	return database
}

var db *gorm.DB

func GormSetup(c *cli.Context) {
	database := parseDb(c)
	var dsn string
	var dialector gorm.Dialector
	if database.DbType == "sqlite" {
		dbSqlite := database.Sqlite
		dsn = dbSqlite.Path
		dialector = sqlite.Open(dsn)
	}
	if database.DbType == "mysql" {
		dbmysql := database.Mysql
		dsn = dbmysql.User + ":" + dbmysql.Pass +
			"@tcp(" + dbmysql.Host + ":" + dbmysql.Pass + ")" +
			"/" + dbmysql.Database
		if dbmysql.Params != "" {
			dsn = dsn + "?" + dbmysql.Params
		}
		dialector = mysql.Open(dsn)
	}

	if dialector != nil {
		_db, err := gorm.Open(dialector, &gorm.Config{
			Logger: gorm_logger.Default.LogMode(gorm_logger.Info),
		})
		if err != nil {
			logger.Fatalf("failed to open %s. dsn(%s): %v", database.DbType, dsn, err)
		} else {
			db = _db
		}
		// Migrate the schema
		db.AutoMigrate(
			&model.Device{},
			&model.Message{},
			&model.User{},
			&model.UserBind{},
		)
	} else {
		logger.Fatalf("unknown database provider. type(%s)", database.DbType)
	}
}
