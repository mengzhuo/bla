package main

import (
	"flag"
	"fmt"

	"github.com/mengzhuo/bla"
)

const (
	DefaultConfig = "config.ini"
)

var (
	configPath = flag.String("config", DefaultConfig, "default config path")
	version    = flag.Bool("version", false, "show version")
)

func main() {
	// defer profile.Start().Stop()
	flag.Parse()

	if *version {
		fmt.Println("bla version:", bla.Version)
		return
	}
	bla.ListenAndServe(*configPath)
}
