package persistence

import (
	"gorm.io/gorm"
)

var SysDao IDbEntityUpdate

func InitSysDao(sysDb *gorm.DB) {
	SysDao = &mysqlPersistence{
		DB: sysDb,
	}
}
