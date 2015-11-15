// Package bla provides ...
package bla

import (
	"flag"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

const Version = "0.1 alpha"

func New() {

	flag.Parse()
	LoadConfig(*configPath)

	server := &Server{
		Docs:        make(map[string]*Doc, 0),
		Tags:        make(map[string][]*Doc, 0),
		StaticFiles: make([]string, 0),
		Version:     Version,
		StopWatch:   make(chan bool)}

	LFatal(server.LoadTempalte())
	server.LoadStaticFiles()
	server.LoadAllDocs()
	LFatal(server.SaveAllDocs())
	server.MakeHome()
	Log(Cfg)
	server.Watch()
	http.Handle("/", http.FileServer(http.Dir(Cfg.PublicPath)))
	http.ListenAndServe(Cfg.Addr, nil)
}

func (s *Server) LoadStaticFiles() {

	Log("Rebuild from: ", Cfg.TemplatePath)
	LErr(os.RemoveAll(Cfg.PublicPath))
	LErr(os.MkdirAll(Cfg.PublicPath, 0755))
	orig, err := filepath.Abs(filepath.Join(Cfg.TemplatePath, "asset"))
	dest, err := filepath.Abs(filepath.Join(Cfg.PublicPath, "asset"))
	LErr(err)

	Log("Rebuild link:", orig, " -> ", dest)
	LErr(os.Symlink(orig, dest))
	Log("Loading Asset files")
	filepath.Walk(orig,
		func(path string, info os.FileInfo, e error) (err error) {

			if info == nil || !info.Mode().IsRegular() {
				return
			}
			if !(filepath.Ext(path) == ".css" || filepath.Ext(path) == ".js") {
				return
			}
			Log("|- ", path)
			s.StaticFiles = append(s.StaticFiles, path)
			return
		})
}

func (s *Server) LoadTempalte() (err error) {

	var funcMap = template.FuncMap{
		"base": path.Base,
		"ext":  filepath.Ext,
	}
	root := filepath.Join(Cfg.TemplatePath, "root.tmpl")
	parsed := func(fname string) (t *template.Template) {
		t, err = template.New("").Funcs(funcMap).ParseFiles(root, filepath.Join(Cfg.TemplatePath, fname))
		if err != nil {
			LFatal(err)
		}
		return t
	}

	s.template.home = parsed("home.tmpl")
	s.template.doc = parsed("doc.tmpl")
	return
}

type Server struct {
	Docs        map[string]*Doc
	Tags        map[string][]*Doc
	StaticFiles []string
	template    struct{ home, doc *template.Template }
	Version     string
	StopWatch   chan bool
}
