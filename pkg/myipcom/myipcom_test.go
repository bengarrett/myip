package myipcom_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/bengarrett/myip/pkg/myipcom"
)

func BenchmarkRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx, timeout := context.WithTimeout(context.Background(), 5*time.Second)
		s, err := myipcom.Request(ctx, timeout, myipcom.Link)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(s)
	}
}

// ExampleIPv4 demonstrates an IPv4 address request with a 5 second timeout.
func ExampleIPv4() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s, err := myipcom.IPv4(ctx, cancel)
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
		s, err := myipcom.IPv4(ctx4, cancel4)
		if err != nil {
			log.Printf("\n%s\n", err)
		}
		s4 <- s
	}()

	go func() {
		s, err := myipcom.IPv6(ctx6, cancel6)
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
	if _, err := myipcom.IPv4(ctx, timeout); !errors.Is(err, nil) {
		t.Errorf("IPv4() = %v, want %v", err, nil)
	}
}

func TestCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	s, err := myipcom.IPv4(ctx, cancel)
	if s != "" || err != nil {
		t.Errorf("IPv4() error = %v, want error string", err)
	}
	if want := context.Canceled; !errors.Is(ctx.Err(), want) {
		t.Errorf("IPv4() context.error = %v, want %v", err, want)
	}
}

func TestError(t *testing.T) {
	ctx, timeout := context.WithTimeout(context.Background(), 30*time.Second)
	if _, err := myipcom.Request(ctx, timeout, "invalid url"); errors.Is(err, nil) {
		t.Errorf("Request() = %v, want an error", err)
	}
}

func TestRequestS(t *testing.T) {
	tests := []struct {
		name    string
		domain  string
		isValid bool
		wantErr string
	}{
		{"empty", "", false, "unsupported protocol scheme"},
		{"html", "https://example.com", false, "invalid character"},
		{"404", "https://api.myip.com" + "/abcdef", false, "404 not found"},
		{"okay", "https://api.myip.com", true, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, timeout := context.WithTimeout(context.Background(), 5*time.Second)
			gotS, err := myipcom.RequestS(ctx, timeout, tt.domain)
			if err != nil && tt.wantErr != "" && !strings.Contains(fmt.Sprint(err), tt.wantErr) {
				t.Errorf("RequestS() error = %v, want %v", err, tt.wantErr)
			}
			if bool(gotS != "") != tt.isValid {
				t.Errorf("RequestS() = %v, want an ip addr: %v", gotS, tt.isValid)
			}
		})
	}
}

func TestValid(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		wantErr error
	}{
		{"no ip", "", myipcom.ErrNoIP},
		{"invalid", "1.1", myipcom.ErrInvalid},
		{"ipv6", "2001:db8:8714:3a90::12", nil},
		{"ipv4", "1.1.1.1", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := myipcom.Valid(false, tt.ip)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Valid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
