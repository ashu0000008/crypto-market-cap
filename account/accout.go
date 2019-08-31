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

	user := orm.User{}
	db.Table("user").Where("device_id = ?", deviceIdString).First(&user)
	if "" != user.Device_id {
		return
	}

	user = orm.User{Device_id: deviceIdString}
	db.Table("user").Create(&user)
}
