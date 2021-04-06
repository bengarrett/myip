package seeip

import (
	"fmt"
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
	tests := []struct {
		name    string
		domain  string
		isValid bool
		wantErr bool
	}{
		{"empty", "", false, true},
		{"html", "example.com", false, true},
		{"404", "ip4.seeip.org/abcdef", false, true},
		{"okay", "ip4.seeip.org", true, false},
	}
	for _, tt := range tests {
		d := tt.domain
		t.Run(tt.name, func(t *testing.T) {
			gotS, err := get(d)
			if bool(err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, want %v", err, tt.wantErr)
			}
			gotV, _ := valid(gotS)
			if gotV != tt.isValid {
				t.Errorf("get() = %v, want %v", gotS, tt.isValid)
			}
		})
	}
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
