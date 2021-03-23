package main

import (
	"fmt"

	"github.com/bengarrett/myip/lib/ipify"
	"github.com/bengarrett/myip/lib/myipcom"
	"github.com/bengarrett/myip/lib/myipio"
)

func main() {
	if err := ipify.Get(); err != nil {
		fmt.Println(err)
	}
	if err := myipcom.Get(); err != nil {
		fmt.Println(err)
	}
	if s := myipio.IPv4(); s != "" {
		fmt.Println(s)
	}
}
