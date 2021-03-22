package geolite2

import (
	"fmt"
	"net"

	"github.com/oschwald/maxminddb-golang"
)

func Country(ip string) (string, error) {
	db, err := maxminddb.Open("assets/GeoLite2-Country/GeoLite2-Country.mmdb")
	if err != nil {
		return "", err
	}
	defer db.Close()

	pip := net.ParseIP(ip)

	var record struct {
		Country struct {
			Names map[string]string `maxminddb:"names"`
		} `maxminddb:"country"`
	} // Or any appropriate struct

	err = db.Lookup(pip, &record)
	if err != nil {
		return "", err
	}
	return record.Country.Names["en"], nil
}

func City(ip string) (string, error) {
	db, err := maxminddb.Open("assets/GeoLite2-City/GeoLite2-City.mmdb")
	if err != nil {
		return "", err
	}
	defer db.Close()

	pip := net.ParseIP(ip)

	var record struct {
		Country struct {
			Names map[string]string `maxminddb:"names"`
		} `maxminddb:"country"`
		City struct {
			Names map[string]string `maxminddb:"names"`
		} `maxminddb:"city"`
	} // Or any appropriate struct

	err = db.Lookup(pip, &record)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s, %s", record.City.Names["en"], record.Country.Names["en"]), nil
}
