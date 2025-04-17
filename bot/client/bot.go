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
	rsrc resources
}

type resources struct {
	proxies *lb.RoundRobin[string]
	accepts *lb.RoundRobin[string]
	headers *lb.RoundRobin[string]
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
			instance.rsrc.accepts = lb
			continue
		case 1:
			instance.rsrc.headers = lb
		case 2:
			instance.rsrc.proxies = lb
		}
	}

	return instance
}
