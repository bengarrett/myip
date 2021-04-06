package ipify

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// https://api.ipify.org
//
// Output:
// 1.1.1.1

var (
	domain                   = "api.ipify.org"
	ErrNoIP                  = errors.New("ip address is empty")
	ErrInvalid               = errors.New("ip address is invalid")
	ErrStatus                = errors.New("unusual ipify.org server response")
	Timeout    time.Duration = 5
)

// IPv4 returns the Internet facing IP address using the free ipify.org service.
func IPv4() string {
	s, err := get()
	if err != nil {
		if _, ok := err.(*url.Error); ok {
			if strings.Contains(err.Error(), "context deadline exceeded") {
				fmt.Printf("\n%s: timeout\n", domain)
				return ""
			}
			fmt.Printf("\n%s: %s\n", domain, err)
		}
		return ""
	}

	return s
}

func get() (string, error) {
	c := &http.Client{
		Timeout: Timeout * time.Second,
	}
	res, err := c.Get("https://" + domain)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s, %w", strings.ToLower(res.Status), ErrStatus)
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
		return false, ErrNoIP
	}
	if net.ParseIP(ip) == nil {
		return false, ErrInvalid
	}

	return true, nil
}
