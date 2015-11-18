package bla

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Doc struct {
	RealPath   string // as id
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

	if info == nil || !info.Mode().IsRegular() {
		return
	}

	if !strings.HasSuffix(filepath.Base(path), ".html") {
		return
	}

	//Log(path)

	f, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return
	}
	defer f.Close()

	doc := &Doc{RealPath: path,
		Related: make([]*Doc, 0)}

	parsed, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Printf("Parse failed:%v", err)
		return
	}

	header := parsed.Find("h1").First()
	doc.Title = header.Text()
	doc.headerNode = header
	doc.Path = Cfg.BasePath + santiSpace(doc.Title)

	if t := parsed.Find(".date").First().Text(); t != "" {
		doc.Time, err = time.Parse("2006-01-02", t)
		if err != nil {
			log.Print(err)
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

	//Log(doc)
	doc.Parsed = parsed

	s.DocLock.Lock()
	defer s.DocLock.Unlock()

	s.Docs[doc.RealPath] = doc
	return
}

func (s *Server) makeRelated() {

	// clean all docs
	s.sortedDocs = s.sortedDocs[:0]

	for _, d := range s.Docs {
		// reset related docs
		d.Related = d.Related[:0]
		s.sortedDocs = append(s.sortedDocs, d)
		for _, t := range d.Tags {
			for _, rd := range s.Tags[t] {
				if !docInDocs(rd, d.Related) && rd != d {
					d.Related = append(d.Related, rd)
				}
			}
		}
	}

	sort.Sort(docsByTime(s.sortedDocs))

}

func (s *Server) LoadAllDocs() (err error) {

	start := time.Now()
	s.Tags = make(map[string][]*Doc)

	filepath.Walk(Cfg.ContentPath, s.LoadDoc)

	// make related doc
	s.makeRelated()
	log.Printf("Load %d docs in %s", len(s.sortedDocs), time.Since(start))

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

func (s *Server) SaveAllDocs() (err error) {

	r := s.newRootData()

	err = os.MkdirAll(path.Join(Cfg.PublicPath, Cfg.BasePath), 0755)

	if err != nil {
		log.Fatal(err)
	}
	for _, d := range s.Docs {

		r.Docs = []*Doc{d}
		docPubPath := path.Join(Cfg.PublicPath, d.Path)
		if docPubPath == path.Dir(docPubPath) {
			continue
		}

		f, err := os.Create(docPubPath)
		if err != nil {
			log.Print(err)
		}

		if err = s.template.doc.ExecuteTemplate(f, "root", r); err != nil {
			log.Print(err)
		}
		f.Close()
		os.Chtimes(docPubPath, d.Time, d.Time)
	}

	return
}

func (s *Server) MakeHome() (err error) {

	log.Print("making home of ", len(s.sortedDocs))
	n := Cfg.HomeArticles
	if len(s.sortedDocs) < Cfg.HomeArticles {
		n = len(s.sortedDocs)
	}
	r := s.newRootData()

	r.Docs = s.sortedDocs[0:n]

	f, err := os.Create(path.Join(Cfg.PublicPath, Cfg.BasePath, "index.html"))
	if err != nil {
		log.Print(err)
	}
	if err = s.template.home.ExecuteTemplate(f, "root", r); err != nil {
		log.Print(err)
	}

	j := 1
	for i := 0; i < len(s.Docs); i += Cfg.HomeArticles {
		f, err := os.Create(path.Join(Cfg.PublicPath, Cfg.BasePath, fmt.Sprintf("index-%d", j)))
		if err != nil {
			log.Print(err)
		}

		n := i + Cfg.HomeArticles
		if i+Cfg.HomeArticles > len(s.sortedDocs) {
			n = len(s.sortedDocs) - 1
		}

		r.Docs = s.sortedDocs[i:n]

		if err = s.template.index.ExecuteTemplate(f, "root", r); err != nil {
			log.Print(err)
		}
		j++
	}

	return
}

type docsByTime []*Doc

func (s docsByTime) Len() int           { return len(s) }
func (s docsByTime) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s docsByTime) Less(i, j int) bool { return s[i].Time.After(s[j].Time) }
