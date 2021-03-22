package myipcom

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const domain = "api.myip.com"

func Get() error {
	for i := 1; i <= 3; i++ {
		if err := get(); err != nil {
			fmt.Printf("%d. %s: %s\n", i, domain, err)
			continue
		}
		break
	}
	return nil
}

func get() error {
	res, err := http.Get("https://" + domain)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		StatusErr := errors.New("unusual myip.com server response")
		return fmt.Errorf("%s, %w", strings.ToLower(res.Status), StatusErr)
	}

	r, err := parse(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(r.Country)

	// ip, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("%s", string(ip))

	// parse json

	// c, err := geolite2.City(string(ip))
	// if err != nil {
	// 	fmt.Println(": City error:", err)
	// 	return nil
	// }
	// if c != "" {
	// 	fmt.Printf(" %s", c)
	// }
	fmt.Println()

	return nil
}

// {"ip":"118.209.50.85","country":"Australia","cc":"AU"}
type Result struct {
	IP      string `json:"ip,omitempty" validate:"required"`
	Country string `json:"country"`
	ISOCode string `json:"cc"`
}

var result Result

func parse(r io.Reader) (Result, error) {
	jsonParser := json.NewDecoder(r)
	if err := jsonParser.Decode(&result); err != nil {
		return Result{}, err
	}
	return result, nil
}
