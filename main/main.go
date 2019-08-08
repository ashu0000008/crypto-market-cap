package main

import (
	"github.com/ashu0000008/crypto-market-cap/db"
)

func main() {
	db, err := db.Setup()
	if err != nil {
	} else {
		defer db.Close()
	}
	//fetchers.Fetch()
}
