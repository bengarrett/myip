package main

import (
	"fmt"
	"testing"
)

// example.com
const (
	example = "93.184.216.34"
	norwell = "93.184.216.34, Norwell, United States"
)

//nolint:unparam
func BenchmarkStd(b *testing.B) {
	var p ping
	p.standard()
}

//nolint:unparam
func BenchmarkFirst(b *testing.B) {
	var p ping
	p.mode.first = true
	p.first()
}

//nolint:unparam
func BenchmarkSimple(b *testing.B) {
	var p ping
	p.mode.simple = true
	p.standard()
}

//nolint:unparam
func BenchmarkSimpleAndFirst(b *testing.B) {
	var p ping
	p.mode.first = true
	p.mode.simple = true
	p.first()
}

func Test_self(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"expected", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := self()
			if (err != nil) != tt.wantErr {
				t.Errorf("self() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_ping_count(t *testing.T) {
	e := []string{}
	first := modes{true, false}
	std := modes{false, false}
	simp := modes{false, true}
	type fields struct {
		results  []string
		complete int
		Print    string
		mode     modes
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"empty", fields{}, "(0/4) "},
		{"zero", fields{e, 0, "", std}, "(0/4) "},
		{"one", fields{e, 1, "", std}, "\r(1/4) "},
		{"example1", fields{e, 1, example, std}, fmt.Sprintf("\r(1/4) %s", example)},
		{"example4", fields{e, 4, example, std}, fmt.Sprintf("\r(4/4) %s", example)},
		{"example5", fields{e, 5, example, std}, fmt.Sprintf("\r(5/4) %s", example)},
		{"first0", fields{e, 0, "", first}, "(0/1) "},
		{"first1", fields{e, 1, example, first}, fmt.Sprintf("\r(1/1) %s", example)},
		{"simple0", fields{e, 0, example, simp}, ""},
		{"simple1", fields{e, 1, example, simp}, fmt.Sprintf("\r%s", example)},
		{"simple4", fields{e, 4, example, simp}, fmt.Sprintf("\r%s", example)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ping{
				results:  tt.fields.results,
				complete: tt.fields.complete,
				Print:    tt.fields.Print,
				mode:     tt.fields.mode,
			}
			if got := p.count(); got != tt.want {
				t.Errorf("ping.count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contains(t *testing.T) {
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
			if got := contains(tt.args.a, tt.args.x); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ping_simple(t *testing.T) {
	const ok = example
	type fields struct {
		results  []string
		complete int
		Print    string
		mode     modes
	}
	tests := []struct {
		name   string
		fields fields
		ip     string
		want   string
	}{
		{"empty", fields{}, "", ""},
		{"single ip", fields{}, example, example},
		{"multiple ips", fields{[]string{"a", "b"}, 2, ok, modes{}}, example, fmt.Sprintf("%s. %s", ok, example)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ping{
				results:  tt.fields.results,
				complete: tt.fields.complete,
				Print:    tt.fields.Print,
				mode:     tt.fields.mode,
			}
			if got := p.simple(tt.ip); got != tt.want {
				t.Errorf("ping.simple() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ping_city(t *testing.T) {
	const ok = norwell
	type fields struct {
		results  []string
		complete int
		Print    string
		mode     modes
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
		{"examples", fields{[]string{"a", "b"}, 2, ok, modes{}}, example, fmt.Sprintf("%s. %s", ok, ok), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ping{
				results:  tt.fields.results,
				complete: tt.fields.complete,
				Print:    tt.fields.Print,
				mode:     tt.fields.mode,
			}
			got, err := p.city(tt.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("ping.city() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ping.city() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ping_parse(t *testing.T) {
	type fields struct {
		results  []string
		complete int
		Print    string
		mode     modes
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
			p := &ping{
				results:  tt.fields.results,
				complete: tt.fields.complete,
				Print:    tt.fields.Print,
				mode:     tt.fields.mode,
			}
			if got := p.parse(tt.ip); got != tt.want {
				t.Errorf("ping.parse(%v) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}
