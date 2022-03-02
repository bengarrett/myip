// Package main determine your Internet-facing
// IP address and location from multiple sources.
// © Ben Garrett https://github.com/bengarrett/myip
package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/bengarrett/myip/pkg/ping"
)

type modes struct {
	first   bool
	ipv6    bool
	timeout int64
}

const (
	// Default HTTP request timeout value in milliseconds.
	httpTimeout = 5000
	// Tabwriter padding using spaces.
	padding = 4
)

var (
	version = "0.0.0"
	commit  = "unset" // nolint: gochecknoglobals
	date    = "unset" // nolint: gochecknoglobals
)

func main() {
	msInSec := func(i int) int {
		const second = 1000
		return i / second
	}
	var (
		mode modes
		p    ping.Ping
	)
	flag.BoolVar(&mode.first, "first", false, "returns the first reported IP address and its location")
	flag.BoolVar(&mode.ipv6, "ipv6", false, "return an IPv6 address instead of IPv4")
	flag.BoolVar(&p.Raw, "simple", false, "simple mode only displays the IP address")
	flag.Int64Var(&mode.timeout, "timeout", httpTimeout,
		fmt.Sprintf("https request timeout in milliseconds (default: %d [%d seconds])", httpTimeout, msInSec(httpTimeout)))
	ver := flag.Bool("version", false, "version and information for this program")
	f := flag.Bool("f", false, "alias for first")
	i := flag.Bool("i", false, "alias for ipv6")
	s := flag.Bool("s", false, "alias for simple")
	t := flag.Int64("t", 0, "alias for timeout")
	v := flag.Bool("v", false, "alias for version")

	flag.Usage = func() {
		const alias = 1
		fmt.Fprintln(os.Stderr, "MyIP Usage:")
		fmt.Fprintln(os.Stderr, "    myip [options]:")
		fmt.Fprintln(os.Stderr, "")
		w := tabwriter.NewWriter(os.Stderr, 0, 0, padding, ' ', 0)
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
	// aliases
	if *i {
		mode.ipv6 = true
	}
	if *s {
		p.Raw = true
	}
	if *t > 0 {
		mode.timeout = *t
	}
	// first mode
	if mode.first || *f {
		mode.first = true
		fmt.Print(p.InitRequest())
		r := p.Request(mode.timeout, mode.ipv6)
		s := p.Parse(r) // rename to FmtRequest
		fmt.Println(s)
		return
	}
	// standard mode
	p.Standard(mode.timeout, mode.ipv6)
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
