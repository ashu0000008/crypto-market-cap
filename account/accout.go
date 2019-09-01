package account

import (
	"github.com/ashu0000008/crypto-market-cap/account/orm"
)

func Check2AddUser(deviceIdString string) {
	db, err := orm.Connect2db()
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	if "" == deviceIdString {
		deviceIdString = "anonymous"
	}

	user := orm.User{}
	db.Table("user").Where("device_id = ?", deviceIdString).First(&user)
	if "" != user.Device_id {
		return
	}

	user = orm.User{Device_id: deviceIdString}
	db.Table("user").Create(&user)
}

func AddFavorite(deviceId string, symbol string) {
	db, err := orm.Connect2db()
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	item := orm.Favorite{}
	db.Table("favorite").Where("device_id = ? && symbol = ?", deviceId, symbol).First(&item)
	if "" != item.Symbol {
		return
	}

	item = orm.Favorite{Device_id: deviceId, Symbol: symbol}
	db.Table("favorite").Create(&item)
}

func RemoveFavorite(deviceId string, symbol string) {
	db, err := orm.Connect2db()
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	item := orm.Favorite{}
	db.Table("favorite").Where("device_id = ? && symbol = ?", deviceId, symbol).First(&item)
	if "" == item.Symbol {
		return
	}

	item = orm.Favorite{Device_id: deviceId, Symbol: symbol}
	db.Table("favorite").Delete(&item)
}

func GetFavorite(deviceId string) []orm.Favorite {
	db, err := orm.Connect2db()
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	var items []orm.Favorite
	db.Table("favorite").Where("device_id = ?", deviceId).Find(&items)

	return items
}
