package bla

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/gin-gonic/gin"

	"golang.org/x/tools/present"
)

type Doc struct {
	*present.Doc
	Permalink    string
	HTML         template.HTML
	Related      []*Doc
	Newer, Older *Doc
}

type Server struct {
	cfg      *Config
	router   *gin.Engine
	log      *log.Logger
	template struct{ index *template.Template }
	docs     []*Doc
	docPaths map[string]*Doc
	docTags  map[string][]*Doc
}

func New(config *Config) (server *Server) {

	router := gin.Default()
	//router.Group("/admin", gin.BasicAuth(gin.Accounts{config.Username: config.Password}))
	router.Static("/", config.PublicPath)

	return &Server{cfg: config,
		router:   router,
		log:      log.New(os.Stderr, "[Bla] ", log.LstdFlags),
		docs:     make([]*Doc, 0),
		docPaths: make(map[string]*Doc),
		docTags:  make(map[string][]*Doc)}
}

func (s *Server) Log(err error) {
	if err != nil {
		s.log.Print(err)
	}
}

func (s *Server) LoadTemplate() (err error) {

	root := filepath.Join(s.cfg.TemplatePath, "root.tmpl")
	parse := func(name string) (t *template.Template, err error) {
		t = template.New("")
		return t.ParseFiles(root, filepath.Join(s.cfg.TemplatePath, name))
	}

	s.template.index, err = parse("index.tmpl")
	return
}

func (s *Server) Run() {
	s.LoadDocs()
	s.Log(s.LoadTemplate())
	s.Log(s.MakeIndex())
	s.log.Fatal(s.router.Run(s.cfg.Addr))
}

func (s *Server) LoadDocs() (err error) {

	const ext = ".article"
	startAt := time.Now()

	fn := func(p string, info os.FileInfo, e error) (err error) {
		if filepath.Ext(p) != ext {
			return nil
		}

		f, err := os.Open(p)
		if err != nil {
			return
		}
		defer f.Close()

		pd, err := present.Parse(f, p, 0)
		if err != nil {
			return
		}
		tp := filepath.ToSlash(info.Name()[0 : len(info.Name())-len(ext)])

		d := &Doc{
			Doc:       pd,
			Permalink: s.cfg.BaseURL + tp,
		}
		s.docs = append(s.docs, d)

		for _, tag := range d.Tags {
			if _, ok := s.docTags[tag]; ok {
				s.docTags[tag] = append(s.docTags[tag], d)
			} else {
				s.docTags[tag] = []*Doc{d}
			}
		}
		return nil
	}

	err = filepath.Walk(s.cfg.ContentPath, fn)
	if err != nil {
		s.log.Print(err)
		return
	}
	// sort by Time
	sort.Sort(docsByTime(s.docs))

	for i, d := range s.docs {
		if i != 0 {
			d.Newer = s.docs[i-1]
		}
		if i+1 < len(s.docs) {
			d.Older = s.docs[i+1]
		}

		for _, tag := range d.Tags {
			for _, rd := range s.docTags[tag] {
				if rd != d {
					d.Related = append(d.Related, rd)
				}
			}
		}
	}
	s.log.Printf("Loaded %d docs in %s", len(s.docs), time.Since(startAt))
	return
}

type rootData struct {
	BasePath string
	Doc      *Doc
	Data     interface{}
}

func (s *Server) MakeIndex() (err error) {

	buf := bytes.NewBuffer([]byte{})
	d := &rootData{BasePath: s.cfg.BasePath}

	n := s.cfg.HomeArticles
	if s.cfg.HomeArticles >= len(s.docs) {
		n = len(s.docs)
	}
	d.Data = s.docs[0:n]
	s.template.index.ExecuteTemplate(buf, "index", d)

	err = ioutil.WriteFile(filepath.Join(s.cfg.PublicPath, "index.html"), buf.Bytes(), 0666)
	return
}

func (s *Server) MakeHTML() {

}

// docsByTime implements sort.Interface, sorting Docs by their Time field.
type docsByTime []*Doc

func (s docsByTime) Len() int           { return len(s) }
func (s docsByTime) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s docsByTime) Less(i, j int) bool { return s[i].Time.After(s[j].Time) }
