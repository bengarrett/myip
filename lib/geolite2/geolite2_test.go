package geolite2

import (
	_ "embed"
	"fmt"
	"testing"
)

// example.com
const example = "93.184.216.34"

func BenchmarkCountry(b *testing.B) {
	s, err := Country(example)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(s)
}
func BenchmarkCity(b *testing.B) {
	s, err := City(example)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(s)
}

func TestCountry(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		want    string
		wantErr bool
	}{
		{"empty", "", "", true},
		{"invalid", "1.1.1", "", true},
		{"valid", example, "United States", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Country(tt.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("Country() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Country() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCity(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		want    string
		wantErr bool
	}{
		{"empty", "", "", true},
		{"invalid", "1.1.1", "", true},
		{"valid", example, "Norwell, United States", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := City(tt.ip)
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
