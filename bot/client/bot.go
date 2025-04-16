package client

import (
	roundrobin "github.com/hlts2/round-robin"
)

type Bot struct {
	proxies roundrobin.RoundRobin
	accepts roundrobin.RoundRobin
	headers roundrobin.RoundRobin
}
