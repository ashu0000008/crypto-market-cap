package account

import (
	"encoding/json"
	"fmt"
)

func ApiGetFavorite(device string) string {
	result := GetFavorite(device)
	data, err := json.Marshal(result)
	if err != nil {
		fmt.Println("json marshal failed,err:", err)
		return "marshal error"
	}
	return string(data)
}

func ApiAddFavorite(device string, symbol string) bool {
	AddFavorite(device, symbol)
	return true
}

func ApiDeleteFavorite(device string, symbol string) bool {
	RemoveFavorite(device, symbol)
	return true
}
