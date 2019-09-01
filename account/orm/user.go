package orm

type User struct {
	Device_id string `json:"device_id"`
	Member_id string
	Id        int64 `gorm:"AUTO_INCREMENT"`
}
