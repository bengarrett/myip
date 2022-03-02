package ping

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bengarrett/myip/pkg/geolite2"
	"github.com/bengarrett/myip/pkg/ipify"
	"github.com/bengarrett/myip/pkg/myipcom"
	"github.com/bengarrett/myip/pkg/myipio"
	"github.com/bengarrett/myip/pkg/seeip"
)

type jobs uint8

const (
	job1 jobs = iota
	job2
	job3
	job4
)

type Ping struct {
	Complete int
	Print    string
	Raw      bool
	Results  []string
}

// City prints the IP address with its geographic location,
// both the country and the city.
func (p Ping) City(ip string) (string, error) {
	c, err := geolite2.City(ip)
	if err != nil {
		return "", err
	}
	if len(p.Results) > 1 {
		return fmt.Sprintf("%s. %s, %s", p.Print, ip, c), nil
	}
	return fmt.Sprintf("%s, %s", ip, c), nil
}

// InitRequest
func (p Ping) InitRequest() string {
	if p.Raw {
		if p.Complete > 0 {
			return fmt.Sprintf("\r%s", p.Print)
		}
		return ""
	}
	const total = 1
	if p.Complete == 0 {
		return fmt.Sprintf("(0/%d) ", total)
	}
	p.Complete = 1
	// (1/1) 93.184.216.34, Norwell, United States
	return fmt.Sprintf("\r(%d/%d) %s", p.Complete, total, p.Print)
}

// Count prints out the request job counts and the resulting IP addresses.
func (p Ping) count() string {
	// simple prints only the ip addresses
	if p.Raw {
		if p.Complete > 0 {
			return fmt.Sprintf("\r%s", p.Print)
		}
		return ""
	}
	// standard prints the ip addresses with request complete counts
	total := 4
	if p.Complete == 0 {
		return fmt.Sprintf("(0/%d) ", total)
	}
	// (1/4) 93.184.216.34, Norwell, United States
	return fmt.Sprintf("\r(%d/%d) %s", p.Complete, total, p.Print)
}

// Print only unique IP addresses from the request results.
func (p *Ping) Parse(ip string) string {
	p.Complete++
	if ip == "" {
		return ""
	}
	if !Contains(p.Results, ip) {
		p.Results = append(p.Results, ip)
		if p.Raw {
			p.Print = p.Simple(ip)
			return p.count()
		}
		var err error
		p.Print, err = p.City(ip)
		if err != nil {
			return fmt.Sprintf("city %q error: %s\n", ip, err)
		}
	}
	return p.count()
}

// Simple prints the IP address.
func (p Ping) Simple(ip string) string {
	if len(p.Results) > 1 {
		return fmt.Sprintf("%s. %s", p.Print, ip)
	}
	return ip
}

// Fast waits for the fastest concurrent request to complete
// before aborting and canceling the others.
func (p Ping) Request(timeoutMS int64, ipv6 bool) string {
	//fmt.Print(p.countToFirst())
	c := make(chan string)
	timeout := time.Duration(timeoutMS) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	if ipv6 {
		go p.workerV6(ctx, cancel, job1, c)
		go p.workerV6(ctx, cancel, job2, c)
		go p.workerV6(ctx, cancel, job3, c)
		go p.workerV6(ctx, cancel, job4, c)
		s := <-c
		cancel()
		<-c
		<-c
		<-c
		fmt.Println()
		return s
	}
	go p.workerV4(ctx, cancel, job1, c)
	go p.workerV4(ctx, cancel, job2, c)
	go p.workerV4(ctx, cancel, job3, c)
	go p.workerV4(ctx, cancel, job4, c)
	s := <-c
	cancel()
	<-c
	<-c
	<-c
	return s
}

// Standard waits for all the concurrent requests to complete.
func (p Ping) Standard(timeoutMS int64, ipv6 bool) {
	fmt.Print(p.count())
	c := make(chan string)
	timeout := time.Duration(timeoutMS) * time.Millisecond
	ctx1, cancel1 := context.WithTimeout(context.Background(), timeout)
	ctx2, cancel2 := context.WithTimeout(context.Background(), timeout)
	ctx3, cancel3 := context.WithTimeout(context.Background(), timeout)
	ctx4, cancel4 := context.WithTimeout(context.Background(), timeout)
	if ipv6 {
		go p.workerV6(ctx1, cancel1, job1, c)
		go p.workerV6(ctx2, cancel2, job2, c)
		go p.workerV6(ctx3, cancel3, job3, c)
		go p.workerV6(ctx4, cancel4, job4, c)
		<-c
		<-c
		<-c
		<-c
		fmt.Println()
		return
	}
	go p.workerV4(ctx1, cancel1, job1, c)
	go p.workerV4(ctx2, cancel2, job2, c)
	go p.workerV4(ctx3, cancel3, job3, c)
	go p.workerV4(ctx4, cancel4, job4, c)
	<-c
	<-c
	<-c
	<-c
	fmt.Println()
}

func (p *Ping) workerV4(ctx context.Context, cancel context.CancelFunc, job jobs, c chan string) {
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
	//fmt.Print(p.Parse(s))
	c <- s
}

func (p *Ping) workerV6(ctx context.Context, cancel context.CancelFunc, job jobs, c chan string) {
	var s string
	var err error
	switch job {
	case job1:
		s, err = ipify.IPv6(ctx, cancel)
	case job2:
		s, err = myipcom.IPv6(ctx, cancel)
	case job3:
		s, err = myipio.IPv6(ctx, cancel)
	case job4:
		s, err = seeip.IPv6(ctx, cancel)
	}
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Print(p.Parse(s))
	c <- s
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
