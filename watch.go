package bla

import (
	"log"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/fsnotify.v1"
)

func (s *Server) RebuildDoc(e fsnotify.Event) {

	if e.Op == fsnotify.Chmod || e.Op == fsnotify.Create {
		return
	}

	if base := filepath.Base(e.Name); strings.HasPrefix(base, ".") || !(strings.HasSuffix(base, ".html") || strings.HasSuffix(base, ".tmpl")) {
		return
	}
	log.Print("Events on ", e)

	s.LoadAllDocs()
	s.SaveAllDocs()
	s.MakeHome()
}

func (s *Server) Watch() {

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}

	if err := watcher.Add(Cfg.ContentPath); err != nil {
		log.Fatal(err)
	}
	if err := watcher.Add(Cfg.TemplatePath); err != nil {
		log.Fatal(err)
	}

	go func(watcher *fsnotify.Watcher) {
		defer watcher.Close()

		contentPath := path.Clean(Cfg.ContentPath)
		tmplPath := path.Clean(Cfg.TemplatePath)

		for {
			select {
			case e := <-watcher.Events:

				switch path.Dir(path.Clean(e.Name)) {
				case tmplPath:
					s.LoadTempalte()
					fallthrough
				case contentPath:
					s.RebuildDoc(e)
				}
			case err := <-watcher.Errors:
				log.Print(err)
			case <-s.StopWatch:
				return
			}
		}
	}(watcher)
}

func (s *Server) ConfigWatch() {
	// This will reload the whole server
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	watcher.Add(*configPath)
	log.Printf("watch config:%s", *configPath)

	go func(watcher *fsnotify.Watcher) {

		defer watcher.Close()

		for {
			e := <-watcher.Events

			if e.Op == fsnotify.Rename || e.Op == fsnotify.Remove {
				continue
			}

			log.Print("config updated ", e)
			s.StopWatch <- true

			LoadConfig(*configPath)
			watcher.Remove(*configPath)

			s.LoadAllDocs()
			s.SaveAllDocs()
			s.MakeHome()

			go s.Watch()
			watcher.Add(*configPath)
		}

	}(watcher)
}
