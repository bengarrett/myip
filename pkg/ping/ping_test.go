package ping_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bengarrett/myip/pkg/ping"
)

const (
	example = "93.184.216.34"
	norwell = "93.184.216.34, Norwell, United States"
)

func TestCity(t *testing.T) {
	const ok = norwell
	type fields struct {
		Print   string
		Results []string
	}
	tests := []struct {
		name    string
		fields  fields
		ip      string
		want    string
		wantErr bool
	}{
		{"empty", fields{}, "", "", true},
		{"example", fields{}, example, ok, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ping.City(tt.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("City() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("City() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	type args struct {
		a []string
		x string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"null", args{}, false},
		{"empty", args{[]string{}, "x"}, false},
		{"no match", args{[]string{"alpha", "beta", "c", "e", "0"}, "x"}, false},
		{"no match", args{[]string{"alpha", "beta", "xx", "e", "0"}, "x"}, false},
		{"match", args{[]string{"alpha", "beta", "c", "x", "e", "0"}, "x"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ping.Contains(tt.args.a, tt.args.x); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSprint(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want string
	}{
		{"empty", "", ""},
		{"invalid", "a.b.c.d", "(1/1) : invalid ip address"},
		{"no geo-location", "0.0.0.0", "(1/1) 0.0.0.0"},
		{"example", example, fmt.Sprintf("(1/1) %s", norwell)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strings.TrimSpace(ping.Sprint(tt.ip)); got != tt.want {
				t.Errorf("Sprint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSprints(t *testing.T) {
	type args struct {
		ip        string
		completed int
		raw       bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{}, ""},
		{"invalid", args{"a.b.c.d", 1, false}, "(1/4) a.b.c.d, invalid ip address"},
		{"invalid raw", args{"a.b.c.d", 1, true}, "(1/4) a.b.c.d"},
		{"no geo-location", args{"0.0.0.0", 2, false}, "(2/4) 0.0.0.0"},
		{"example", args{example, 4, false}, fmt.Sprintf("(4/4) %s", norwell)},
		{"example raw", args{example, 4, true}, fmt.Sprintf("(4/4) %s", example)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strings.TrimSpace(ping.Sprints(tt.args.ip, tt.args.completed, tt.args.raw)); got != tt.want {
				t.Errorf("Sprints() = %v, want %v", got, tt.want)
			}
		})
	}
}
