package robin_test

import (
	"fmt"
	"log"
	"testing"
	"z3ntl3/storm/lb/robin"
)

func TestRobin(t *testing.T) {
	lb := robin.New(1, 2, 3, 4)
	for p := range lb.Iter() {
		fmt.Printf("%d\n", *p)
	}

	reset := lb.Next()
	if *reset != 1 {
		log.Fatal("fail")
	}

	fmt.Printf("reset: %d\n", *reset)
}
