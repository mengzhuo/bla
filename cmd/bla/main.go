// Package bla provides ...
package main

import (
	"flag"
	"log"

	"github.com/mengzhuo/bla"
)

var (
	dbPath  = flag.String("path", "./entries", "")
	address = flag.String("addr", "127.0.0.1:8080", "")
)

func main() {
	flag.Parse()
	server := bla.NewServer(&bla.Cfg{DBPath: *dbPath,
		Address: *address})
	log.Fatal(server.Run())
}
