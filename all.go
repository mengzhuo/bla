package bla

import (
	"os"
	"path/filepath"
)

func generateAllPage(s *Handler, publicPath string) (err error) {

	f, err := os.Create(filepath.Join(publicPath, "all"))
	if err != nil {
		return
	}
	defer f.Close()

	return s.tpl.ExecuteTemplate(f, "all",
		&mulDocData{s, "", s.sortDocs})
}
