package myipcom

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// https://api.myip.com
//
// Output:
// {"ip":"118.209.50.85","country":"Australia","cc":"AU"}

// Result of IPv4 query.
type Result struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	ISOCode string `json:"cc"`
}

var (
	result     Result
	domain                   = "api.myip.com"
	ErrNoIP                  = errors.New("ip address is empty")
	ErrInvalid               = errors.New("ip address is invalid")
	Timeout    time.Duration = 5
)

// IPv4 returns the Internet facing IP address using the free myip.com service.
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
		StatusErr := errors.New("unusual myip.com server response")
		return "", fmt.Errorf("%s, %w", strings.ToLower(res.Status), StatusErr)
	}
	r, err := parse(res.Body)
	if err != nil {
		return "", err
	}
	if ok, err := r.valid(); !ok {
		return r.IP, err
	}

	return r.IP, nil
}

func parse(r io.Reader) (Result, error) {
	jsonParser := json.NewDecoder(r)
	if err := jsonParser.Decode(&result); err != nil {
		return Result{}, err
	}
	return result, nil
}

func (r Result) valid() (bool, error) {
	if r.IP == "" {
		return false, ErrNoIP
	}
	if net.ParseIP(r.IP) == nil {
		return false, ErrInvalid
	}

	return true, nil
}
