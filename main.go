package main

import (
	"fmt"

	"github.com/bengarrett/myip/lib/ipify"
	"github.com/bengarrett/myip/lib/myipcom"
	"github.com/bengarrett/myip/lib/myipio"
)

func main() {
	if s := ipify.IPv4(); s != "" {
		fmt.Println(s)
	}
	if s := myipcom.IPv4(); s != "" {
		fmt.Println(s)
	}
	if s := myipio.IPv4(); s != "" {
		fmt.Println(s)
	}
}
