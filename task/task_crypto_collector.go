package task

import (
	"database/sql"
	"encoding/json"
	"github.com/ashu0000008/crypto-market-cap/db"
	"github.com/ashu0000008/crypto-market-cap/fetchers"
)

//加密币收集任务
func TaskCryptoCollectorStart() {
	StartTask(20, taskCryptoCollector)
}

func taskCryptoCollector() {
	response := fetchers.Fetch()
	str := []byte(response)
	dataContainer := ResponseData{}
	json.Unmarshal(str, &dataContainer)

	db, err := db.Setup()
	if err != nil {
		print(err)
		return
	}
	defer db.Close()

	//创建表
	//sql_table := "CREATE TABLE IF NOT EXISTS coininfo(cid INTEGER PRIMARY KEY AUTOINCREMENT,cname VARCHAR(64),symbol VARCHAR(64),slug VARCHAR(64),maxsupply VARCHAR(64),circulatingsupply VARCHAR(64),platform VARCHAR(64),cmcrank VARCHAR(64),tags VARCHAR(64));"
	//_, err = db.Exec(sql_table)
	//if err != nil{
	//	print(err)
	//	return
	//}

	for _, coin := range dataContainer.Data {
		collectCoin(db, coin)
	}
}

func collectCoin(db *sql.DB, coin CoinInfo) {
	//币已存在数据库，就无需插入了
	target := coin.Symbol
	var name string
	err := db.QueryRow("select cname from coininfo where symbol = ?", target).Scan(&name)
	if name != "" {
		return
	}

	stmt, err := db.Prepare("INSERT INTO coininfo(cname, symbol, slug, maxsupply, circulatingsupply, platform, cmcrank) values(?,?,?,?,?,?,?)")
	if err != nil {
		print(err)
		return
	}

	_, err = stmt.Exec(coin.Name, coin.Symbol, coin.Slug, coin.MaxSupply, coin.CirculatingSupply, coin.Platform, coin.CmcRank)
	if err != nil {
		print(err)
		return
	}

}

type CoinInfo struct {
	Id                interface{} `json:"id"`
	Name              interface{} `json:"name"`
	Symbol            interface{} `json:"symbol"`
	Slug              interface{} `json:"slug"`
	MaxSupply         interface{} `json:"max_supply"`
	CirculatingSupply interface{} `json:"circulating_supply"`
	Platform          interface{} `json:"platform"`
	CmcRank           interface{} `json:"cmc_rank"`
	Tags              interface{} `json:"tags"`
}

type ResponseData struct {
	Data []CoinInfo `json:"data"`
}
