package bla

import (
	"os"
	"path/filepath"
)

func generateAllPage(s *Handler) (err error) {

	f, err := os.Create(filepath.Join(s.publicPath, "all"))
	if err != nil {
		return
	}
	defer f.Close()

	return s.tpl.ExecuteTemplate(f, "all",
		&mulDocData{s, "", s.sortDocs})
}
