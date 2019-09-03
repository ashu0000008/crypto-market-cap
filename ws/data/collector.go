package data

func StartDataCollect(chanCollector chan string, chanData chan string) {
	go HuobiInit(chanCollector, chanData)
}
