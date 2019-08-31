package orm

import (
	"github.com/jinzhu/gorm"
)

func Connect2db() (*gorm.DB, error) {

	db, err := gorm.Open("mysql", "root:abc1234@(localhost:3306)/user?charset=utf8")
	if err != nil {
		panic("连接数据库失败")
	}

	return db, nil
}
