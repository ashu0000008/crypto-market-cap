package task

import (
	"database/sql"
	"encoding/json"
	"github.com/ashu0000008/crypto-market-cap/db"
	"github.com/ashu0000008/crypto-market-cap/fetchers"
	"strings"
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
		CollectQuote(db, coin)
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

	var platform = coin.Platform.Name
	_, err = stmt.Exec(coin.Name, coin.Symbol, coin.Slug, coin.MaxSupply, coin.CirculatingSupply, platform, coin.CmcRank)
	if err != nil {
		print(err)
		return
	}
}

func CollectQuote(db *sql.DB, coin CoinInfo) {
	var data = coin.Quote.USD
	if data.MarketCap == 0 {
		return
	}

	tmp := strings.Split(data.Date, "T")
	if tmp == nil || len(tmp) != 2 {
		return
	}

	var id int
	err := db.QueryRow("select id from coincap where symbol = ? and date1 = ?", coin.Symbol, tmp[0]).Scan(&id)
	if id != 0 {
		return
	}

	stmt, err := db.Prepare("INSERT INTO coincap(symbol, date1, price, volume, marketcap, percent1h, percent24h, percent7d) values(?,?,?,?,?,?,?,?)")
	if err != nil {
		print(err)
		return
	}

	_, err = stmt.Exec(coin.Symbol, tmp[0], data.Price, data.Volume, data.MarketCap, data.Percent1h, data.Percent24h, data.Percent7d)
	if err != nil {
		print(err)
		return
	}
}

type CoinInfo struct {
	Id                interface{}  `json:"id"`
	Name              interface{}  `json:"name"`
	Symbol            interface{}  `json:"symbol"`
	Slug              interface{}  `json:"slug"`
	MaxSupply         interface{}  `json:"max_supply"`
	CirculatingSupply interface{}  `json:"circulating_supply"`
	Platform          PlatformInfo `json:"platform"`
	CmcRank           interface{}  `json:"cmc_rank"`
	Tags              interface{}  `json:"tags"`
	Quote             QuoteInfo    `json:"quote"`
}

type ResponseData struct {
	Data []CoinInfo `json:"data"`
}

type PlatformInfo struct {
	Name string `json:"name"`
}

type QuoteInfo struct {
	USD QuoteUsd `json:"USD"`
}

type QuoteUsd struct {
	Price      float64 `json:"price"`
	Volume     float64 `json:"volume_24h"`
	Percent1h  float64 `json:"percent_change_1h"`
	Percent24h float64 `json:"percent_change_24h"`
	Percent7d  float64 `json:"percent_change_7d"`
	MarketCap  float64 `json:"market_cap"`
	Date       string  `json:"last_updated"`
}
