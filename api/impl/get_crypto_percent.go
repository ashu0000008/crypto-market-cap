package impl

import (
	"fmt"
	"github.com/ashu0000008/crypto-market-cap/db"
	"strings"
)

func GetPercent(symbol string) float64 {
	db, err := db.Setup()
	if err != nil {
		fmt.Print(err.Error())
		return 0
	}
	defer db.Close()

	rows, err := db.Query("select symbol, marketcap from coincap where date1 = (select date1 from coincap order by date1 desc limit 1)")
	if err != nil {
		fmt.Print(err.Error())
		return 0
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

	return target / total
}
