package main

import (
	"crypto/tls"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/mengzhuo/bla"
)

const (
	DefaultConfig = "config.json"
)

var (
	certfile      = flag.String("cert", "", "cert file path")
	keyfile       = flag.String("key", "", "cert file path")
	addr          = flag.String("addr", ":8080", "listen port")
	configPath    = flag.String("config", DefaultConfig, "default config path")
	accessLogPath = flag.String("accesslog", "", "default access log path")

	logPool = sync.Pool{New: func() interface{} { return &LogWriter{nil, 200} }}
	tlsCert *tls.Certificate
	server  *http.Server
)

func main() {

	flag.Parse()
	log.Printf("pid:%d", os.Getpid())

	h := bla.NewHandler(*configPath)
	log.Print("addr: ", *addr)
	log.Print("cfg: ", *configPath)

	lh := logTimeAndStatus(h)
	server = &http.Server{Addr: *addr, Handler: lh}

	if *certfile != "" && *keyfile != "" {

		server.TLSConfig = &tls.Config{}
		server.TLSConfig.GetCertificate = getCertificate
		log.Printf("TLS:%s, %s", *certfile, *keyfile)
		watchReloadCert()
		log.Fatal(server.ListenAndServeTLS(*certfile, *keyfile))
	}
	server.ListenAndServe()

}

func logTimeAndStatus(handler http.Handler) http.Handler {

	var (
		writer io.Writer
		err    error
	)

	if *accessLogPath != "" {
		writer, err = os.OpenFile(*accessLogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		writer = os.Stdout
	}

	accessLogger := log.New(writer, "", log.LstdFlags)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		writer := logPool.Get().(*LogWriter)
		writer.ResponseWriter = w
		writer.statusCode = 200

		handler.ServeHTTP(writer, r)
		accessLogger.Printf("%s %s %s %s %d",
			r.RemoteAddr, r.Method, r.URL.Path,
			time.Now().Sub(start), writer.statusCode)
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

func loadCertificate() {

	log.Println("Loading new certs")
	cert, err := tls.LoadX509KeyPair(*certfile, *keyfile)
	if err != nil {
		log.Println("load cert failed keep old", err)
		return
	}
	tlsCert = &cert
}

func watchReloadCert() {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGUSR1)
		for range c {
			log.Print("got reload cert signal")
			loadCertificate()
		}
	}()
}

func getCertificate(ch *tls.ClientHelloInfo) (cert *tls.Certificate, err error) {
	return tlsCert, nil
}
