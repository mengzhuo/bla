package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	ini "gopkg.in/ini.v1"

	"github.com/mengzhuo/bla"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	DefaultConfig = "config.ini"
)

var (
	configPath = flag.String("config", DefaultConfig, "default config path")
	version    = flag.Bool("version", false, "show version")

	logPool = sync.Pool{New: func() interface{} { return &LogWriter{nil, 200} }}
	tlsCert *tls.Certificate
	server  *http.Server
	Version string = "dev"
)

func main() {
	//defer profile.Start().Stop()
	flag.Parse()

	if *version {
		fmt.Println("bla version:", Version)
		return
	}

	raw, err := ini.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	raw.MapTo(cfg)

	log.Printf("pid:%d", os.Getpid())

	if cfg.MetricListen != "" {
		go loadMetric(cfg.MetricListen)
	}

	log.Printf("Server:%v", cfg)

	h := bla.NewHandler(*configPath)

	lh := logTimeAndStatus(h)
	server = &http.Server{Addr: cfg.Listen, Handler: lh}

	if cfg.Certfile != "" && cfg.Keyfile != "" {
		// for higher score in ssllab
		server.TLSConfig = &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP384, tls.CurveP256, tls.CurveP521},
			PreferServerCipherSuites: true,
		}
		server.TLSConfig.GetCertificate = getCertificate
		log.Printf("TLS:%s, %s", cfg.Certfile, cfg.Keyfile)
		watchReloadCert()
		log.Fatal(server.ListenAndServeTLS(cfg.Certfile, cfg.Keyfile))
	}
	server.ListenAndServe()

}

func logTimeAndStatus(handler http.Handler) http.Handler {

	var (
		writer io.Writer
		err    error
	)

	if cfg.AccessLogPath != "" {
		var file *os.File
		file, err = os.OpenFile(cfg.AccessLogPath,
			os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			log.Fatal(err)
		}
		writer = bufio.NewWriter(file)
		log.Printf("Access Log to file: %s", cfg.AccessLogPath)
	} else {
		writer = os.Stdout
	}

	accessLogger := log.New(writer, "", log.LstdFlags)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		writer := logPool.Get().(*LogWriter)
		writer.ResponseWriter = w
		writer.statusCode = 200

		if cfg.Certfile != "" {
			writer.ResponseWriter.Header().Add("Strict-Transport-Security", "max-age=31536000")
		}
		handler.ServeHTTP(writer, r)

		delta := time.Now().Sub(start)
		accessLogger.Printf("%s %s %s %s %d",
			r.RemoteAddr, r.Method, r.URL.Path,
			delta, writer.statusCode)

		httpRequestDurationSeconds.Observe(delta.Seconds())
		httpRequestCount.Inc()

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
	cert, err := tls.LoadX509KeyPair(cfg.Certfile, cfg.Keyfile)
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

// Config ---------------------

var cfg *ServerConfig

type ServerConfig struct {
	Certfile      string
	Keyfile       string
	Listen        string
	MetricListen  string
	AccessLogPath string
}

func init() {
	cfg = &ServerConfig{
		"", "", ":8080",
		"access.log", ""}
}

func loadMetric(addr string) {

	prometheus.MustRegister(httpRequestCount)
	prometheus.MustRegister(httpRequestDurationSeconds)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(addr, nil)
}

var (
	httpRequestCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "http",
		Subsystem: "request",
		Name:      "requests_count",
		Help:      "The total number of http request",
	})
	httpRequestDurationSeconds = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace: "http",
			Subsystem: "request",
			Name:      "duration_seconds",
			Help:      "The request duration distribution",
		})
)
