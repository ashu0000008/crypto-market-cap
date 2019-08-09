package task

import "github.com/ashu0000008/crypto-market-cap/fetchers"

//加密币收集任务
func TaskCryptoCollectorStart() {
	StartTask(20, taskCryptoCollector)
}

func taskCryptoCollector() {
	fetchers.Fetch()
}
