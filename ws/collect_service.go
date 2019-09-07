package main

import (
	"github.com/ashu0000008/crypto-market-cap/ws/config"
	"github.com/ashu0000008/crypto-market-cap/ws/data"
	"github.com/ashu0000008/crypto-market-cap/ws/redisops"
	"time"
)

func main() {
	redisClient := redisops.RedisConnect()
	chanCollector := redisClient.Subscribe(config.REDIS_QUOTE_MANAGER_NAME)

	data.StartDataCollect1(chanCollector, redisClient)

	for {
		time.Sleep(time.Minute * 1)
	}
}
