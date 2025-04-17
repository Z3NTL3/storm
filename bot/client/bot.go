package client

import (
	// todo make own lb impl because the below one is not accomodating our needs

	"io"
	"log"
	"os"
	"path"
	"strings"
	"z3ntl3/storm/globals"
	lb "z3ntl3/storm/lb/robin"
)

type Bot struct {
	Rsrc resources
}

type resources struct {
	Proxies *lb.RoundRobin[string]
	Accepts *lb.RoundRobin[string]
	Headers *lb.RoundRobin[string]
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
		if i == 0 {

		}

		parts := strings.Split(string(contents), "\n")
		lb := lb.New(parts...)

		switch i {
		case 0:
			instance.Rsrc.Accepts = lb
			continue
		case 1:
			instance.Rsrc.Headers = lb
		case 2:
			instance.Rsrc.Proxies = lb
		}
	}

	return instance
}

// Dedicated to be spawned on a goroutine
// The body is thread-safe
//
// Arguments should be passed by [Bot]'s Next method on the fields as they comfort [lb.RoundRobin]
func (c *Bot) Stress(proxy, headers, accepts string) {
	// client := fasthttp.Client{}
	// req := fasthttp.AcquireRequest()
}
