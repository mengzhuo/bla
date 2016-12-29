package bla

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func (h *Handler) generateSiteMap() (err error) {

	buf := bytes.NewBuffer(nil)
	for _, d := range h.sortDocs {
		fmt.Fprintf(buf, "%s%s/%s", "https://", h.Cfg.HostName,
			path.Join(h.Cfg.BaseURL, d.SlugTitle))
	}

	err = ioutil.WriteFile(filepath.Join(h.publicPath, "sitemap.txt"),
		buf.Bytes(), os.ModePerm)
	ErrOrOk("gernate site map", err)
	return
}
