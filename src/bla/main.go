package bla

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Initialize
type Handler struct {
	cfgPath string
	cfg     *Config
	static  http.Handler

	docs     map[string]*Doc
	sortDocs []*Doc
	tags     map[string][]*Doc
}

func NewHandler(cfgPath string) *Handler {

	h := &Handler{cfgPath: cfgPath}
	h.loadConfig()
	h.loadDoc()
	return h
}

func (s *Handler) loadDoc() {
	s.sortDocs = []*Doc{}
	s.docs = map[string]*Doc{}
}

func (h *Handler) loadConfig() {

	f, err := os.Open(h.cfgPath)
	if err != nil && os.IsExist(err) {
		log.Panic(err)
	}

	log.Print("loading config")
	cfg := DefaultConfig()
	dec := json.NewDecoder(f)
	err = dec.Decode(cfg)

	if err != nil {
		log.Panic(err)
	}

	h.static = http.FileServer(http.Dir(cfg.AssetPath))
	h.cfg = cfg
	log.Printf("%s", cfg)

}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	switch p := strings.TrimPrefix(r.URL.Path, s.cfg.BaseURL); p {
	case "/":
		s.ServeHome(w, r)
	case "/tag":
		s.ServeTag(w, r)
	default:
		if _, ok := s.docs[p]; ok {

		} else {
			s.static.ServeHTTP(w, r)
		}
	}

	duration := time.Now().Sub(start)
	w.Header().Add("Response-In", duration.String())
	log.Printf("%s %s %s", r.Method, r.URL.Path, duration.String())

}

func Error(err error, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "internal error:%s", err)
	log.Printf("ERR:%s %s", err, r.URL.Path)
}

func (s *Handler) ServeHome(w http.ResponseWriter, r *http.Request) {

	if len(s.sortDocs) < s.cfg.HomeDocCount {
		s.sortDocs
	}

}
