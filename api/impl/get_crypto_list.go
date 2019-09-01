package impl

import (
	"encoding/json"
	"fmt"
	"github.com/ashu0000008/crypto-market-cap/db"
)

func GetCryptoListImpl(page int, size int) string {
	db, err := db.Setup()
	if err != nil {
		print(err)
		return "db setup error"
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM coininfo")
	if err != nil {
		return "db query error:" + err.Error()
	}

	var coinList []CoinInfo
	for rows.Next() {
		var info CoinInfo
		err = rows.Scan(&info.Cid, &info.Cname, &info.Symbol, &info.Slug, &info.MaxSupply, &info.CirculatingSupply, &info.Platform, &info.Cmcrank)
		if err != nil {
			fmt.Print(err.Error())
			continue
		}
		coinList = append(coinList, info)
	}

	data, err := json.Marshal(coinList)
	if err != nil {
		fmt.Println("json marshal failed,err:", err)
		return "marshal error"
	}

	fmt.Print("GetCryptoListImpl catch record ", len(coinList), "\r\n")
	return string(data)
}

type CoinInfo struct {
	Cid               int         `json:"cid"`
	Cname             string      `json:"cname"`
	Symbol            string      `json:"symbol"`
	Slug              string      `json:"slug"`
	MaxSupply         interface{} `json:"maxsupply"`
	CirculatingSupply float64     `json:"circulatingsupply"`
	Platform          string      `json:"platform"`
	Cmcrank           int         `json:"cmcrank"`
}
