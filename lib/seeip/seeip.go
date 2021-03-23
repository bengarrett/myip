package seeip

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

// https://seeip.org

const domain = "ip4.seeip.org"

// IPv4 returns the Internet facing IP address using the free seeip.org service.
func IPv4() (s string) {
	var err error
	for i := 1; i <= 3; i++ {
		if s, err = get(); err != nil {
			fmt.Printf("%d. %s: %s\n", i, domain, err)
			continue
		}
		break
	}

	return s
}

func get() (string, error) {
	res, err := http.Get("https://" + domain)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		StatusErr := errors.New("unusual seeip.org server response")
		return "", fmt.Errorf("%s, %w", strings.ToLower(res.Status), StatusErr)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	ip := string(b)
	if ok, err := valid(ip); !ok {
		return ip, err
	}

	return ip, nil
}

func valid(ip string) (bool, error) {
	if ip == "" {
		return false, errors.New("ip address is empty")
	}
	if net.ParseIP(ip) == nil {
		return false, errors.New("ip address is invalid")
	}

	return true, nil
}
