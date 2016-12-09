package bla

import (
	"net/http"
	"strings"
)

type Server struct {
	cfg    *Config
	static http.Handler
}

func buildSite(cfg *Config) {
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch p := strings.TrimPrefix(r.URL.Path, s.cfg.BaseURL); p {
	case "/":
	case "/index":
	case "/search":
	case "/stats":
	case "/upload":
	case "/edit":
	case "/admin":
	default:
		s.static.ServeHTTP(w, r)
		return
	}
}
