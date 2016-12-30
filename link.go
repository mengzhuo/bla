package bla

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func wrapLinkToPublic(s *Handler, public string) filepath.WalkFunc {

	return func(path string, info os.FileInfo, err error) error {
		if path == s.Cfg.RootPath {
			return nil
		}

		if strings.Count(path, "/")-strings.Count(s.Cfg.RootPath, "/") > 1 {
			return nil
		}

		switch base := filepath.Base(path); base {
		case "template", "docs", ".public", "certs":
			return nil
		default:
			realPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			target := filepath.Join(public, base)
			log.Printf("link %s -> %s", realPath, target)

			err = os.Symlink(realPath, target)
			if err != nil {
				if os.IsExist(err) {
					return nil
				}
				log.Fatal(err)
			}
		}
		return nil
	}
}
