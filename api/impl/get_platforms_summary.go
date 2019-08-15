package impl

import (
	"encoding/json"
	"fmt"
	"github.com/ashu0000008/crypto-market-cap/db"
	"sort"
)

func GetCryptoPlatformsSummaryImpl() string {
	db, err := db.Setup()
	if err != nil {
		fmt.Print(err.Error())
		return "db setup error"
	}
	defer db.Close()

	rows, err := db.Query("select symbol, marketcap from coincap where date1 = (select date1 from coincap order by date1 desc limit 1)")
	if err != nil {
		fmt.Print(err.Error())
		return "db query error"
	}

	var capMap map[string]PlatformInfo
	capMap = make(map[string]PlatformInfo)

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

		info := capMap[platform]
		info.Symbol = platform
		info.Num += 1
		info.Cap += cap
		capMap[platform] = info
	}

	//将结果排序
	var platforms []PlatformInfo
	for _, v := range capMap {
		platforms = append(platforms, v)
	}
	sort.Sort(PlatformSlice(platforms))

	data, err := json.Marshal(platforms)
	if err != nil {
		fmt.Println("json marshal failed,err:", err)
		return "marshal error"
	}

	fmt.Print("GetCryptoPlatformsSummaryImpl catch record ", len(platforms), "\r\n")
	return string(data)
}

type PlatformInfo struct {
	Symbol string  `json:"symbol"`
	Cap    float64 `json:"cap"`
	Num    int     `json:"num"`
}

type PlatformSlice []PlatformInfo

func (a PlatformSlice) Less(i, j int) bool {
	return a[j].Cap < a[i].Cap
}

func (a PlatformSlice) Len() int {
	return len(a)
}

func (a PlatformSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
