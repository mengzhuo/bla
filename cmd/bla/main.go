package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mengzhuo/bla"
)

const (
	DefaultConfig = "config.ini"
)

var (
	configPath = flag.String("config", DefaultConfig, "default config path")
	varg       = flag.Bool("v", false, "show version")
	Version    = "dev"
)

func main() {
	// defer profile.Start().Stop()
	flag.Parse()

	if *varg {
		fmt.Println("bla version:", Version)
		return
	}
	go listenToUSR1()
	go bla.ListenLEHTTP(*configPath)
	bla.ListenAndServe(*configPath)
}

func listenToUSR1() {

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGUSR1)
		for range c {
			log.Print("got reload cert signal")
			bla.LoadCertificate()
		}
	}()
}
