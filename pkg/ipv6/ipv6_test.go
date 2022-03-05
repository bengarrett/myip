package ipv6_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bengarrett/myip/pkg/ipv6"
)

const timeout = 5 * time.Second

func BenchmarkAll(b *testing.B) {
	to := int64(timeout)
	ipv6.All(to, false)
	ipv6.All(to, true)
	fmt.Println()
}

func BenchmarkOne(b *testing.B) {
	to := int64(timeout)
	fmt.Println(ipv6.One(to))
}
