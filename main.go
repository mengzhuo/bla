package bla

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/fsnotify/fsnotify"
)

const (
	Version = "dev"
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

	loadConfig(h)
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
		err := saveAll(s)
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
					if strings.HasPrefix(event.Name, ".") {
						continue
					}
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
					err := saveAll(s)
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
	old, err := filepath.Glob(filepath.Join(os.TempDir(), "bla_*"))
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

type handleFunc func(s *Handler, public string) error

func saveAll(s *Handler) (err error) {

	var newPub string
	newPub, err = ioutil.TempDir("", fmt.Sprintf("bla_%s_", s.Cfg.UserName))
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	for k, function := range map[string]handleFunc{
		"index":    generateIndex,
		"all_page": generateAllPage,
		"sitemap":  generateSiteMap,
		"tag":      generateTagPage,
		"docs":     generateSingle,
	} {
		start := time.Now()
		err = function(s, newPub)
		cost := time.Now().Sub(start)
		if err != nil {
			log.Printf("generate:%-15s ... [ERR]", k)
			return err
		}
		log.Printf("generate:%-15s ... [OK][%s]", k, cost)
	}

	filepath.Walk(s.Cfg.RootPath, wrapLinkToPublic(s, newPub))

	s.publicPath = newPub
	s.public = http.FileServer(http.Dir(s.publicPath))
	clearOldTmp(s.publicPath)
	return nil
}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if strings.HasPrefix(r.URL.Path, "/fs") {
		s.webfs.ServeHTTP(w, r)
		httpRequestCount.WithLabelValues("/fs").Inc()
	} else {
		s.public.ServeHTTP(w, r)
		httpRequestCount.WithLabelValues("/doc").Inc()
	}
}
