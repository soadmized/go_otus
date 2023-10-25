package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout")
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	go waitForSigint(cancel)

	args := flag.Args()
	if len(args) != 2 {
		log.Printf("usage: go-telnet --timeout host port")
	}

	addr := fmt.Sprintf("%s:%s", args[0], args[1])
	cl := NewTelnetClient(addr, timeout, os.Stdin, os.Stdout)

	err := cl.Connect()
	if err != nil {
		log.Printf("connect error: %v", err)
		return
	}

	defer func() {
		err := cl.Close()
		if err != nil {
			log.Printf("close error: %v", err)
		}
	}()

	go func() {
		err := cl.Receive()
		if err != nil {
			log.Printf("receive error: %v", err)
		}

		cancel()
	}()

	go func() {
		err := cl.Send()
		if err != nil {
			log.Printf("send error: %v", err)
		}

		cancel()
	}()

	<-ctx.Done()
}

func waitForSigint(cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	cancel()
}
