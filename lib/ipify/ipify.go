package ipify

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bengarrett/myip/lib/geolite2"
)

const domain = "api.ipify.org"

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
		StatusErr := errors.New("unusual ipify.org server response")
		return fmt.Errorf("%s, %w", strings.ToLower(res.Status), StatusErr)
	}

	ip, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Printf("%s", string(ip))
	c, err := geolite2.City(string(ip) + "x")
	if err != nil {
		fmt.Println(": City error:", err)
		return nil
	}
	if c != "" {
		fmt.Printf(" %s", c)
	}
	fmt.Println()

	return nil
}
