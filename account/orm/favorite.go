package orm

type Favorite struct {
	Id        int64  `gorm:"AUTO_INCREMENT"`
	Device_id string `json:"device_id"`
	Symbol    string `json:"symbol"`
}
