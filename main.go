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

func version() {
	fmt.Printf("MyIP v%s\n", "0.0")
}

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

func (p ping) count() {
	if p.mode.simple {
		if p.complete > 0 {
			fmt.Printf("\r%s", p.Print)
		}
		return
	}
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

func (p *ping) request1(c chan string) {
	s := ipify.IPv4()
	p.print(s)
	c <- s
}
func (p *ping) request2(c chan string) {
	s := myipcom.IPv4()
	p.print(s)
	c <- s
}
func (p *ping) request3(c chan string) {
	s := myipio.IPv4()
	p.print(s)
	c <- s
}

func (p *ping) request4(c chan string) {
	s := seeip.IPv4()
	p.print(s)
	c <- s
}

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
