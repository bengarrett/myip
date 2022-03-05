// Package geolite2 returns the city or country location
// of an IP address.
// The IP region data is from GeoLite2 created by MaxMind,
// available from https://www.maxmind.com.
// © Ben Garrett https://github.com/bengarrett/myip
package geolite2

import (

	// Embed GeoLite2 databases.
	_ "embed"
	"errors"
	"fmt"
	"net"

	"github.com/oschwald/maxminddb-golang"
)

var ErrInvalid = errors.New("ip address is an invalid textual representation")

const lang = "en"

//go:embed db/GeoLite2-Country/GeoLite2-Country.mmdb
var country []byte

//go:embed db/GeoLite2-City/GeoLite2-City.mmdb
var city []byte

func Country(ip string) (string, error) {
	// use FromBytes() instead of Open("file.mmdb")
	db, err := maxminddb.FromBytes(country)
	if err != nil {
		return "", err
	}
	defer db.Close()

	pip := net.ParseIP(ip)
	if pip == nil {
		return "", ErrInvalid
	}

	var record struct {
		Country struct {
			Names map[string]string `maxminddb:"names"`
		} `maxminddb:"country"`
	} // Or any appropriate struct

	err = db.Lookup(pip, &record)
	if err != nil {
		return "", err
	}
	return record.Country.Names[lang], nil
}

func City(ip string) (string, error) {
	db, err := maxminddb.FromBytes(city)
	if err != nil {
		return "", err
	}
	defer db.Close()

	pip := net.ParseIP(ip)
	if pip == nil {
		return "", ErrInvalid
	}

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
	ct, co := record.City.Names[lang], record.Country.Names[lang]
	switch {
	case ct != "" && co != "":
		return fmt.Sprintf("%s, %s", ct, co), nil
	case co != "":
		return co, nil
	case ct != "":
		return ct, nil
	default:
		return "", nil
	}
}
