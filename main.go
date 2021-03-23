package main

import (
	"fmt"

	"github.com/bengarrett/myip/lib/geolite2"
	"github.com/bengarrett/myip/lib/ipify"
	"github.com/bengarrett/myip/lib/myipcom"
	"github.com/bengarrett/myip/lib/myipio"
)

type Queries struct {
	Results []string
	Done    int
	Print   string
}

func (q *Queries) request1(c chan string) {
	s := ipify.IPv4()
	q.store(s)
	c <- s
}
func (q *Queries) request2(c chan string) {
	s := myipcom.IPv4()
	q.store(s)
	c <- s
}
func (q *Queries) request3(c chan string) {
	s := myipio.IPv4()
	q.store(s)
	c <- s
}

func (q *Queries) request4(c chan string) {
	s := "7.134.10.1"
	q.store(s)
	c <- s
}

func (q *Queries) store(ip string) {
	if !contains(q.Results, ip) {
		q.Results = append(q.Results, ip)
		c, err := geolite2.City(ip)
		fmt.Println(err)
		if len(q.Results) > 1 {
			q.Print = fmt.Sprintf("%s. %s, %s", q.Print, ip, c)
			return
		}
		q.Print = fmt.Sprintf("%s, %s", ip, c)
	}
	q.Done++
	count(q.Done, q.Print)
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func count(i int, s string) {
	if i == 0 {
		fmt.Print("(0/4) ")
		return
	}
	fmt.Printf("\r(%d/4) %s", i, s)
}

func main() {
	var q Queries
	count(0, "")
	c := make(chan string)
	go q.request1(c)
	go q.request2(c)
	go q.request3(c)
	go q.request4(c)
	_, _, _, _ = <-c, <-c, <-c, <-c
	fmt.Println()
	// fmt.Println(z, y, x)
	//fmt.Println(ips)
}
