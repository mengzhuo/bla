package bla

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func generateSiteMap(h *Handler, public string) (err error) {

	buf := bytes.NewBuffer(nil)
	for _, d := range h.sortDocs {
		fmt.Fprintf(buf, "%s%s/%s\n", "https://", h.Cfg.HostName,
			path.Join(h.Cfg.BaseURL, d.SlugTitle))
	}

	err = ioutil.WriteFile(filepath.Join(public, "sitemap.txt"),
		buf.Bytes(), os.ModePerm)
	return
}
