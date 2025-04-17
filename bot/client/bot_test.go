package client_test

import (
	"fmt"
	"log"
	"os"
	"path"
	"testing"
	"z3ntl3/storm/bot/client"
	"z3ntl3/storm/cli"
)

func TestClient(t *testing.T) {
	os.Args = append(os.Args, "--target=test", "--timeout=5s")
	if err := os.Chdir(path.Join(os.Getenv("HOME"), "Documents", "storm")); err != nil {
		log.Fatal(err)
	}

	cli.Execute()
	client := client.New()

	for item := range client.Rsrc.Proxies.Iter() {
		fmt.Printf("proxies: %+v\n", *item)
	}

	fmt.Printf("reset: %s\n", *client.Rsrc.Proxies.Next())

	for item := range client.Rsrc.Headers.Iter() {
		fmt.Printf("headers: %+v\n", *item)
	}

	for item := range client.Rsrc.Accepts.Iter() {
		fmt.Printf("accepts: %+v\n", *item)
	}
}
