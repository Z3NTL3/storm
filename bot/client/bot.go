package client

import (
	lb "github.com/thegeekyasian/round-robin-go"
)

type Bot struct {
	rsrc resources
}

type resources struct {
	proxies lb.RoundRobin[string]
	accepts lb.RoundRobin[string]
	headers lb.RoundRobin[string]
}
