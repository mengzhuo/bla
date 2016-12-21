package main

import (
	"bla"
	"flag"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	DefaultConfig = "config.json"
)

var (
	certfile   = flag.String("cert", "", "cert file path")
	keyfile    = flag.String("key", "", "cert file path")
	addr       = flag.String("addr", ":8080", "listen port")
	configPath = flag.String("config", DefaultConfig, "default config path")
	logPool    = sync.Pool{New: func() interface{} { return &LogWriter{nil, 200} }}
)

func main() {

	flag.Parse()
	h := bla.NewHandler(*configPath)
	log.Print("addr: ", *addr)
	log.Print("cfg: ", *configPath)

	lh := logTimeAndStatus(h)

	if *certfile != "" && *keyfile != "" {
		log.Printf("TLS:%s, %s", *certfile, *keyfile)
		http.ListenAndServeTLS(*addr, *certfile, *keyfile, lh)
		return
	}
	http.ListenAndServe(*addr, lh)

}

func logTimeAndStatus(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		writer := logPool.Get().(*LogWriter)
		writer.ResponseWriter = w
		writer.statusCode = 200

		handler.ServeHTTP(writer, r)
		log.Printf("%s %s %s %s %d", r.RemoteAddr,
			r.Method, r.URL.Path, time.Now().Sub(start), writer.statusCode)
		logPool.Put(writer)
	})
}

type LogWriter struct {
	http.ResponseWriter
	statusCode int
}

func (l *LogWriter) WriteHeader(i int) {
	l.statusCode = i
	l.ResponseWriter.WriteHeader(i)
}
