package client

import (
	// todo make own lb impl because the below one is not accomodating our needs
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

func New() *Bot {
	// cwd, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// instance := new(Bot)
	// for i, path_ := range []string{globals.Accepts, globals.Headers, globals.ProxyFile} {
	// 	f, err := os.Open(path.Join(cwd, path_))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	contents, err := io.ReadAll(f)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	contents = []byte(strings.Trim(string(contents), "\r\n"))
	// 	if i == 0 {
	// 		parts := strings.Split(string(contents), "\n")
	// 		lb, err = lb.New[string](parts...)
	// 	}
	// }
}
