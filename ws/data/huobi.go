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

		var target string
		if strings.EqualFold(input, "btc") {
			target = "btcusdt"
		} else {
			target = input + "btc"
		}

		topic := strings.ReplaceAll("market.btcusdt.trade.detail", "btcusdt", target)
		_ = market.Subscribe(topic, func(topic string, json *huobiapi.JSON) {
			dataReceive(chanData, topic, json)
		})
		_, _ = market.Request(topic)
	}
}

func dataReceive(chanData chan string, topic string, json *huobiapi.JSON) {
	market := strings.Split(topic, ".")[1]
	price, _ := json.Get("tick").Get("data").GetIndex(0).Get("price").Float64()
	//fmt.Println(market, price)

	var prec int
	if strings.Contains(topic, "usdt") {
		prec = 2
	} else {
		prec = 8
	}
	priceString := strconv.FormatFloat(price, 'f', prec, 64)
	//fmt.Println(priceString)
	chanData <- market + "-" + priceString
}
