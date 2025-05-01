package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"z3ntl3/storm/cli"
)

func main() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		<-exit

		fmt.Println("received exit signal, exiting now")
		os.Exit(0)
	}()

	cli.Execute()
}
