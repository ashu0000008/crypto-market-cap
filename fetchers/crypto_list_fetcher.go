package fetchers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func Fetch() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest", nil)
	if err != nil {
		log.Print(err)
		return ""
	}

	q := url.Values{}
	q.Add("start", "1")
	q.Add("limit", "1000")
	q.Add("convert", "USD")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", "37b2d469-562a-4556-8886-b6e79658bc36")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		return ""
	}
	fmt.Println(resp.Status)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBody))
	return string(respBody)
}
