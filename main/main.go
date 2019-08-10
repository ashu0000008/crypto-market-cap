package main

import (
	"github.com/ashu0000008/crypto-market-cap/task"
	"time"
)

func main() {
	//启动
	task.TaskCryptoCollectorStart()

	//主线程不退出
	for {
		time.Sleep(1 * time.Hour)
		print("sleepy...\r\n")
	}
}
