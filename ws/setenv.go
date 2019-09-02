package main

import (
	"fmt"
	"os"
)

func main() {
	err := os.Setenv("BINANCE_SECRET", "lsmH3Xq92RjdJvDTWpfmIUBBAV9xTou5u1LEqXvvtAK1EsLd3XoBIEqeXv1WB1Kv")
	if err != nil {
		fmt.Println(err.Error())
	}

	err = os.Setenv("BINANCE_APIKEY", "cuI6fvo2HxKTtl4o5aH78bc3oojea2H7VEmmVnpkiMWSItb6E6oxTowMRvxqDo0q")
	if err != nil {
		fmt.Println(err.Error())
	}

	myenv := os.Getenv("BINANCE_SECRET")
	fmt.Println(myenv)
}
