package client

import (
	lb "github.com/hlts2/round-robin"
)

type Bot struct {
	rsrc resources
}

type resources struct {
	proxies lb.RoundRobin
	accepts lb.RoundRobin
	headers lb.RoundRobin
}
