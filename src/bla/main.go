package bla

import (
	"encoding/json"
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

	"golang.org/x/net/webdav"

	"github.com/fsnotify/fsnotify"
)

// Initialize
type Handler struct {
	cfgPath string
	Cfg     *Config
	public  http.Handler
	webfs   http.Handler
	tpl     *template.Template

	publicPath   string
	templatePath string
	docPath      string

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
	h.loadWebDav()

	return h
}

func (s *Handler) loadWebDav() {

	fs := webdav.Dir(s.Cfg.RootPath)
	ls := webdav.NewMemLS()

	handler := &webdav.Handler{
		Prefix:     "/fs",
		FileSystem: fs,
		LockSystem: ls,
	}
	s.webfs = handler
}

func (s *Handler) watch() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		// loadData minial interval is 1 second
		ticker := time.NewTicker(500 * time.Millisecond)
		mod := true

		for {
			select {
			case event := <-watcher.Events:
				switch ext := filepath.Ext(event.Name); ext {
				case ".md", ".json", ".tmpl":
					log.Println("modified file:", event.Name)
					mod = true
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
				return
			case <-ticker.C:
				if !mod {
					continue
				}
				mod = false
				s.loadData()
				s.loadTemplate()
				err := s.saveAll()
				if err != nil {
					log.Fatal("can't save docs:", err)
				}
			}
		}
	}()

	watcher.Add(s.docPath)
	watcher.Add(s.templatePath)
	watcher.Add(s.cfgPath)

	if err != nil {
		log.Print(err)
	}

}

func (s *Handler) saveAll() (err error) {

	log.Print("Saving all docs...")
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Println("Deleting public...")
	err = os.RemoveAll(s.publicPath)
	if err != nil {
		return
	}

	var f *os.File
	err = os.MkdirAll(s.publicPath, 0755)
	if err != nil {
		return
	}

	log.Printf("saving all docs...")
	for _, doc := range s.docs {
		f, err = os.Create(filepath.Join(s.publicPath, doc.SlugTitle))
		if err != nil {
			return
		}
		if err = s.tpl.ExecuteTemplate(f, "single",
			&singleData{s, doc.Title, doc}); err != nil {
			return
		}
		f.Close()
	}

	var docs []*Doc
	if len(s.sortDocs) > s.Cfg.HomeDocCount {
		docs = s.sortDocs[:s.Cfg.HomeDocCount]
	} else {
		docs = s.sortDocs
	}

	f, err = os.Create(filepath.Join(s.publicPath, "index.html"))
	if err != nil {
		return
	}
	defer f.Close()

	if err = s.tpl.ExecuteTemplate(f, "index",
		&mulDocData{s, "", docs}); err != nil {
		return
	}

	log.Printf("saving all tags...")
	err = os.MkdirAll(filepath.Join(s.publicPath, "tags"), 0755)
	if err != nil {
		return
	}

	for tagName, docs := range s.tags {
		f, err = os.Create(filepath.Join(s.publicPath, "/tags/", tagName))
		if err != nil {
			return
		}
		sort.Sort(docsByTime(docs))
		if err = s.tpl.ExecuteTemplate(f, "tag_page",
			&tagData{s, tagName, docs, tagName}); err != nil {
			return
		}
		f.Close()
	}

	log.Printf("linking...%v", s.Cfg.LinkPath)
	for _, p := range s.Cfg.LinkPath {
		var realPath string
		realPath, err = filepath.Abs(p)
		if err != nil {
			return
		}
		err = os.Symlink(realPath, filepath.Join(s.publicPath, p))
		if err != nil {
			return
		}
	}
	log.Println("save completed")
	return nil
}

func (s *Handler) loadData() {
	log.Print("Loading docs from:", s.docPath)

	s.mu.Lock()
	s.sortDocs = []*Doc{}
	s.docs = map[string]*Doc{}
	s.tags = map[string][]*Doc{}
	s.mu.Unlock()

	f, err := os.Open(s.docPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = filepath.Walk(s.docPath, s.docWalker)
	if err != nil {
		log.Print(err)
	}
	sort.Sort(docsByTime(s.sortDocs))
	log.Print("End Loading docs from:", s.docPath)
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

	h.publicPath = filepath.Join(cfg.RootPath, "public")
	h.templatePath = filepath.Join(cfg.RootPath, "template")
	h.docPath = filepath.Join(cfg.RootPath, "docs")

	h.public = http.FileServer(http.Dir(h.publicPath))

	h.Cfg = cfg
	log.Printf("%#v", *cfg)

}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if strings.HasPrefix(r.URL.Path, "/fs") {
		s.webfs.ServeHTTP(w, r)
	} else {
		s.public.ServeHTTP(w, r)
	}
}

func (s *Handler) loadTemplate() {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Printf("loding template:%s", s.templatePath)
	tpl, err := template.ParseGlob(s.templatePath + "/*.tmpl")
	if err != nil {
		log.Print(err)
		return
	}
	s.tpl = tpl
}
