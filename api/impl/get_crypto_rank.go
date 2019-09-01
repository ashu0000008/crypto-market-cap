package impl

import (
	"encoding/json"
	"fmt"
	"github.com/ashu0000008/crypto-market-cap/db"
)

func GetCryptoRankImpl(page int, size int) string {
	db, err := db.Setup()
	if err != nil {
		print(err)
		return "db setup error"
	}
	defer db.Close()

	stmt, err := db.Prepare("select symbol, marketcap from coincap where date1 = (select date1 from coincap order by date1 desc limit 1) limit ?, ?")
	rows, err := stmt.Query(size*page, size)
	if err != nil {
		fmt.Print(err.Error())
		return "db qurey error"
	}

	var tokens []TokenInfo

	for rows.Next() {
		var symbolTemp string
		var cap float64
		err = rows.Scan(&symbolTemp, &cap)
		if err != nil {
			fmt.Print(err.Error())
			continue
		}

		v := TokenInfo{Symbol: symbolTemp, Cap: cap}
		tokens = append(tokens, v)
	}

	data, err := json.Marshal(tokens)
	if err != nil {
		fmt.Println("json marshal failed,err:", err)
		return "marshal error"
	}

	fmt.Print("GetCryptoListImpl catch record ", len(tokens), "\r\n")
	return string(data)
}
