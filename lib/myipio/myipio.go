package myipio

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"path"
	"strings"
)

// https://api.my-ip.io/ip.json

// {
// 	"success": true,
// 	"ip": "165.227.66.230",
// 	"type": "IPv4"
// }

// Result of
type Result struct {
	Success bool   `json:"success"`
	IP      string `json:"ip"`
	Type    string `json:"type"`
}

var result Result

const domain = "api.my-ip.io"

func IPv4() (s string) {
	var err error
	for i := 1; i <= 3; i++ {
		s, err = get()
		if err != nil {
			fmt.Printf("%d. %s: %s\n", i, domain, err)
			continue
		}
		break
	}

	return s
}

func get() (string, error) {
	resp, err := http.Get("https://" + path.Join(domain, "ip.json"))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		StatusErr := errors.New("unusual my-ip.io server response")
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
	jsonParser := json.NewDecoder(r)
	if err := jsonParser.Decode(&result); err != nil {
		return Result{}, err
	}

	return result, nil
}

func (r Result) valid() (bool, error) {
	if r.IP == "" {
		return false, errors.New("ip address is empty")
	}
	if !r.Success {
		return false, errors.New("ip address is unsuccessful")
	}
	if strings.ToLower(r.Type) != "ipv4" {
		return false, errors.New("ip address is not ipv4")
	}
	if net.ParseIP(r.IP) == nil {
		return false, errors.New("ip address is invalid")
	}

	return true, nil
}
