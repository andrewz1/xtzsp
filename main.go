package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/andrewz1/xtzsp/xrecv"
)

var (
	sig = []os.Signal{
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGINT,
		syscall.SIGHUP,
	}
)

func sigSetup() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, sig...)
	return c
}

func sigWait(c chan os.Signal) os.Signal {
	for s := range c {
		for _, ss := range sig {
			if s == ss {
				return s
			}
		}
	}
	return syscall.SIGQUIT
}

func main() {
	c := sigSetup()
	if err := xrecv.Init(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(sigWait(c))
}
