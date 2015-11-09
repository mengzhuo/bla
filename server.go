// Package bla provides ...
package bla

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

var (
	Cfg        *Config
	configPath = flag.String("c", "./config.json", "")
)

func LoadConfig(path string) (err error) {

	d, err := ioutil.ReadFile(path)
	LFatal(err)
	return json.Unmarshal(d, &Cfg)
}

func New() {

	flag.Parse()
	LFatal(LoadConfig(*configPath))

	server := &Server{docs: make([]*Doc, 0),
		Tags: make(map[string][]*Doc, 0)}

	LFatal(server.LoadTempalte())
	server.LoadAllDocs()
	server.MakeHome()
	Log(Cfg)

}

func (s *Server) LoadTempalte() (err error) {

	root := filepath.Join(Cfg.TemplatePath, "root.tmpl")
	parse := func(name string) (*template.Template, error) {
		t := template.New("")
		return t.ParseFiles(root, filepath.Join(Cfg.TemplatePath, name))
	}

	s.template.home, err = parse("home.tmpl")
	if err != nil {
		return
	}
	return
}

type Server struct {
	docs     []*Doc
	Tags     map[string][]*Doc
	template struct{ home *template.Template }
}

type rootData struct {
	Cfg  *Config
	Docs []*Doc
}

func (s *Server) MakeHome() {
	r := &rootData{Cfg: Cfg}
	n := Cfg.HomeArticles
	if len(s.docs) < Cfg.HomeArticles {
		n = len(s.docs)
	}
	r.Docs = s.docs[0:n]

	f, err := os.OpenFile(filepath.Join(Cfg.PublicPath, "index.html"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		LErr(err)
		return
	}
	LErr(s.template.home.ExecuteTemplate(f, "", r))
}
