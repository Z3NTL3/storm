package client

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
	"z3ntl3/storm/globals"
	lb "z3ntl3/storm/lb/robin"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

type Bot struct {
	Rsrc resources
}

type resources struct {
	Proxies *lb.RoundRobin[string]
	Accepts *lb.RoundRobin[string]
	Headers *lb.RoundRobin[string]
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
	for i, path_ := range []string{globals.Accepts, globals.Headers, globals.ProxyFile} {
		f, err := os.Open(path.Join(cwd, path_))
		if err != nil {
			log.Fatal(err)
		}

		contents, err := io.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}

		contents = []byte(strings.Trim(string(contents), "\r\n"))
		parts := strings.Split(string(contents), "\n")
		lb := lb.New(parts...)

		switch i {
		case 0:
			instance.Rsrc.Accepts = lb
		case 1:
			instance.Rsrc.Headers = lb
		case 2:
			instance.Rsrc.Proxies = lb
		default:
			log.Fatal("could not match any data to use for the stress test")
		}
	}

	return instance
}

// Dedicated to be spawned on a goroutine
// The body is thread-safe
//
// Arguments should be passed by [Bot]'s Next method on the fields as they comfort [lb.RoundRobin]
func (c *Bot) Stress(proxy string, th_id uint64, pool_msg chan MessageContext) {
	var client fasthttp.Client

	proxyURI, err := url.Parse(proxy)
	if err != nil {
		pool_msg <- MessageContext{
			Err: err,
		}
	}

	switch proxyURI.Scheme {
	case "http", "https":
		client.Dial = fasthttpproxy.FasthttpHTTPDialerDualStack(proxy)
	case "socks5":
		client.Dial = fasthttpproxy.FasthttpSocksDialerDualStack(proxy)
	default:
		pool_msg <- MessageContext{
			Err:  errors.New("unsupported proxy protocol: socks4, use socks5 or http/https instead"),
			Kill: true,
		}
	}

	req := fasthttp.AcquireRequest()
	for header := range c.Rsrc.Headers.Iter() {
		if !strings.Contains(*header, " ") {
			continue
		}
		h := strings.Split(*header, " ")
		req.Header.Set(h[0], h[1])
	}

	for accept := range c.Rsrc.Accepts.Iter() {
		req.Header.Set("Accept", *accept)
	}

	// do not observe response as to save memory
	err = client.Do(req, nil)
	if err != nil {
		pool_msg <- MessageContext{
			Msg: fmt.Sprintf("thread [%d] sending payload responded with: %s", th_id, err),
		}
		return
	}

	pool_msg <- MessageContext{
		Msg: fmt.Sprintf("thread [%d] payload successfully sent: %s", th_id, req.Header.String()),
	}
}
