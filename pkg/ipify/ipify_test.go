package ipify

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
)

func BenchmarkRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx, timeout := context.WithTimeout(context.Background(), 5*time.Second)
		s, err := request(ctx, timeout, linkv4)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(s))
	}
}

// ExampleIPv4 demonstrates an IPv4 address request with a 5 second timeout.
func ExampleIPv4() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s, err := IPv4(ctx, cancel)
	if err != nil {
		log.Printf("\n%s\n", err)
	}
	fmt.Println(s)
}

// ExampleIPv6 demonstrates cocurrent IPv4 and IPv6 address requests with a 5 second timeout.
func ExampleIPv6() {
	ctx4, cancel4 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel4()

	ctx6, cancel6 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel6()

	s4 := make(chan string)
	s6 := make(chan string)

	go func() {
		s, err := IPv4(ctx4, cancel4)
		if err != nil {
			log.Printf("\n%s\n", err)
		}
		s4 <- s
	}()

	go func() {
		s, err := IPv6(ctx6, cancel6)
		if err != nil {
			log.Printf("\n%s\n", err)
		}
		s6 <- s
	}()

	ip4 := <-s4
	ip6 := <-s6
	fmt.Println(ip4)
	fmt.Println(ip6)
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
	s, err := IPv4(ctx, cancel)
	if s != "" || err != nil {
		t.Errorf("IPv4() error = %v, want error string", err)
	}
	if want := context.Canceled; !errors.Is(ctx.Err(), want) {
		t.Errorf("IPv4() context.error = %v, want %v", err, want)
	}
}

func TestError(t *testing.T) {
	ctx, timeout := context.WithTimeout(context.Background(), 30*time.Second)
	if _, err := Request(ctx, timeout, "invalid url"); errors.Is(err, nil) {
		t.Errorf("Request() = %v, want an error", err)
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
				t.Errorf("request() error = %v, want %v", err, tt.wantErr)
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
