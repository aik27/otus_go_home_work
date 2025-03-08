package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go-telnet [--timeout=10s] host port")
	}

	address := net.JoinHostPort(flag.Arg(0), flag.Arg(1))
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		fmt.Fprintln(os.Stderr, "Connection error:", err)
		os.Exit(1)
	}

	defer func() {
		err := client.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Connection close error:", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		err := client.Send()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Send error:", err)
		}
	}()

	go func() {
		err := client.Receive()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Receive error:", err)
		}
	}()

	<-ctx.Done()
}
