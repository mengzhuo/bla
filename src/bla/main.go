package bla

import (
	"net/http"
	"path/filepath"
	"strings"
)

type Server struct {
	cfg    *Config
	static http.Handler
}

func buildSite(cfg *Config) {
	filepath.Walk(cfg.ContentPath, s.loadDocAndPage)
}

func (s *Site) loadDocAndPage(p string, info *os.Info, e error) error {
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
