package bla

import (
	"strings"
	"testing"
)

func TestNewRawDoc(t *testing.T) {
	ff := `Title=Hello World
	Time=2016-03-20T06:12:44Z
Tags=golang, bla, epoch
Public=true
+++
### Hohohoho
	`
	doc, err := newDoc(strings.NewReader(ff))
	if err != nil {
		t.Error(err)
	}
	t.Log(doc)
	t.Log(string(doc.Content))
}
