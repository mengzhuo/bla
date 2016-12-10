package bla

import (
	"strings"
	"testing"
)

func TestNewRawDoc(t *testing.T) {
	ff := `Hello World
2016-03-20
golang, bla, epoch
draft

Hohohoho
	`
	doc, err := newRawDoc(strings.NewReader(ff))
	if err != nil {
		t.Error(err)
	}
	t.Log(doc)
	t.Log(string(doc.Content))
}
