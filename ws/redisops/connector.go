package redisops

import (
	"fmt"
	"github.com/go-redis/redis"
)

func RedisConnect() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "abc1234",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	return client
}
