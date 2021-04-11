package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

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

var (
	version = "0.0.0" // nolint:gochecknoglobals
	commit  = "unset" // nolint:gochecknoglobals
	date    = "unset" // nolint:gochecknoglobals
)

func main() {
	var p ping
	flag.BoolVar(&p.mode.first, "first", false, "returns the first reported IP address and its location")
	flag.BoolVar(&p.mode.simple, "simple", false, "simple mode only displays the IP address")
	ver := flag.Bool("version", false, "version and information for this program")
	f := flag.Bool("f", false, "alias for first")
	s := flag.Bool("s", false, "alias for simple")
	v := flag.Bool("v", false, "alias for version")

	flag.Usage = func() {
		const alias = 1
		fmt.Fprintln(os.Stderr, "MyIP Usage:")
		fmt.Fprintln(os.Stderr, "    myip [options]:")
		fmt.Fprintln(os.Stderr, "")
		w := tabwriter.NewWriter(os.Stderr, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "    -h, --help\tshow this list of options")
		flag.VisitAll(func(f *flag.Flag) {
			if len(f.Name) == alias {
				return
			}
			fmt.Fprintf(w, "    -%v, --%v\t%v\n", f.Name[:1], f.Name, f.Usage)
		})
		w.Flush()
	}
	flag.Parse()

	// version information
	if *ver || *v {
		info()
		return
	}
	// simple alias
	if *s {
		p.mode.simple = true
	}
	// first mode
	if p.mode.first || *f {
		p.mode.first = true
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
	return exe, nil
}

// Info prints out the program information and version.
func info() {
	const copyright = "\u00A9"
	fmt.Printf("MyIP Tetrad v%s\n%s 2021 Ben Garrett\n", version, copyright)
	fmt.Printf("https://github.com/bengarrett/myip\n\n")
	fmt.Printf("build: %s (%s)\n", commit, date)
	exe, err := self()
	if err != nil {
		fmt.Printf("path: %s\n", err)
		return
	}
	fmt.Printf("path:  %s\n", exe)
}

// Fast waits for the fastest concurrent request to complete
// before aborting and cancelling the others.
func (p ping) first() {
	fmt.Print(p.count())
	c := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	go p.worker(ctx, cancel, job1, c)
	go p.worker(ctx, cancel, job2, c)
	go p.worker(ctx, cancel, job3, c)
	go p.worker(ctx, cancel, job4, c)
	<-c
	cancel()
	<-c
	<-c
	<-c
	fmt.Println()
}

// Standard waits for all the concurrent requests to complete.
func (p ping) standard() {
	fmt.Print(p.count())
	c := make(chan string)
	ctx1, to1 := context.WithTimeout(context.Background(), 1*time.Second)
	ctx2, to2 := context.WithTimeout(context.Background(), 5*time.Second)
	ctx3, to3 := context.WithTimeout(context.Background(), 5*time.Second)
	ctx4, to4 := context.WithTimeout(context.Background(), 5*time.Second)
	go p.worker(ctx1, to1, job1, c)
	go p.worker(ctx2, to2, job2, c)
	go p.worker(ctx3, to3, job3, c)
	go p.worker(ctx4, to4, job4, c)
	<-c
	<-c
	<-c
	<-c
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

func (p *ping) worker(ctx context.Context, cancel context.CancelFunc, job jobs, c chan string) {
	var s string
	var err error
	switch job {
	case job1:
		s, err = ipify.IPv4(ctx, cancel)
	case job2:
		s, err = myipcom.IPv4(ctx, cancel)
	case job3:
		s, err = myipio.IPv4(ctx, cancel)
	case job4:
		s, err = seeip.IPv4(ctx, cancel)
	}
	if err != nil {
		log.Fatalf("\n%s\n", err)
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

// Simple prints the IP address.
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
