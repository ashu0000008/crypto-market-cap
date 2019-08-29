package impl

import (
	"encoding/json"
	"fmt"
	"github.com/ashu0000008/crypto-market-cap/db"
	"sort"
	"strings"
)

func GetPercentHistory(symbol string) string {
	db, err := db.Setup()
	if err != nil {
		fmt.Print(err.Error())
		return ""
	}
	defer db.Close()

	var set map[string]bool
	set = make(map[string]bool)
	for loop := 0; ; loop++ {
		stmt, err := db.Prepare("select date1 from coincap LIMIT ?, ?")
		rows, err := stmt.Query(2000*loop, 2000)
		if err != nil {
			return "db query error:" + err.Error()
		}

		var count int
		count = 0
		for rows.Next() {
			count++
			var date string
			err = rows.Scan(&date)
			if err != nil {
				fmt.Print(err.Error())
				continue
			}
			set[date] = true
		}

		if 0 == count {
			break
		}
	}

	result := make(map[string]CoinPercent)
	for key := range set {

		rows, err := db.Query("select symbol, marketcap from coincap where date1 = ?", key)
		if err != nil {
			return "db query for date1 error:" + err.Error()
		}

		var total float64
		var target float64

		total = 0
		target = 0

		for rows.Next() {
			var symbolTemp string
			var cap float64
			err = rows.Scan(&symbolTemp, &cap)
			if err != nil {
				fmt.Print(err.Error())
				continue
			}

			total += cap
			if strings.EqualFold(symbol, symbolTemp) {
				target = cap
			}
		}

		info := result[key]
		info.Date = key
		info.Percent = target / total
		result[key] = info
	}

	//将结果排序
	var coinPercents []CoinPercent
	for _, v := range result {
		coinPercents = append(coinPercents, v)
	}
	sort.Sort(CoinPercentSlice(coinPercents))

	data, err := json.Marshal(coinPercents)
	if err != nil {
		fmt.Println("json marshal failed,err:", err)
		return "marshal error"
	}

	fmt.Print("GetCryptoPlatformsSummaryImpl catch record ", len(coinPercents), "\r\n")
	return string(data)
}

type CoinPercent struct {
	Date    string  `json:"date"`
	Percent float64 `json:"percent"`
}

type CoinPercentSlice []CoinPercent

func (a CoinPercentSlice) Less(i, j int) bool {
	return a[j].Date > a[i].Date
}

func (a CoinPercentSlice) Len() int {
	return len(a)
}

func (a CoinPercentSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
