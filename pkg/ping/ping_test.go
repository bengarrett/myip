package ping_test

import (
	"fmt"
	"testing"

	"github.com/bengarrett/myip/pkg/ping"
)

const (
	example = "93.184.216.34"
	norwell = "93.184.216.34, Norwell, United States"
	timeout = 5000
)

//nolint:unparam
func BenchmarkStd(b *testing.B) {
	var p ping.Ping
	p.Standard(timeout, false)
}

//nolint:unparam
func BenchmarkFirst(b *testing.B) {
	var p ping.Ping
	p.Request(timeout, false)
}

//nolint:unparam
func BenchmarkSimple(b *testing.B) {
	var p ping.Ping
	p.Standard(timeout, false)
}

//nolint:unparam
func BenchmarkSimpleAndFirst(b *testing.B) {
	var p ping.Ping
	p.Raw = true
	p.Request(timeout, false)
}

func Test_Contains(t *testing.T) {
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

func Test_ping_Simple(t *testing.T) {
	const ok = example
	type fields struct {
		Results []string
		Print   string
	}
	tests := []struct {
		name   string
		fields fields
		ip     string
		want   string
	}{
		{"empty", fields{}, "", ""},
		{"single ip", fields{}, example, example},
		{"multiple ips", fields{[]string{"a", "b"}, ok}, example, fmt.Sprintf("%s. %s", ok, example)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ping.Ping{
				Results: tt.fields.Results,
				Print:   tt.fields.Print,
			}
			if got := p.Simple(tt.ip); got != tt.want {
				t.Errorf("ping.Simple() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ping_City(t *testing.T) {
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
		{"examples", fields{Print: ok, Results: []string{"a", "b"}},
			example, fmt.Sprintf("%s. %s", ok, ok), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ping.Ping{
				Print:   tt.fields.Print,
				Results: tt.fields.Results,
			}
			got, err := p.City(tt.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("ping.City() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ping.City() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ping_Parse(t *testing.T) {
	type fields struct {
		Print string
	}
	tests := []struct {
		name   string
		fields fields
		ip     string
		want   string
	}{
		{"empty", fields{}, "", ""},
		{"ip", fields{}, example, fmt.Sprintf("\r(1/4) %s", norwell)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ping.Ping{
				Print: tt.fields.Print,
			}
			if got := p.Parse(tt.ip); got != tt.want {
				t.Errorf("ping.Parse(%v) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}
