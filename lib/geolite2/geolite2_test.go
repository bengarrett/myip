package geolite2

import (
	_ "embed"
	"fmt"
	"log"
	"testing"
)

type jobs uint8

const (
	cities jobs = iota
	countries
)

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

func ExampleCity() {
	s, err := City("93.184.216.34")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)
	// Output: Norwell, United States
}

func ExampleCountry() {
	s, err := Country("93.184.216.34")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)
	// Output: United States
}

func TestLocations(t *testing.T) {
	tests := []struct {
		name    string
		job     jobs
		ip      string
		want    string
		wantErr bool
	}{
		{"empty country", countries, "", "", true},
		{"invalid country", countries, "1.1.1", "", true},
		{"valid country", countries, example, "United States", false},
		{"empty city", cities, "", "", true},
		{"invalid city", cities, "1.1.1", "", true},
		{"valid city", cities, example, "Norwell, United States", false},
	}
	var got string
	var err error
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.job {
			case cities:
				got, err = City(tt.ip)
			case countries:
				got, err = Country(tt.ip)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Locations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Locations() = %v, want %v", got, tt.want)
			}
		})
	}
}
