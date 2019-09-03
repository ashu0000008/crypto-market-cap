package data

import (
	"fmt"
	"github.com/leizongmin/huobiapi"
	"github.com/leizongmin/huobiapi/market"
	"strconv"
	"strings"
)

func HuobiInit(chanCollector chan string, chanData chan string) {
	// 创建客户端实例
	marketHuobi, err := huobiapi.NewMarket()
	if err != nil {
		panic(err)
	}
	// 订阅主题
	_ = marketHuobi.Subscribe("market.btcusdt.trade.detail", func(topic string, json *huobiapi.JSON) {
		dataReceive(chanData, topic, json)
	})

	// 请求数据
	json, err := marketHuobi.Request("market.btcusdt.detail")
	if err != nil {
		panic(err)
	} else {
		fmt.Println(json)
	}

	go getData(chanCollector, chanData, marketHuobi)

	// 进入阻塞等待，这样不会导致进程退出
	marketHuobi.Loop()
}

func getData(chanCollector chan string, chanData chan string, market *market.Market) {
	for {
		input := <-chanCollector
		topic := strings.ReplaceAll("market.btcusdt.trade.detail", "btcusdt", input+"btc")
		_ = market.Subscribe(topic, func(topic string, json *huobiapi.JSON) {
			dataReceive(chanData, topic, json)
		})
		_, _ = market.Request(topic)
	}
}

func dataReceive(chanData chan string, topic string, json *huobiapi.JSON) {
	market := strings.Split(topic, ".")[1]
	price, _ := json.Get("tick").Get("data").Get("price").Float64()
	fmt.Println(market, price)

	priceString := strconv.FormatFloat(price, 'E', 2, 64)
	chanData <- market + "-" + priceString
}
