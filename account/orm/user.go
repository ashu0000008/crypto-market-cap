package orm

type User struct {
	Device_id string
	Member_id string
	Id        int64 `gorm:"AUTO_INCREMENT"`
}
