package bla

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"

	ini "gopkg.in/ini.v1"
)

// Let's encrypted web

func ListenLEHTTP(cfgPath string) {

	raw, err := ini.Load(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &ServerConfig{}
	raw.MapTo(cfg)

	if cfg.ListenLEAddr == "" || cfg.LEDir == "" {
		log.Printf("LE not configured, %s, %s", cfg.ListenLEAddr, cfg.LEDir)
		return
	}

	dir := http.FileServer(http.Dir(cfg.LEDir))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		log.Printf("LE: %s", r.URL)
		if strings.HasPrefix(r.URL.Path, "/.well-known") {
			dir.ServeHTTP(w, r)
		} else {
			w.Header().Set("Location", fmt.Sprintf("https://%s", path.Join(r.Host, r.URL.Path)))
			w.WriteHeader(http.StatusMovedPermanently)
		}
	})
	http.ListenAndServe(cfg.ListenLEAddr, nil)
}
