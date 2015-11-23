// Package bla provides ...
package bla

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
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

	admin := map[string]http.HandlerFunc{
		".add":  server.Add,
		".edit": server.Edit,
		".quit": server.Quit,
	}
	for k, v := range admin {
		http.HandleFunc(path.Join(Cfg.BasePath, k), server.auth(v))
	}

	http.Handle("/", http.FileServer(http.Dir(Cfg.PublicPath)))

	if Cfg.TLSKeyFile != "" && Cfg.TLSCertFile != "" {
		log.Fatal(http.ListenAndServeTLS(Cfg.Addr, Cfg.TLSCertFile, Cfg.TLSKeyFile, nil))
	} else {
		log.Fatal(http.ListenAndServe(Cfg.Addr, nil))
	}
}

func (s *Server) auth(orig http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if u, p, ok := r.BasicAuth(); ok {
			if u == Cfg.Username && p == Cfg.Password {
				orig(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="`+Cfg.BaseURL+`"`)
		w.WriteHeader(401)
		w.Write([]byte("401 Unauthorized\n"))
	}
}

func (s *Server) Quit(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Location", strings.Replace(Cfg.BaseURL, "://", "://log:out@", 1))
	w.WriteHeader(301)
}

func (s *Server) Add(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		rd := s.newRootData()
		if err := s.template.add.ExecuteTemplate(w, "root", rd); err != nil {
			log.Print(err)
			w.Write([]byte(err.Error()))
			w.WriteHeader(500)
		}
	case "POST":
		r.ParseForm()
		content := r.FormValue("mce_0")
		parsed, err := goquery.NewDocumentFromReader(bytes.NewBuffer([]byte(content)))
		if err != nil {
			log.Printf("Parse failed:%v", err)
			w.WriteHeader(500)
			return
		}
		header := parsed.Find("h1").First().Text()
		if header == "" {
			log.Print("Header not existed", header)
			w.WriteHeader(500)
			return

		}

		f, err := os.Create(filepath.Join(Cfg.ContentPath, santiSpace(header)+".html"))
		if err != nil {
			log.Print(err)
			w.WriteHeader(500)
			return
		}

		f.WriteString(content)
		f.Close()

		w.Header().Set("Location", Cfg.BasePath)
		w.WriteHeader(301)

	default:
		w.WriteHeader(403)

	}
}

func (s *Server) Edit(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		doc := r.URL.Query().Get("doc")
		if doc == "" {

			w.Write([]byte("no doc query"))
			w.WriteHeader(403)
			return
		}

		fp := filepath.Join(Cfg.ContentPath, doc+".html")
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			w.Write([]byte("no such doc"))
			w.WriteHeader(404)
			return
		}

		f, err := os.Open(fp)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(500)
			return
		}
		defer f.Close()

		rd := s.newRootData()
		parsed, _ := goquery.NewDocumentFromReader(f)

		rd.Docs = append(rd.Docs, &Doc{Parsed: parsed})
		if err := s.template.add.ExecuteTemplate(w, "root", rd); err != nil {
			log.Print(err)
			w.Write([]byte(err.Error()))
			w.WriteHeader(500)
			f.Close()
			return
		}

		bak, err := os.Create(fp + ".bak")
		io.Copy(bak, f)
		bak.Close()
		f.Close()
	}
}

func (s *Server) Reset() {

	if err := os.RemoveAll(Cfg.PublicPath); err != nil {
		log.Print(err)
	}
	os.MkdirAll(Cfg.PublicPath, 0755)

	if Cfg.Favicon != "" {
		Link(Cfg.Favicon, filepath.Join(Cfg.PublicPath, "favicon.ico"))
	}

	for _, p := range Cfg.LinkFiles {
		Link(p, filepath.Join(Cfg.PublicPath, filepath.Base(p)))
	}

	uploads, _ := filepath.Abs(Cfg.UploadPath)
	uploadBase := filepath.Base(uploads)

	os.MkdirAll(filepath.Join(Cfg.PublicPath, Cfg.BasePath), 0755)
	Link(Cfg.UploadPath, filepath.Join(Cfg.PublicPath, Cfg.BasePath, uploadBase))
}

func (s *Server) LoadStaticFiles() {

	log.Print("Rebuild from: ", Cfg.TemplatePath)
	Link(filepath.Join(Cfg.TemplatePath, "asset"),
		filepath.Join(Cfg.PublicPath, Cfg.BasePath, "asset"))

	log.Print("Loading Asset files")

	orig, err := filepath.Abs(filepath.Join(Cfg.TemplatePath, "asset"))
	if err != nil {
		log.Fatal(err)
	}

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
		"now":  time.Now,
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

	s.template.add = parsed("add.tmpl")
	return
}

type Server struct {
	Docs        map[string]*Doc
	sortedDocs  []*Doc
	Tags        map[string][]*Doc
	StaticFiles []string
	template    struct{ home, doc, index, add *template.Template }
	Version     string
	StopWatch   chan bool
	DocLock     *sync.RWMutex
	TagLock     *sync.RWMutex
}

func Link(orig, dest string) error {

	absOrig, err := filepath.Abs(orig)
	if err != nil {
		return err
	}
	absDest, err := filepath.Abs(dest)
	if err != nil {
		return err
	}

	if err := os.Symlink(absOrig, absDest); err != nil {
		log.Print(err)
		return err
	} else {
		log.Print("Rebuild link:", absOrig, " -> ", absDest)
	}
	return nil
}
