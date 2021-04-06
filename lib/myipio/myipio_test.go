package myipio

import (
	"fmt"
	"strings"
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

func TestIPv4(t *testing.T) {
	wantS, _ := get()
	if gotS := IPv4(); gotS != wantS {
		t.Errorf("IPv4() = %v, want %v", gotS, wantS)
	}
}

func Test_get(t *testing.T) {
	tests := []struct {
		name    string
		domain  string
		isValid bool
		wantErr string
	}{
		{"empty", "", false, "no such host"},
		{"html", "example.com", false, "404 not found"},
		{"404", "api.my-ip.io/abcdef", false, "404 not found"},
		{"okay", "api.my-ip.io", true, ""},
	}
	for _, tt := range tests {
		domain = tt.domain
		t.Run(tt.name, func(t *testing.T) {
			gotS, err := get()
			if err != nil && tt.wantErr != "" && !strings.Contains(fmt.Sprint(err), tt.wantErr) {
				t.Errorf("IPv4() error = %v, want %v", err, tt.wantErr)
			}
			if bool(gotS != "") != tt.isValid {
				t.Errorf("IPv4() = %v, want an ip addr: %v", gotS, tt.isValid)
			}
		})
	}
}

func TestResult_valid(t *testing.T) {
	const addr = "1.1.1.1"
	const ipv4 = "IPv4"
	type fields struct {
		Success bool
		IP      string
		Type    string
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr error
	}{
		{"no ip", fields{false, "", ipv4}, false, ErrNoIP},
		{"fail", fields{false, addr, ipv4}, false, ErrNoSuccess},
		{"not ipv4", fields{true, addr, "IPv6"}, false, ErrNoIPv4},
		{"invalid", fields{true, "1", ipv4}, false, ErrInvalid},
		{"valid", fields{true, addr, ipv4}, true, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Result{
				Success: tt.fields.Success,
				IP:      tt.fields.IP,
				Type:    tt.fields.Type,
			}
			got, err := r.valid()
			if err != tt.wantErr {
				t.Errorf("Result.valid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Result.valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
