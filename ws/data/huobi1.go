package data

import (
	"fmt"
	"github.com/ashu0000008/crypto-market-cap/ws/config"
	"github.com/go-redis/redis"
	"github.com/leizongmin/huobiapi"
	"github.com/leizongmin/huobiapi/market"
	"strconv"
	"strings"
)

func HuobiInit1(chanCollector *redis.PubSub, redisClient *redis.Client) {
	// 创建客户端实例
	marketHuobi, err := huobiapi.NewMarket()
	if err != nil {
		panic(err)
	}
	// 订阅主题
	_ = marketHuobi.Subscribe("market.btcusdt.trade.detail", func(topic string, json *huobiapi.JSON) {
		dataReceive1(redisClient, topic, json)
	})

	// 请求数据
	json, err := marketHuobi.Request("market.btcusdt.detail")
	if err != nil {
		panic(err)
	} else {
		fmt.Println(json)
	}

	go getData1(chanCollector, redisClient, marketHuobi)

	// 进入阻塞等待，这样不会导致进程退出
	marketHuobi.Loop()
}

func getData1(chanCollector *redis.PubSub, redisClient *redis.Client, market *market.Market) {
	for {
		input, _ := chanCollector.ReceiveMessage()

		var target string
		if strings.EqualFold(input.Payload, "btc") {
			target = "btcusdt"
		} else {
			target = input.Payload + "btc"
		}

		topic := strings.ReplaceAll("market.btcusdt.trade.detail", "btcusdt", target)
		_ = market.Subscribe(topic, func(topic string, json *huobiapi.JSON) {
			dataReceive1(redisClient, topic, json)
		})
		_, _ = market.Request(topic)
	}
}

func dataReceive1(redisClient *redis.Client, topic string, json *huobiapi.JSON) {
	market := strings.Split(topic, ".")[1]

	data, err := json.Get("tick").Get("data").Array()
	if nil != err {
		return
	}

	for index := range data {
		price, _ := json.Get("tick").Get("data").GetIndex(index).Get("price").Float64()
		fmt.Println(market, price)

		var prec int
		if strings.Contains(topic, "usdt") {
			prec = 2
		} else {
			prec = 8
		}
		priceString := strconv.FormatFloat(price, 'f', prec, 64)
		//fmt.Println(priceString)
		data := market + "-" + priceString

		redisClient.Publish(config.REDIS_QUOTE_DATA_NAME, data)
	}

	//price, _ := json.Get("tick").Get("data").GetIndex(0).Get("price").Float64()
	//fmt.Println(market, price)
	//
	//var prec int
	//if strings.Contains(topic, "usdt") {
	//	prec = 2
	//} else {
	//	prec = 8
	//}
	//priceString := strconv.FormatFloat(price, 'f', prec, 64)
	////fmt.Println(priceString)
	//data := market + "-" + priceString
	//
	//redisClient.Publish(config.REDIS_QUOTE_DATA_NAME, data)
}
