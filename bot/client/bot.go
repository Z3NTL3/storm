package client

import (
	lb "github.com/hlts2/round-robin"
)

type Bot struct {
	proxies lb.RoundRobin
	accepts lb.RoundRobin
	headers lb.RoundRobin
}
