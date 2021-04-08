package myipcom

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"
	"testing"
)

func BenchmarkGet(b *testing.B) {
	s, err := get(domain)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(s)
}

func TestIPv4(t *testing.T) {
	wantS, _ := get(domain)
	if gotS := IPv4(); gotS != wantS {
		t.Errorf("IPv4() = %v, want %v", gotS, wantS)
	}
}

func Test_get(t *testing.T) {
	u, err := url.Parse(domain)
	if err != nil {
		t.Errorf("failed to parse domain %q, %s", domain, err)
	}
	tests := []struct {
		name    string
		domain  string
		isValid bool
		wantErr string
	}{
		{"empty", "", false, "no Host in request URL"},
		{"html", "example.com", false, "invalid character"},
		{"404", path.Join(u.Path, "abcdef"), false, "404 not found"},
		{"okay", domain, true, ""},
	}
	for _, tt := range tests {
		d := tt.domain
		t.Run(tt.name, func(t *testing.T) {
			gotS, err := get(d)
			if err != nil && tt.wantErr != "" && !strings.Contains(fmt.Sprint(err), tt.wantErr) {
				t.Errorf("get() error = %v, want %v", err, tt.wantErr)
			}
			if bool(gotS != "") != tt.isValid {
				t.Errorf("get() = %v, want an ip addr: %v", gotS, tt.isValid)
			}
		})
	}
}

func TestResult_valid(t *testing.T) {
	const (
		addr = "1.1.1.1"
		c    = "Australia"
		iso  = "AU"
	)
	type fields struct {
		IP      string
		Country string
		ISOCode string
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr error
	}{
		{"no ip", fields{"", "", ""}, false, ErrNoIP},
		{"invalid", fields{"1.1", "", ""}, false, ErrInvalid},
		{"ipv6", fields{"2001:db8:8714:3a90::12", "", ""}, false, ErrNoIPv4},
		{"valid", fields{addr, c, iso}, true, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Result{
				IP:      tt.fields.IP,
				Country: tt.fields.Country,
				ISOCode: tt.fields.ISOCode,
			}
			got, err := r.valid()
			if got != tt.want {
				t.Errorf("Result.valid() = %v, want %v", got, tt.want)
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Result.valid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
