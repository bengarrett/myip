package ipv6

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/bengarrett/myip/pkg/ipify"
	"github.com/bengarrett/myip/pkg/myipcom"
	"github.com/bengarrett/myip/pkg/myipio"
	"github.com/bengarrett/myip/pkg/ping"
	"github.com/bengarrett/myip/pkg/seeip"
)

type jobs uint8

const (
	job1 jobs = iota
	job2
	job3
	job4
)

type query struct {
	complete int
	results  []string
	raw      bool
}

func (q *query) worker(ctx context.Context, cancel context.CancelFunc, j jobs, c chan string) {
	ip, err := job(ctx, cancel, j)
	q.complete++
	if err != nil {
		s := ping.Sprints(err.Error(), q.complete, false)
		if q.complete == 1 {
			fmt.Fprint(os.Stdout, s)
		} else {
			fmt.Fprintf(os.Stdout, "\n%s", s)
		}
		c <- ""
		return
	}
	s := ping.Sprints(ip, q.complete, q.raw)
	newIP := !ping.Contains(q.results, ip)
	if newIP {
		q.results = append(q.results, ip)
	}
	if newIP && len(q.results) > 1 {
		fmt.Fprintf(os.Stdout, "\n%s", s)
	} else {
		fmt.Fprint(os.Stdout, s)
	}
	c <- ip
}

// All queries four different services for an IPv4 address and
// as the replies come in, it prints the results to standard output.
// Enabling raw will exclude the city and country location.
func All(timeoutMS int64, raw bool) {
	q := query{raw: raw}
	c := make(chan string)
	timeout := time.Duration(timeoutMS) * time.Millisecond
	ctx1, cancel1 := context.WithTimeout(context.Background(), timeout)
	ctx2, cancel2 := context.WithTimeout(context.Background(), timeout)
	ctx3, cancel3 := context.WithTimeout(context.Background(), timeout)
	ctx4, cancel4 := context.WithTimeout(context.Background(), timeout)
	go q.worker(ctx1, cancel1, job1, c)
	go q.worker(ctx2, cancel2, job2, c)
	go q.worker(ctx3, cancel3, job3, c)
	go q.worker(ctx4, cancel4, job4, c)
	<-c
	<-c
	<-c
	<-c
}

// One queries four different services for an IPv4 address and
// returns the result of the quickest reply. All other requests
// are then aborted.
func One(timeoutMS int64) string {
	c := make(chan string)
	timeout := time.Duration(timeoutMS) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	go worker1(ctx, cancel, job1, c)
	go worker1(ctx, cancel, job2, c)
	go worker1(ctx, cancel, job3, c)
	go worker1(ctx, cancel, job4, c)
	s := ""
	for {
		if s != "" {
			cancel()
			return s
		}
		s = <-c
	}
}

func worker1(ctx context.Context, cancel context.CancelFunc, j jobs, c chan string) {
	s, err := job(ctx, cancel, j)
	if err != nil {
		c <- err.Error()
	}
	c <- s
}

func job(ctx context.Context, cancel context.CancelFunc, j jobs) (string, error) {
	switch j {
	case job1:
		return ipify.IPv6(ctx, cancel)
	case job2:
		return myipcom.IPv6(ctx, cancel)
	case job3:
		return myipio.IPv6(ctx, cancel)
	case job4:
		return seeip.IPv6(ctx, cancel)
	}
	return "", nil
}
