package main

import (
	"fmt"

	"github.com/bengarrett/myip/lib/ipify"
	"github.com/bengarrett/myip/lib/myipcom"
)

func main() {
	if err := ipify.Get(); err != nil {
		fmt.Println(err)
	}
	if err := myipcom.Get(); err != nil {
		fmt.Println(err)
	}
}
