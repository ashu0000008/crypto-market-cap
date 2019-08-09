package main

import (
	"github.com/ashu0000008/crypto-market-cap/db"
	"github.com/ashu0000008/crypto-market-cap/task"
	"time"
)

func main() {
	db, err := db.Setup()
	if err != nil {
	} else {
		defer db.Close()
	}

	//启动
	task.TaskCryptoCollectorStart()

	//主线程不退出
	for {
		time.Sleep(1 * time.Minute)
		print("sleepy...\r\n")
	}
}
