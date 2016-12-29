package bla

import (
	"os"
	"path/filepath"
	"sort"
)

func generateTagPage(s *Handler) (err error) {

	defer ErrOrOk("generate tags", err)

	for tagName, docs := range s.tags {
		err = makeTagPage(s, tagName, docs)
		if err != nil {
			return
		}
	}
	return nil
}

func makeTagPage(s *Handler, tagName string, docs []*Doc) (err error) {
	var f *os.File

	f, err = os.Create(filepath.Join(s.publicPath, "/tags/", tagName))
	if err != nil {
		return
	}
	defer f.Close()
	sort.Sort(docsByTime(docs))
	err = s.tpl.ExecuteTemplate(f, "tag_page",
		&tagData{s, tagName, docs, tagName})
	return
}
