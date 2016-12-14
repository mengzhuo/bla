package bla

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
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

type Page struct {
	Title   string
	Content string
}

type docsByTime []*Doc

func (s docsByTime) Len() int           { return len(s) }
func (s docsByTime) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s docsByTime) Less(i, j int) bool { return s[i].Time.After(s[j].Time) }
