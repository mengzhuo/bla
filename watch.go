package bla

import (
	"os"
	"path"

	"gopkg.in/fsnotify.v1"
)

func (s *Server) RebuildDoc(e fsnotify.Event) {

	Log("rebuilding doc:", e)
	if e.Op != fsnotify.Chmod {

		p := path.Clean(e.Name)

		Log(p)

		if e.Op == fsnotify.Write {
			f, err := os.Open(p)
			LErr(err)
			info, err := f.Stat()
			LErr(err)
			s.LoadDoc(p, info, nil)
		}
	}

	s.MakeHome()
}

func (s *Server) ReloadConfig() {

	s.StopWatch <- true

}

func (s *Server) Watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		LFatal(err)
	}

	LFatal(watcher.Add(Cfg.ContentPath))
	LFatal(watcher.Add(Cfg.TemplatePath))

	go func(watcher *fsnotify.Watcher) {
		defer watcher.Close()

		contentPath := path.Clean(Cfg.ContentPath)
		tmplPath := path.Clean(Cfg.TemplatePath)
		Log("CP: ", contentPath)

		for {
			select {
			case e := <-watcher.Events:

				Log("EP: ", path.Clean(e.Name))

				switch path.Dir(path.Clean(e.Name)) {
				case tmplPath:
				case contentPath:
					s.RebuildDoc(e)
				}
			case err := <-watcher.Errors:
				LErr(err)
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
		LFatal(err)
	}
	defer watcher.Close()

	LFatal(watcher.Add(*configPath))

}
