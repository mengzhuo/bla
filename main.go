package bla

import (
	"io/ioutil"
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

	ini "gopkg.in/ini.v1"

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
	loadWebDav(h)

	return h
}

func (s *Handler) watch() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		// init all docs
		loadData(s)
		loadTemplate(s)
		err := s.saveAll()
		if err != nil {
			log.Fatal("can't save docs:", err)
		}
		// loadData minial interval is 1 second
		ticker := time.NewTicker(time.Second)
		docChange := false
		rootChange := false

		for {

			select {
			case event := <-watcher.Events:
				switch ext := filepath.Ext(event.Name); ext {
				case ".md", ".tmpl":
					log.Println("modified file:", event.Name)
					docChange = true
				case ".swp":
					continue
				}

				rootChange = true
			case err := <-watcher.Errors:
				log.Println("error:", err)
			case <-ticker.C:
				if docChange {
					docChange = false
					loadData(s)
					loadTemplate(s)
				}

				if rootChange {
					rootChange = false
					err := s.saveAll()
					if err != nil {
						log.Print("can't save docs:", err)
						continue
					}
				}
			}
		}
	}()

	watcher.Add(s.Cfg.RootPath)
	watcher.Add(s.docPath)
	watcher.Add(s.templatePath)

	if err != nil {
		log.Print(err)
	}
}

func clearOldTmp(exclude string) (err error) {
	realExcluded, err := filepath.Abs(exclude)
	if err != nil {
		return err
	}
	old, err := filepath.Glob(filepath.Join(os.TempDir(), "bla*"))
	if err != nil {
		return err
	}

	for _, path := range old {
		realPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		if realPath == realExcluded {
			continue
		}
		log.Println("removing old public", realPath)
		os.RemoveAll(realPath)
	}
	return nil
}

type handleFunc func(s *Handler) error

func (s *Handler) saveAll() (err error) {

	s.publicPath, err = ioutil.TempDir("", "bla")
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	err = os.MkdirAll(s.publicPath, 0700)
	if err != nil {
		return
	}

	for k, function := range map[string]handleFunc{
		"index":    generateIndex,
		"all page": generateAllPage,
		"sitemap":  generateSiteMap,
		"tag":      generateTagPage} {

		err = function(s)
		if err != nil {
			log.Printf("generate:%-15s ... [ERR]", k)
			return err
		}
		log.Printf("generate:%-15s ... [OK]", k)
	}

	log.Printf("linking all dir in %s", s.Cfg.RootPath)
	filepath.Walk(s.Cfg.RootPath, s.linkToPublic)
	log.Println("save completed")

	s.public = http.FileServer(http.Dir(s.publicPath))
	clearOldTmp(s.publicPath)
	return nil
}

func loadData(s *Handler) {
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

	rawCfg, err := ini.Load(h.cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("loading config")
	cfg := DefaultConfig()

	rawCfg.MapTo(cfg)

	h.templatePath = filepath.Join(cfg.RootPath, "template")
	h.docPath = filepath.Join(cfg.RootPath, "docs")

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

func loadTemplate(s *Handler) {
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
