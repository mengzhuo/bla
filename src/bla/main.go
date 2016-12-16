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
	"sync"
	"text/template"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Initialize
type Handler struct {
	cfgPath string
	Cfg     *Config
	extLib  http.Handler
	public  http.Handler
	tpl     *template.Template

	mu       sync.RWMutex
	docs     map[string]*Doc
	sortDocs []*Doc
	tags     map[string][]*Doc
}

func NewHandler(cfgPath string) *Handler {

	h := &Handler{
		cfgPath: cfgPath,
		mu:      sync.RWMutex{},
	}

	h.loadConfig()
	h.watch()

	return h
}

func (s *Handler) watch() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		// loadData minial interval is 1 second
		ticker := time.NewTicker(time.Second)
		mod := false
		s.loadData()
		s.loadTemplate()
		for {
			select {
			case event := <-watcher.Events:
				switch ext := filepath.Ext(event.Name); ext {
				case ".md":
					fallthrough
				case ".tmpl":
					log.Println("modified file:", event.Name)
					mod = true
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
				return
			case <-ticker.C:
				if mod {
					mod = false
					s.loadData()
					s.loadTemplate()
				}
			}
		}
	}()

	watcher.Add(s.Cfg.DocPath)
	watcher.Add(s.Cfg.TemplatePath)

	if err != nil {
		log.Print(err)
	}

}

func (s *Handler) loadData() {
	log.Print("Loading docs from:", s.Cfg.DocPath)

	s.mu.Lock()
	s.sortDocs = []*Doc{}
	s.docs = map[string]*Doc{}
	s.tags = map[string][]*Doc{}
	s.mu.Unlock()

	f, err := os.Open(s.Cfg.DocPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = filepath.Walk(s.Cfg.DocPath, s.docWalker)
	if err != nil {
		log.Print(err)
	}
	sort.Sort(docsByTime(s.sortDocs))
	log.Print("End Loading docs from:", s.Cfg.DocPath)
}

func (s *Handler) docWalker(p string, info os.FileInfo, err error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

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

	h.extLib = http.StripPrefix("/libs", http.FileServer(http.Dir(cfg.ExternalLibPath)))
	h.public = http.FileServer(http.Dir(cfg.DocPath))

	h.Cfg = cfg
	log.Printf("%#v", *cfg)

}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	/*
		if r.Host != s.Cfg.HostName {
			w.Header().Add("Location", "https://"+s.Cfg.HostName)
			w.WriteHeader(301)
			return
		}
	*/

	cnt := strings.Count(r.URL.Path, "/")
	if cnt == 1 {
		switch r.URL.Path {
		case "/":
			s.ServeHome(w, r)
		case "/robots.txt":
		case "/favicon.ico":
		default:
			s.ServeDoc(w, r)
		}
	} else {
		if r.URL.Path[:5] == "/libs" {
			s.extLib.ServeHTTP(w, r)
		} else {
			s.public.ServeHTTP(w, r)
		}
	}

	duration := time.Now().Sub(start)
	w.Header().Add("Response-In", duration.String())
	log.Printf("%s %s %s", r.Method, r.URL.Path, duration.String())
}

func (s *Handler) ServeDoc(w http.ResponseWriter, r *http.Request) {

	docName := strings.TrimLeft(r.URL.Path, "/")
	s.mu.RLock()
	defer s.mu.RUnlock()
	doc, ok := s.docs[docName]
	if !ok {
		http.NotFound(w, r)
		return
	}

	if err := s.tpl.ExecuteTemplate(w, "single", &rootData{s, nil, doc}); err != nil {
		Error(err, w, r)
	}
}

func Error(err error, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "internal error:%s", err)
	log.Printf("ERR:%s %s", err, r.URL.Path)
}

func (s *Handler) ServeHome(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	docs := s.sortDocs
	if len(s.sortDocs) > s.Cfg.HomeDocCount {
		docs = docs[:s.Cfg.HomeDocCount]
	}
	s.mu.RUnlock()

	if err := s.tpl.ExecuteTemplate(w, "index", &rootData{s, docs, nil}); err != nil {
		Error(err, w, r)
	}
}

func (s *Handler) loadTemplate() {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Printf("loding template:%s", s.Cfg.TemplatePath)
	tpl, err := template.ParseGlob(s.Cfg.TemplatePath + "/*.tmpl")
	if err != nil {
		log.Print(err)
		return
	}
	s.tpl = tpl
}

type rootData struct {
	Hdl  *Handler
	Docs []*Doc
	Doc  *Doc
}
