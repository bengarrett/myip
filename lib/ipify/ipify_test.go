package ipify

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func BenchmarkRequest(b *testing.B) {
	ctx, timeout := context.WithTimeout(context.Background(), 5*time.Second)
	s, err := request(ctx, timeout, link)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(s)
}

func TestIPv4(t *testing.T) {
	rc, rto := context.WithTimeout(context.Background(), 5*time.Second)
	wantS, _ := request(rc, rto, link)

	ctx, timeout := context.WithTimeout(context.Background(), 5*time.Second)
	if gotS, _ := IPv4(ctx, timeout); gotS != wantS {
		t.Errorf("IPv4() = %v, want %v", gotS, wantS)
	}
}

func Test_get(t *testing.T) {
	tests := []struct {
		name    string
		domain  string
		wantErr string
	}{
		{"empty", "", "unsupported protocol scheme"},
		{"html", "https://example.com", ""},
		{"404", link + "/abcdef", "404 not found"},
		{"okay", link, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, timeout := context.WithTimeout(context.Background(), 5*time.Second)
			_, err := request(ctx, timeout, tt.domain)
			if err != nil && tt.wantErr != "" && !strings.Contains(fmt.Sprint(err), tt.wantErr) {
				t.Errorf("get() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
