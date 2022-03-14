package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/antonovegorv/simple-pinger/pinger"
	"github.com/antonovegorv/simple-pinger/pinger/config"
)

const numberOfPingers = 1

func main() {
	interval := flag.Int("i", 1, "wait interval seconds between sending each packet")
	count := flag.Int("c", 0, "stop after sending count packets")
	ttl := flag.Int("t", 64, "set the IP Time to Live")
	size := flag.Int("s", 64, "size of packet in bytes")
	flag.Parse()

	var hostname string
	if hostname = flag.Arg(0); hostname == "" {
		fmt.Println("You have to provide hostname to ping with")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(numberOfPingers)

	errorsChan := make(chan error, 1)

	p := pinger.New(ctx, wg, errorsChan, config.New(
		hostname,
		*interval,
		*count,
		*ttl,
		*size,
	))
	go p.Ping()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

Loop:
	for {
		select {
		case <-termChan:
			cancel()
			break Loop
		case err := <-errorsChan:
			if err != nil {
				fmt.Println(err)
			}
			cancel()
			break Loop
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}

	wg.Wait()

	p.LogStats()
}
