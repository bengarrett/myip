package main

import (
	"flag"
	"fmt"

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
	fast   bool
	simple bool
}

func main() {
	var p ping
	flag.BoolVar(&p.mode.simple, "simple", false, "simple mode only displays an IP address and exits")
	flag.BoolVar(&p.mode.simple, "s", false, "")
	flag.BoolVar(&p.mode.fast, "fast", false, "fast mode returns the first reported IP address and exits")
	flag.BoolVar(&p.mode.fast, "f", false, "")
	ver := flag.Bool("version", false, "version and information for this program")
	v := flag.Bool("v", false, "")
	flag.Parse()
	// version information
	if *ver || *v {
		version()
		return
	}
	// fast mode
	if p.mode.fast {
		p.fast()
		return
	}
	// standard mode
	p.standard()
}

// Version prints out the program information and version.
func version() {
	const copyright = "\u00A9"
	const app = "0.0"
	fmt.Printf("MyIP v%s\n%s 2021 Ben Garrett\n\n", app, copyright)
	fmt.Println("https://github.com/bengarrett/myip")
}

// Fast waits for the fastest concurrent request to complete
// then aborts and closes the others.
func (p ping) fast() {
	p.count()
	c := make(chan string)
	go p.request1(c)
	go p.request2(c)
	go p.request3(c)
	go p.request4(c)
	<-c
	close(c)
	fmt.Println()
}

// Standard waits for all the concurrent requests to complete.
func (p ping) standard() {
	p.count()
	c := make(chan string)
	go p.request1(c)
	go p.request2(c)
	go p.request3(c)
	go p.request4(c)
	_, _, _, _ = <-c, <-c, <-c, <-c
	fmt.Println()
}

// Count prints out the request job counts and the resulting IP addresses.
func (p ping) count() {
	// simple prints only the ip addresses
	if p.mode.simple {
		if p.complete > 0 {
			fmt.Printf("\r%s", p.Print)
		}
		return
	}
	// standard prints the ip addresses with request complete counts
	total := 4
	if p.mode.fast {
		total = 1
	}
	if p.complete == 0 {
		fmt.Printf("(0/%d) ", total)
		return
	}
	fmt.Printf("\r(%d/%d) %s", p.complete, total, p.Print)
}

// Request1 pings ipify.org.
func (p *ping) request1(c chan string) {
	s := ipify.IPv4()
	p.print(s)
	c <- s
}

// Request2 pings myip.com.
func (p *ping) request2(c chan string) {
	s := myipcom.IPv4()
	p.print(s)
	c <- s
}

// Request3 pings my-ip.io.
func (p *ping) request3(c chan string) {
	s := myipio.IPv4()
	p.print(s)
	c <- s
}

// Request4 pings seeip.org.
func (p *ping) request4(c chan string) {
	s := seeip.IPv4()
	p.print(s)
	c <- s
}

// Print only unique IP addresses from the request results.
func (p *ping) print(ip string) {
	p.complete++
	if !contains(p.results, ip) {
		p.results = append(p.results, ip)
		if p.mode.simple {
			p.printSimple(ip)
		} else {
			p.printCity(ip)
		}
	}
	p.count()
}

// PrintCity prints the IP address with its geographic location,
// both the country and the city.
func (p *ping) printCity(ip string) {
	c, err := geolite2.City(ip)
	if err != nil {
		fmt.Println(err)
	}
	if len(p.results) > 1 {
		p.Print = fmt.Sprintf("%s. %s, %s", p.Print, ip, c)
		return
	}
	p.Print = fmt.Sprintf("%s, %s", ip, c)
}

// PrintSimple prints the IP address.
func (p *ping) printSimple(ip string) {
	if len(p.results) > 1 {
		p.Print = fmt.Sprintf("%s. %s", p.Print, ip)
		return
	}
	p.Print = ip
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
