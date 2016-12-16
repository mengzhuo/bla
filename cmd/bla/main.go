package main

import (
	"bla"
	"flag"
	"log"
	"net/http"
)

const (
	DefaultConfig = "config.json"
)

var (
	certfile   = flag.String("cert", "", "cert file path")
	keyfile    = flag.String("key", "", "cert file path")
	addr       = flag.String("addr", ":8080", "listen port")
	configPath = flag.String("config", DefaultConfig, "default config path")
)

func main() {

	flag.Parse()
	h := bla.NewHandler(*configPath)
	log.Print("addr: ", *addr)
	log.Print("cfg: ", *configPath)

	if *certfile != "" && *keyfile != "" {
		log.Printf("TLS:%s, %s", *certfile, *keyfile)
		http.ListenAndServeTLS(*addr, *certfile, *keyfile, h)
		return
	}
	http.ListenAndServe(*addr, h)

}
