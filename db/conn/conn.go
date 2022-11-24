package conn

import (
	"database/sql"
	"github.com/kim118000/db/model/sys"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var SysDb *gorm.DB

func InitDbConn(config *DbConfig) *gorm.DB {
	sqlDB, err := sql.Open("mysql", config.Dns)
	sqlDB.SetMaxOpenConns(config.MaxConnNumber)
	sqlDB.SetMaxIdleConns(config.IdleConnNumber)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{

	})
	if err != nil {
		panic(err)
	}
	SysDb = db

	return SysDb
}

func initTable(db *gorm.DB) {
	db.AutoMigrate(&sys.User{})

}
