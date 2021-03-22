package main

import (
	"fmt"

	"github.com/bengarrett/myip/lib/ipify"
)

func main() {
	if err := ipify.Get(); err != nil {
		fmt.Println(err)
	}
}
