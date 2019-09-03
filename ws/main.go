package main

import (
	"github.com/ashu0000008/crypto-market-cap/ws/data"
	"github.com/ashu0000008/crypto-market-cap/ws/server"
)

func main() {
	chanCollector := make(chan string) //传递订阅关系
	chanData := make(chan string)

	data.StartDataCollect(chanCollector, chanData)
	server.StartWSServer(chanCollector, chanData)
}
