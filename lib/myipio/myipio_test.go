package myipio

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
		fmt.Println(s)
	}
}

// ExampleIPv4 demonstrates an IPv4 address request with a 5 second timeout.
func ExampleIPv4() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s, err := IPv4(ctx, cancel)
	if err != nil {
		log.Fatalf("\n%s\n", err)
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
	want := context.Canceled
	s, err := IPv4(ctx, cancel)
	if s != "" || err != nil {
		t.Errorf("IPv4() s = %v, error = %v, want an empty string with no errors", s, err)
	}
	if !errors.Is(ctx.Err(), want) {
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
		isValid bool
		wantErr string
	}{
		{"empty", "", false, "unsupported protocol scheme"},
		{"html", "https://example.com", false, "invalid character"},
		{"404", "https://api4.my-ip.io/ip.json/abcdef", false, "404 not found"},
		{"okay", "https://api4.my-ip.io/ip.json", true, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, timeout := context.WithTimeout(context.Background(), 5*time.Second)
			gotS, err := request(ctx, timeout, tt.domain)
			if err != nil && tt.wantErr != "" && !strings.Contains(fmt.Sprint(err), tt.wantErr) {
				t.Errorf("request() error = %v, want %v", err, tt.wantErr)
			}
			if bool(gotS.IP != "") != tt.isValid {
				t.Errorf("request() = %v, want an ip addr: %v", gotS, tt.isValid)
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
			got, err := r.valid(linkv4)
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
