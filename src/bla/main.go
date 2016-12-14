package bla

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Initialize
type Handler struct {
	cfgPath string
	cfg     *Config
	static  http.Handler
	tpl     *template.Template

	docs     map[string]*Doc
	sortDocs []*Doc
	tags     map[string][]*Doc
}

func NewHandler(cfgPath string) *Handler {

	h := &Handler{cfgPath: cfgPath}

	h.loadConfig()
	h.loadTemplate()
	h.loadData()
	h.watch()

	return h
}

func (s *Handler) watch() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					s.loadData()
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
				return
			}
		}
	}()

	watcher.Add(s.cfg.DocPath)

	log.Print("watching: ", s.cfg.DocPath)
	if err != nil {
		log.Print(err)
	}

}

func (s *Handler) loadData() {
	s.sortDocs = []*Doc{}
	s.docs = map[string]*Doc{}
	s.tags = map[string][]*Doc{}

	f, err := os.Open(s.cfg.DocPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	log.Print("Loading docs from:", s.cfg.DocPath)
	err = filepath.Walk(s.cfg.DocPath, s.docWalker)
	if err != nil {
		log.Print(err)
	}
	sort.Sort(docsByTime(s.sortDocs))
}

func (s *Handler) docWalker(p string, info os.FileInfo, err error) error {

	start := time.Now()
	if info.IsDir() || filepath.Ext(info.Name()) != ".md" {
		return nil
	}
	var f *os.File
	f, err = os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()

	var doc *Doc
	doc, err = newDoc(f)
	if err != nil {
		return err
	}
	if !doc.Public {
		log.Printf("doc:%s loaded but not public", p)
		return nil
	}

	doc.SlugTitle = path.Base(p)[0 : len(path.Base(p))-3]

	for _, t := range doc.Tags {
		s.tags[t] = append(s.tags[t], doc)
	}
	s.docs[doc.SlugTitle] = doc
	s.sortDocs = append(s.sortDocs, doc)
	log.Printf("loaded doc:%s in %s", doc.SlugTitle,
		time.Now().Sub(start).String())
	return nil
}

func (h *Handler) loadConfig() {

	f, err := os.Open(h.cfgPath)
	if err != nil && os.IsExist(err) {
		log.Panic(err)
	}
	defer f.Close()

	log.Print("loading config")
	cfg := DefaultConfig()
	dec := json.NewDecoder(f)
	err = dec.Decode(cfg)

	if err != nil {
		log.Panic(err)
	}

	h.static = http.FileServer(http.Dir(cfg.AssetPath))
	h.cfg = cfg
	log.Printf("%#v", *cfg)

}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	switch p := strings.TrimPrefix(r.URL.Path, s.cfg.BaseURL); p {
	case "/":
		s.ServeHome(w, r)
	default:
		log.Print("looking for doc:", p)
		if doc, ok := s.docs[strings.TrimLeft(p, "/")]; ok {
			s.ServeDoc(doc, w, r)
		} else {
			s.static.ServeHTTP(w, r)
		}
	}

	duration := time.Now().Sub(start)
	w.Header().Add("Response-In", duration.String())
	log.Printf("%s %s %s", r.Method, r.URL.Path, duration.String())

}

func (s *Handler) ServeDoc(doc *Doc, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, doc.Content)
}

func Error(err error, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "internal error:%s", err)
	log.Printf("ERR:%s %s", err, r.URL.Path)
}

func (s *Handler) ServeHome(w http.ResponseWriter, r *http.Request) {

	docs := s.sortDocs
	if len(s.sortDocs) > s.cfg.HomeDocCount {
		docs = docs[:s.cfg.HomeDocCount]
	}
	s.tpl.ExecuteTemplate(w, "index.tmpl", rootData{s, docs})
}

func (s *Handler) loadTemplate() {
	var err error
	s.tpl, err = template.ParseGlob(s.cfg.TemplatePath + "/*")
	log.Print(s.tpl, err)
	return
}

type rootData struct {
	h    *Handler
	docs []*Doc
}
