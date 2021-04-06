package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bengarrett/myip/lib/geolite2"
	"github.com/bengarrett/myip/lib/ipify"
	"github.com/bengarrett/myip/lib/myipcom"
	"github.com/bengarrett/myip/lib/myipio"
	"github.com/bengarrett/myip/lib/seeip"
)

type ping struct {
	results  []string
	complete int
	Print    string
	mode     modes
}

type modes struct {
	first  bool
	simple bool
}

type jobs uint8

const (
	job1 jobs = iota
	job2
	job3
	job4
)

func main() {
	var p ping
	flag.BoolVar(&p.mode.first, "first", false, "Returns the first reported IP address, its location and exits.")
	flag.BoolVar(&p.mode.simple, "simple", false, "Simple mode only displays an IP address and exits.")
	ver := flag.Bool("version", false, "Version and information for this program.")
	flag.Parse()
	// version information
	if *ver {
		version()
		return
	}
	// first mode
	if p.mode.first {
		p.first()
		return
	}
	// standard mode
	p.standard()
}

func self() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("self error: %w", err)
	}
	return filepath.Dir(exe), nil
}

// Version prints out the program information and version.
func version() {
	const copyright = "\u00A9"
	const app = "0.0"
	fmt.Printf("MyIP v%s\n%s 2021 Ben Garrett\n\n", app, copyright)
	fmt.Println("Web:  https://github.com/bengarrett/myip")
	exe, err := self()
	if err != nil {
		fmt.Printf("Path: %s\n", err)
		return
	}
	fmt.Printf("Path: %s\n", exe)
}

// Fast waits for the fastest concurrent request to complete
// before aborting and closing the others.
func (p ping) first() {
	fmt.Print(p.count())
	c := make(chan string)
	go p.worker(job1, c)
	go p.worker(job2, c)
	go p.worker(job3, c)
	go p.worker(job4, c)
	<-c
	close(c)
	fmt.Println()
}

// Standard waits for all the concurrent requests to complete.
func (p ping) standard() {
	fmt.Print(p.count())
	c := make(chan string)
	go p.worker(job1, c)
	go p.worker(job2, c)
	go p.worker(job3, c)
	go p.worker(job4, c)
	_, _, _, _ = <-c, <-c, <-c, <-c
	fmt.Println()
}

// Count prints out the request job counts and the resulting IP addresses.
func (p ping) count() string {
	// simple prints only the ip addresses
	if p.mode.simple {
		if p.complete > 0 {
			return fmt.Sprintf("\r%s", p.Print)
		}
		return ""
	}
	// standard prints the ip addresses with request complete counts
	total := 4
	if p.mode.first {
		total = 1
	}
	if p.complete == 0 {
		return fmt.Sprintf("(0/%d) ", total)
	}
	// (1/4) 93.184.216.34, Norwell, United States
	return fmt.Sprintf("\r(%d/%d) %s", p.complete, total, p.Print)
}

func (p *ping) worker(i jobs, c chan string) {
	var s string
	switch i {
	case job1:
		s = ipify.IPv4()
	case job2:
		s = myipcom.IPv4()
	case job3:
		s = myipio.IPv4()
	case job4:
		s = seeip.IPv4()
	}
	fmt.Print(p.parse(s))
	c <- s
}

// Print only unique IP addresses from the request results.
func (p *ping) parse(ip string) string {
	p.complete++
	if ip == "" {
		return ""
	}
	if !contains(p.results, ip) {
		p.results = append(p.results, ip)
		if p.mode.simple {
			p.Print = p.simple(ip)
			return p.count()
		}
		var err error
		p.Print, err = p.city(ip)
		if err != nil {
			return fmt.Sprintf("city %q error: %s\n", ip, err)
		}
	}
	return p.count()
}

// City prints the IP address with its geographic location,
// both the country and the city.
func (p ping) city(ip string) (string, error) {
	c, err := geolite2.City(ip)
	if err != nil {
		return "", err
	}
	if len(p.results) > 1 {
		return fmt.Sprintf("%s. %s, %s", p.Print, ip, c), nil
	}
	return fmt.Sprintf("%s, %s", ip, c), nil
}

// PrintSimple prints the IP address.
func (p ping) simple(ip string) string {
	if len(p.results) > 1 {
		return fmt.Sprintf("%s. %s", p.Print, ip)
	}
	return ip
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
