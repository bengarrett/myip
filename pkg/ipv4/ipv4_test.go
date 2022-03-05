package ipv4_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bengarrett/myip/pkg/ipv4"
)

const timeout = 5 * time.Second

func BenchmarkAll(b *testing.B) {
	to := int64(timeout)
	ipv4.All(to, false)
	ipv4.All(to, true)
	fmt.Println()
}

func BenchmarkOne(b *testing.B) {
	to := int64(timeout)
	fmt.Println(ipv4.One(to))
}
