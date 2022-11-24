package persistence

import (
	"gorm.io/gorm"
	"log"
)

type mysqlPersistence struct {
	DB *gorm.DB
}

func (mysql *mysqlPersistence) DoSave(entity any) (int64, error) {
	db := mysql.DB.Create(entity)
	if db.Error != nil {
		log.Printf("insert data error %v", db.Error)
		return db.RowsAffected, db.Error
	}
	return db.RowsAffected, nil
}

func (mysql *mysqlPersistence) DoUpdate(entity any) (int64, error) {
	db := mysql.DB.Updates(entity)
	if db.Error != nil {
		log.Printf("udpate data error %v", db.Error)
		return db.RowsAffected, db.Error
	}
	return db.RowsAffected, nil
}

func (mysql *mysqlPersistence) DoDelete(entity any) (int64, error) {
	db := mysql.DB.Delete(entity)
	if db.Error != nil {
		log.Printf("delete data error %v", db.Error)
		return db.RowsAffected, db.Error
	}
	return db.RowsAffected, nil
}

func (mysql *mysqlPersistence) DoQuery(dest interface{}, query interface{}, args ...interface{}) error {
	db := mysql.DB.Where(query, args).First(dest)
	if db.Error != nil {
		log.Printf("do query data error %v", db.Error)
		return db.Error
	}
	return nil
}