package bla

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Doc struct {
	Path       string
	Title      string
	headerNode *goquery.Selection
	Time       time.Time
	Tags       []string
	Parsed     *goquery.Document
	Related    []*Doc
}

func santiSpace(s string) string {
	return strings.Replace(s, " ", "-", -1)
}

func (d *Doc) String() string {
	return fmt.Sprintf("Path:%s Title:%s Time:%s Tags:%s", d.Path, d.Title, d.Time, d.Tags)
}

func (s *Server) LoadDoc(path string, info os.FileInfo, e error) (err error) {

	if !info.Mode().IsRegular() {
		return
	}

	if strings.HasPrefix(filepath.Base(path), ".") {
		return
	}

	Log(path)

	f, err := os.Open(path)
	if err != nil {
		LErr(err)
		return
	}
	defer f.Close()

	doc := &Doc{}
	parsed, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		LErr(err)
		return
	}

	header := parsed.Find("h1").First()
	doc.Title = header.Text()
	doc.headerNode = header
	doc.Path = Cfg.BasePath + santiSpace(doc.Title)

	if t := parsed.Find(".date").First().Text(); t != "" {
		doc.Time, err = time.Parse("2006-01-02", t)
		if err != nil {
			LErr(err)
			return
		}
	} else {
		doc.Time = info.ModTime()
	}

	parsed.Find(".tag").Each(func(i int, sel *goquery.Selection) {
		t := sel.Text()
		doc.Tags = append(doc.Tags, t)
		if _, ok := s.Tags[t]; ok {
			s.Tags[t] = append(s.Tags[t], doc)
		} else {
			s.Tags[t] = []*Doc{doc}
		}
	})
	// Make header links
	header.ReplaceWithHtml(fmt.Sprintf(`<h1 class="title"><a href="%s">%s</a></h1>`, doc.Path, doc.Title))

	Log(doc)
	doc.Parsed = parsed
	doc.Related = make([]*Doc, 0)
	s.docs = append(s.docs, doc)
	return
}

func (s *Server) makeRelated(d *Doc, t string) {

	for _, rd := range s.Tags[t] {
		if !docInDocs(rd, d.Related) {
			d.Related = append(d.Related, rd)
		}

	}

}

func (s *Server) LoadAllDocs() (err error) {

	filepath.Walk(Cfg.ContentPath, s.LoadDoc)
	sort.Sort(docsByTime(s.docs))

	// make related doc

	for _, d := range s.docs {
		for _, t := range d.Tags {
			s.makeRelated(d, t)
		}
	}

	return
}

type rootData struct {
	Cfg    *Config
	Server *Server
	Docs   []*Doc
}

func (s *Server) newRootData() *rootData {
	return &rootData{Server: s, Cfg: Cfg}
}

func docInDocs(doc *Doc, docs []*Doc) bool {

	for _, d := range docs {
		if d.Title == doc.Title {
			return true
		}
	}
	return false
}

func (s *Server) MakeDocAppend(doc *Doc) (app string) {

	buf := bytes.NewBuffer([]byte{})

	for _, d := range doc.Related {

		fmt.Fprintf(buf, `<a href="%s" class="related">%s</a> `, d.Path, d.Title)
		doc.Related = append(doc.Related, d)
	}
	return buf.String()
}

func (s *Server) SaveAllDocs() (err error) {
	for _, d := range s.docs {

		r := s.newRootData()
		r.Docs = []*Doc{d}

		f, err := os.Create(fmt.Sprintf("%s/%s", Cfg.PublicPath, santiSpace(d.Title)))
		LErr(err)
		LErr(s.template.doc.ExecuteTemplate(f, "root", r))
		f.Close()
	}

	//
	return
}

func (s *Server) MakeHome() (err error) {

	n := Cfg.HomeArticles
	if len(s.docs) < Cfg.HomeArticles {
		n = len(s.docs)
	}
	r := s.newRootData()

	r.Docs = s.docs[0:n]

	f, err := os.Create(fmt.Sprintf("%s/%s", Cfg.PublicPath, "index.html"))
	LErr(s.template.home.ExecuteTemplate(f, "root", r))
	return
}

type docsByTime []*Doc

func (s docsByTime) Len() int           { return len(s) }
func (s docsByTime) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s docsByTime) Less(i, j int) bool { return s[i].Time.After(s[j].Time) }
