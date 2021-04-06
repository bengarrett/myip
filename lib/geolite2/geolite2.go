package geolite2

import (
	// Embed GeoLite2 databases.
	_ "embed"
	"fmt"
	"net"

	"github.com/oschwald/maxminddb-golang"
)

const lang = "en"

//go:embed db/GeoLite2-Country/GeoLite2-Country.mmdb
var country []byte // nolint:gochecknoglobals

//go:embed db/GeoLite2-City/GeoLite2-City.mmdb
var city []byte // nolint:gochecknoglobals

func Country(ip string) (string, error) {
	// use FromBytes() instead of Open("file.mmdb")
	db, err := maxminddb.FromBytes(country)
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
	return record.Country.Names[lang], nil
}

func City(ip string) (string, error) {
	db, err := maxminddb.FromBytes(city)
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
