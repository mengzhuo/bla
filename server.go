package bla

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Cfg struct {
	DBPath  string
	Address string
}

type Server struct {
	cfg    *Cfg
	dbPath string
	mux    *http.ServeMux
}

func NewServer(cfg *Cfg) (server *Server) {
	return &Server{cfg, "",
		http.NewServeMux()}
}

func (s *Server) Run() (err error) {
	// Log
	log.SetPrefix("[BLA]")
	s.InitDir(s.cfg.DBPath)
	return
}

func (s *Server) InitDir(path string) (err error) {

	s.dbPath = filepath.Clean(path)
	log.Printf("Entries path:%s", s.dbPath)
	err = os.MkdirAll(s.dbPath, 0700)
	if os.IsExist(err) {
		return nil
	}
	return
}

func (s *Server) LoadEntries() (err error) {

	return
}

func handleEntries(w http.ResponseWriter, r *http.Request) {
	io.Copy(s.header, w)
	io.Copy(entry, w)
	io.Copy(s.footer, w)
}
