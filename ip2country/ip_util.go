package ip2country

import (
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
)

var dbFile = "./ip2country/GeoLite2-City.mmdb"

func GetCountry(ipAddr string) string {
	db, err := geoip2.Open(dbFile)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer db.Close()
	if ipAddr == "" {
		return ""
	}

	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return ""
	}
	record, err := db.City(ip)
	if err != nil {
		return ""
	}
	result := record.Country.Names["en"] + "-" + record.City.Names["en"]
	return result
}
