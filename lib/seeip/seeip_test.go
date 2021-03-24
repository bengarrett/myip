package seeip

import (
	"fmt"
	"testing"
)

func BenchmarkGet(b *testing.B) {
	s, err := get()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(s)
}

func Test_valid(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		want    bool
		wantErr bool
	}{
		{"empty", "", false, true},
		{"string", "abcde", false, true},
		{"too small", "1.2.3", false, true},
		{"to long", "1.2.3.4.5", false, true},
		{"range", "0.255.255.256", false, true},
		{"ipv4", "5.255.5.88", true, false},
		{"ipv6", "2002:3742:0100::3742:0100", true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := valid(tt.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("valid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPv4(t *testing.T) {
	tests := []struct {
		name    string
		domain  string
		isValid bool
		wantErr string
	}{
		{"empty", "", false, "ip address is empty"},
		{"html", "example.com", false, "ip address is invalid"},
		{"404", "ip4.seeip.org/abcdef", false, "ip address is empty"},
		{"okay", "ip4.seeip.org", true, ""},
	}
	for _, tt := range tests {
		domain = tt.domain
		t.Run(tt.name, func(t *testing.T) {
			gotS := IPv4()
			gotV, err := valid(gotS)
			if err != nil && tt.wantErr != "" && fmt.Sprint(err) != tt.wantErr {
				t.Errorf("IPv4() error = %v, want %v", err, tt.wantErr)
			}
			if gotV != tt.isValid {
				t.Errorf("IPv4() = %v, want %v", gotS, tt.isValid)
			}
		})
	}
}
