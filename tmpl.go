package bla

import (
	"log"
	"text/template"
)

func loadTemplate(s *Handler) {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Printf("loding template:%s", s.templatePath)
	tpl, err := template.ParseGlob(s.templatePath + "/*.tmpl")
	if err != nil {
		log.Print(err)
		return
	}
	s.tpl = tpl
}

type mulDocData struct {
	Hdl   *Handler
	Title string
	Docs  []*Doc
}

type singleData struct {
	Hdl   *Handler
	Title string
	Doc   *Doc
}

type tagData struct {
	Hdl     *Handler
	Title   string
	Docs    []*Doc
	TagName string
}
