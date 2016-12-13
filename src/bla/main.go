package bla

import (
	"net/http"
	"strings"
)

type Site struct {
	cfg    *Config
	static http.Handler
	docs   map[string]*Doc
}

func NewSite(cfg string) *http.ServeMux {
	return &Site{cfg: cfg}
}

func (s *Site) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch p := strings.TrimPrefix(r.URL.Path, s.cfg.BaseURL); p {
	case "/":
	case "/index":
	case "/search":
	case "/upload":
	case "/edit":
	case "/admin":
	default:
		if _, ok := s.docs[p]; ok {

		}
		s.static.ServeHTTP(w, r)
		return
	}
}
