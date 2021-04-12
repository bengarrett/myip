package ipify

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

func BenchmarkRequest(b *testing.B) {
	ctx, timeout := context.WithTimeout(context.Background(), 5*time.Second)
	s, err := request(ctx, timeout, linkv4)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(s))
}

func TestTimeout(t *testing.T) {
	ctx, timeout := context.WithTimeout(context.Background(), 0*time.Second)
	if _, err := IPv4(ctx, timeout); !errors.Is(err, nil) {
		t.Errorf("IPv4() = %v, want %v", err, nil)
	}
}

func TestCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	want := context.Canceled
	s, err := IPv4(ctx, cancel)
	if s != "" || err != nil {
		t.Errorf("IPv4() error = %v, want error string", err)
	}
	if !errors.Is(ctx.Err(), want) {
		t.Errorf("IPv4() context.error = %v, want %v", err, want)
	}
}

func TestError(t *testing.T) {
	ctx, timeout := context.WithTimeout(context.Background(), 30*time.Second)
	if _, err := Request(ctx, timeout, "invalid url"); errors.Is(err, nil) {
		t.Errorf("IPv4() = %v, want an error", err)
	}
}

func Test_request(t *testing.T) {
	tests := []struct {
		name    string
		domain  string
		wantErr string
	}{
		{"empty", "", "unsupported protocol scheme"},
		{"html", "https://example.com", ""},
		{"404", "https://api.ipify.org/abcdef", "404 not found"},
		{"ipv4", linkv4, ""},
		{"ipv6", linkv6, ""},
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

func Test_valid(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		want    bool
		wantErr error
	}{
		{"empty", "", false, ErrNoIP},
		{"invalid", "abc", false, ErrInvalid},
		{"ipv4", "1.1.1.1", true, nil},
		{"ipv6", "0:0:0:0:0:FFFF:0101:0101", true, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := valid(tt.ip)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("valid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
