package bla

import (
	"os"
	"path/filepath"
)

func generateIndex(s *Handler, public string) (err error) {

	var docs []*Doc
	if len(s.sortDocs) > s.Cfg.HomeDocCount {
		docs = s.sortDocs[:s.Cfg.HomeDocCount]
	} else {
		docs = s.sortDocs
	}

	f, err := os.Create(filepath.Join(public, "index.html"))
	if err != nil {
		return
	}
	defer f.Close()

	return s.tpl.ExecuteTemplate(f, "index",
		&mulDocData{s, "", docs})

}

func generateDoc(s *Handler, public string) (err error) {

	for _, doc := range s.docs {
		var f *os.File
		f, err = os.Create(filepath.Join(public, doc.SlugTitle))
		if err != nil {
			return
		}
		defer f.Close()
		err = s.tpl.ExecuteTemplate(f, "single",
			&singleData{s, doc.Title, doc})
		if err != nil {
			return
		}
	}
	return nil
}
