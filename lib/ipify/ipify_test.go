package ipify

import (
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
		wantErr string
	}{
		{"empty", "", "no Host in request URL"},
		{"html", "example.com", ""},
		{"404", path.Join(u.Path, "abcdef"), "404 not found"},
		{"okay", domain, ""},
	}
	for _, tt := range tests {
		d := tt.domain
		t.Run(tt.name, func(t *testing.T) {
			_, err := get(d)
			if err != nil && tt.wantErr != "" && !strings.Contains(fmt.Sprint(err), tt.wantErr) {
				t.Errorf("get() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
