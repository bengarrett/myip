package myipio

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// https://api.my-ip.io/ip.json
//
// Output:
// {
// 	"success": true,
// 	"ip": "100.100.0.0",
// 	"type": "IPv4"
// }

// Result of IPv4 query.
type Result struct {
	Success bool   `json:"success"`
	IP      string `json:"ip"`
	Type    string `json:"type"`
}

var (
	result       Result
	domain                     = "api.my-ip.io"
	ErrNoIP                    = errors.New("ip address is empty")
	ErrNoSuccess               = errors.New("ip address is unsuccessful")
	ErrNoIPv4                  = errors.New("ip address is not ipv4")
	ErrInvalid                 = errors.New("ip address is invalid")
	Timeout      time.Duration = 5
)

// IPv4 returns the Internet facing IP address of the free my-ip.io service.
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
	res, err := c.Get("https://" + path.Join(domain, "ip.json"))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		StatusErr := errors.New("unusual my-ip.io server response")
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
	if !r.Success {
		return false, ErrNoSuccess
	}
	if !strings.EqualFold(r.Type, "ipv4") {
		return false, ErrNoIPv4
	}
	if net.ParseIP(r.IP) == nil {
		return false, ErrInvalid
	}

	return true, nil
}
