package myipcom

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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
	domain     = "api.myip.com"
	ErrNoIP    = errors.New("ip address is empty")
	ErrInvalid = errors.New("ip address is invalid")
)

const timeout time.Duration = 5

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
			return ""
		}
		log.Fatalln(err)
	}

	return s
}

func get() (string, error) {
	c := &http.Client{
		Timeout: timeout * time.Second,
	}
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://"+domain, nil)
	if err != nil {
		return "", err
	}
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		StatusErr := errors.New("unusual myip.com server response")
		return "", fmt.Errorf("%s, %w", strings.ToLower(resp.Status), StatusErr)
	}
	r, err := parse(resp.Body)
	if err != nil {
		return "", err
	}
	if ok, err := r.valid(); !ok {
		return r.IP, err
	}

	return r.IP, nil
}

func parse(r io.Reader) (Result, error) {
	var result Result
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
