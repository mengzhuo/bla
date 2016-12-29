package bla

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/russross/blackfriday"

	ini "gopkg.in/ini.v1"
)

const (
	StatusPublic = "public"
)

type Doc struct {
	Title     string
	SlugTitle string
	Time      time.Time
	Tags      []string
	Public    bool

	Content string
}

func newDoc(r io.Reader) (d *Doc, err error) {

	var buf []byte
	buf, err = ioutil.ReadAll(r)

	idx := bytes.Index(buf, []byte("+++"))
	if idx == -1 {
		return nil, fmt.Errorf("header not found")
	}

	d = &Doc{}
	sec, err := ini.Load(buf[:idx])
	if err != nil {
		return
	}

	err = sec.MapTo(d)
	if err != nil {
		return
	}

	d.Content = string(blackfriday.MarkdownCommon(buf[idx+4:]))
	return
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

type docsByTime []*Doc

func (s docsByTime) Len() int           { return len(s) }
func (s docsByTime) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s docsByTime) Less(i, j int) bool { return s[i].Time.After(s[j].Time) }
