package bla

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
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
}

func (d *Doc) String() string {
	return fmt.Sprintf("Path:%s Title:%s Time:%s Tags:%s", d.Path, d.Title, d.Time, d.Tags)
}

func (s *Server) LoadDoc(path string, info os.FileInfo, e error) (err error) {

	if !info.Mode().IsRegular() {
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
	doc.Path = Cfg.BasePath + doc.Title

	if t := parsed.Find("#datetime").First().Text(); t != "" {
		doc.Time, err = time.Parse("2012-01-02T10:04:05", t)
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

	Log(doc)
	doc.Parsed = parsed

	s.docs = append(s.docs, doc)
	return
}

func (s *Server) LoadAllDocs() (err error) {

	filepath.Walk(Cfg.ContentPath, s.LoadDoc)
	sort.Sort(docsByTime(s.docs))
	return
}

type rootData struct {
	Title   string
	Content string
	Append  string
	Cfg     *Config
}

func (s *Server) SaveAllDocs() (err error) {
	for _, d := range s.docs {

		// Make header links
		d.headerNode.ReplaceWithHtml(fmt.Sprintf(`<h1 class="title"><a href="%s">%s</a></h1>`, d.Path, d.Title))

		// get HTML
		html, err := d.Parsed.First().Html()
		if err != nil {
			LErr(err)
			continue
		}
		r := rootData{Title: d.Title + "-", Content: html, Cfg: Cfg}

		f, err := os.Create(fmt.Sprintf("%s/%s", Cfg.PublicPath, d.Title))
		LErr(s.template.ExecuteTemplate(f, "doc", r))
		f.Close()
	}

	//
	return
}

func (s *Server) MakeHome() (err error) {

	html := ""
	for _, d := range s.docs {
		p, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", Cfg.PublicPath, d.Title))
		LErr(err)
		html += string(p)
	}

	r := rootData{Content: html, Cfg: Cfg}

	f, err := os.Create(fmt.Sprintf("%s/%s", Cfg.PublicPath, "index.html"))
	LErr(s.template.ExecuteTemplate(f, "doc", r))
	return
}

type docsByTime []*Doc

func (s docsByTime) Len() int           { return len(s) }
func (s docsByTime) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s docsByTime) Less(i, j int) bool { return s[i].Time.After(s[j].Time) }
