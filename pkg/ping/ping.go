// Package ping contains string functions for the
// ipv4 and ipv6 packages.
// © Ben Garrett https://github.com/bengarrett/myip
package ping

import (
	"errors"
	"fmt"

	"github.com/bengarrett/myip/pkg/geolite2"
)

var ErrInvalid = errors.New("invalid ip address")

const (
	Zero  = "(0/4) " // Zero returns a pre-ping string.
	Zero1 = "(0/1) " // Zero1 returns a pre-ping string for the first flag.
)

// City prints the IP address with its geographic location
// with both a country and city.
func City(ip string) (string, error) {
	c, err := geolite2.City(ip)
	if errors.Is(err, geolite2.ErrInvalid) {
		return "", ErrInvalid
	} else if err != nil {
		return "", fmt.Errorf("geo error for %s: %w", ip, err)
	}
	if c == "" {
		// reserved IP addresses that have no geolocaations
		// for example 0.0.0.0, 127.0.0.1
		return ip, nil
	}
	return fmt.Sprintf("%s, %s", ip, c), nil
}

// Contains returns true if x is found in the string array.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// Sprint returns a formatted IP address for the One request.
func Sprint(ip string) string {
	if ip == "" {
		return ""
	}
	city, err := City(ip)
	if err != nil {
		return fmt.Sprintf("\r(1/1) %s: %s", city, err)
	}
	// (1/1) 93.184.216.34, Norwell, United States
	return fmt.Sprintf("\r(1/1) %s", city)
}

// Sprint returns a formatted IP address for the All requests.
// The completed value is displayed as the number of finished requests.
// Enabling raw returns the IP address without any city or country information.
func Sprints(ip string, completed int, raw bool) string {
	if ip == "" {
		return ""
	}
	if raw {
		return count(completed, ip)
	}
	s, err := City(ip)
	if err != nil {
		return fmt.Sprintf("%s, %s", count(completed, ip), err)
	}
	return count(completed, s)
}

// Count returns a formatted job count and IP address.
func count(completed int, s string) string {
	const total = 4
	// (1/4) 93.184.216.34, Norwell, United States
	return fmt.Sprintf("\r(%d/%d) %s", completed, total, s)
}
