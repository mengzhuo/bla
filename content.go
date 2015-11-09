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
	Path  string
	Title string
	Time  time.Time
	Tags  []string
	HTML  []byte
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

	doc.Title = parsed.Find("h1").First().Text()
	doc.Path = Cfg.BasePath + info.Name()

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

	f.Seek(0, 0)
	doc.HTML, err = ioutil.ReadAll(f)
	if err != nil {
		LErr(err)
		return
	}

	Log(doc)
	s.docs = append(s.docs, doc)
	return
}

func (s *Server) LoadAllDocs() (err error) {

	filepath.Walk(Cfg.ContentPath, s.LoadDoc)
	sort.Sort(docsByTime(s.docs))
	return
}

type docsByTime []*Doc

func (s docsByTime) Len() int           { return len(s) }
func (s docsByTime) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s docsByTime) Less(i, j int) bool { return s[i].Time.After(s[j].Time) }
