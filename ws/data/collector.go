package data

import "github.com/go-redis/redis"

func StartDataCollect(chanCollector chan string, chanData chan string) {
	go HuobiInit(chanCollector, chanData)
}

func StartDataCollect1(chanCollector *redis.PubSub, redisClient *redis.Client) {
	go HuobiInit1(chanCollector, redisClient)
}
