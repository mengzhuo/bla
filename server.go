// Package bla provides ...
package bla

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
	"text/template"
)

const Version = "0.1 alpha"

func New() {

	flag.Parse()
	LoadConfig(*configPath)

	server := &Server{
		Docs:        make(map[string]*Doc, 0),
		sortedDocs:  make([]*Doc, 0),
		DocLock:     &sync.RWMutex{},
		Tags:        make(map[string][]*Doc, 0),
		TagLock:     &sync.RWMutex{},
		StaticFiles: make([]string, 0),
		Version:     Version,
		StopWatch:   make(chan bool)}

	server.Reset()
	if err := server.LoadTempalte(); err != nil {
		log.Fatal(err)
	}
	server.LoadStaticFiles()
	server.LoadAllDocs()
	if err := server.SaveAllDocs(); err != nil {
		log.Fatal(err)
	}
	server.MakeHome()
	log.Print(Cfg)
	server.Watch()
	server.ConfigWatch()
	http.HandleFunc(path.Join(Cfg.BasePath, "/.login"), Login)
	http.Handle("/", http.FileServer(http.Dir(Cfg.PublicPath)))
	http.ListenAndServe(Cfg.Addr, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	log.Print(r.BasicAuth())
}

func (s *Server) Reset() {

	if err := os.RemoveAll(Cfg.PublicPath); err != nil {
		log.Print(err)
	}
	os.MkdirAll(Cfg.PublicPath, 0755)

	if Cfg.Favicon != "" {
		fav, _ := filepath.Abs(Cfg.Favicon)
		dest, _ := filepath.Abs(filepath.Join(Cfg.PublicPath, "favicon.ico"))

		if err := os.Symlink(fav, dest); err != nil {
			log.Print(err)
		} else {
			log.Print("Rebuild link:", fav, " -> ", dest)
		}
	}

	uploads, _ := filepath.Abs(Cfg.UploadPath)
	dest, _ := filepath.Abs(filepath.Join(Cfg.PublicPath, Cfg.BasePath, "uploads"))

	os.MkdirAll(filepath.Join(Cfg.PublicPath, Cfg.BasePath), 0755)

	if err := os.Symlink(uploads, dest); err != nil {
		log.Print(err)
	} else {
		log.Print("Rebuild link:", uploads, " -> ", dest)
	}
}

func (s *Server) LoadStaticFiles() {

	log.Print("Rebuild from: ", Cfg.TemplatePath)
	orig, err := filepath.Abs(filepath.Join(Cfg.TemplatePath, "asset"))
	dest, err := filepath.Abs(filepath.Join(Cfg.PublicPath, "asset"))
	if err != nil {
		log.Print(err)
	}

	log.Print("Rebuild link:", orig, " -> ", dest)
	if err := os.Symlink(orig, dest); err != nil {
		log.Print(err)
	}
	log.Print("Loading Asset files")
	filepath.Walk(orig,
		func(path string, info os.FileInfo, e error) (err error) {

			if info == nil || !info.Mode().IsRegular() {
				return
			}
			if !(filepath.Ext(path) == ".css" || filepath.Ext(path) == ".js") {
				return
			}
			log.Print("|- ", path)
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
			log.Fatal(err)
		}
		return t
	}

	s.template.home = parsed("home.tmpl")
	s.template.doc = parsed("doc.tmpl")
	s.template.index = parsed("index.tmpl")
	return
}

type Server struct {
	Docs        map[string]*Doc
	sortedDocs  []*Doc
	Tags        map[string][]*Doc
	StaticFiles []string
	template    struct{ home, doc, index *template.Template }
	Version     string
	StopWatch   chan bool
	DocLock     *sync.RWMutex
	TagLock     *sync.RWMutex
}
