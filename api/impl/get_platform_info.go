package impl

import (
	"encoding/json"
	"fmt"
	"github.com/ashu0000008/crypto-market-cap/db"
	"sort"
)

func GetPlatformInfo(target string) string {
	db, err := db.Setup()
	if err != nil {
		fmt.Print(err.Error())
		return ""
	}
	defer db.Close()

	rows, err := db.Query("select symbol, marketcap from coincap where date1 = (select date1 from coincap order by date1 desc limit 1)")
	if err != nil {
		fmt.Print(err.Error())
		return "db query error"
	}

	var capMap map[string]TokenInfo
	capMap = make(map[string]TokenInfo)

	for rows.Next() {
		var symbolTemp string
		var cap float64
		err = rows.Scan(&symbolTemp, &cap)
		if err != nil {
			fmt.Print(err.Error())
			continue
		}

		//根据symbol拿到platform
		var platform string
		err = db.QueryRow("select platform from coininfo where symbol = ?", symbolTemp).Scan(&platform)
		if err != nil || platform == "" {
			continue
		}

		if target != platform {
			continue
		}

		info := capMap[symbolTemp]
		info.Symbol = symbolTemp
		info.Cap = cap
		capMap[symbolTemp] = info
	}

	//将结果排序
	var tokens []TokenInfo
	for _, v := range capMap {
		tokens = append(tokens, v)
	}
	sort.Sort(TokenSlice(tokens))

	data, err := json.Marshal(tokens)
	if err != nil {
		fmt.Println("json marshal failed,err:", err)
		return "marshal error"
	}

	fmt.Print("GetCryptoPlatformsSummaryImpl catch record ", len(tokens), "\r\n")
	return string(data)
}

type TokenInfo struct {
	Symbol string  `json:"symbol"`
	Cap    float64 `json:"cap"`
}

type TokenSlice []TokenInfo

func (a TokenSlice) Less(i, j int) bool {
	return a[j].Cap < a[i].Cap
}

func (a TokenSlice) Len() int {
	return len(a)
}

func (a TokenSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
