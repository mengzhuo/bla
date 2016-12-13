package main

import (
	"bla"
	"flag"
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
	bla.Init(*configPath)

	if *certfile != "" && *keyfile != "" {
		http.ListenAndServeTLS(*addr, *certfile, *keyfile, bla.Handler)
		return
	}
	http.ListenAndServe(*addr, bla.Handler)

}
