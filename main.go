package main

import (
	"fmt"

	"github.com/bengarrett/myip/lib/geolite2"
	"github.com/bengarrett/myip/lib/ipify"
	"github.com/bengarrett/myip/lib/myipcom"
	"github.com/bengarrett/myip/lib/myipio"
)

var ips []string
var done int
var list string

//sort.SearchStrings(a, "A")

func request1(c chan string) {
	s := ipify.IPv4()
	store(s)
	done++
	count(done)
	c <- s
}
func request2(c chan string) {
	s := myipcom.IPv4()
	store(s)
	done++
	count(done)
	c <- s
}
func request3(c chan string) {
	s := myipio.IPv4()
	store(s)
	done++
	count(done)
	c <- s
}

func request4(c chan string) {
	s := "7.134.10.1"
	store(s)
	done++
	count(done)
	c <- s
}

func store(ip string) {
	if !contains(ips, ip) {
		ips = append(ips, ip)
		c, _ := geolite2.City(ip)
		if len(ips) > 1 {
			list = fmt.Sprintf("%s. %s, %s", list, ip, c)
			return
		}
		list = fmt.Sprintf("%s, %s", ip, c)
	}
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func count(i int) {
	if i == 0 {
		fmt.Print("(0/4) ")
		return
	}
	fmt.Printf("\r(%d/4) %s", i, list)
}

func main() {
	count(0)
	c := make(chan string)
	go request1(c)
	go request2(c)
	go request3(c)
	go request4(c)
	_, _, _, _ = <-c, <-c, <-c, <-c
	fmt.Println()
	// fmt.Println(z, y, x)
	//fmt.Println(ips)
}
