// Package bla provides ...
package bla

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
)

var (
	Cfg        *Config
	configPath = flag.String("c", "./bla.config.json", "")
)

type Config struct {
	Addr     string `json:"address"`
	BaseURL  string `json:"base_url"`
	BasePath string `json:"base_path"`

	Username string `json:"username"`
	Password string `json:"password"`

	HomeArticles int    `json:"home_articles"`
	Title        string `json:"title"`
	Footer       string `json:"footer"`
	DocURLFormat string `json:"doc_url_format"`

	ContentPath  string `json:"content_path"`
	UploadPath   string `json:"upload_path"`
	TemplatePath string `json:"template_path"`

	PublicPath    string   `json:"public_path"` // all parsed html/content etc...
	Favicon       string   `json:"favicon"`
	LinkFiles     []string `json:"link_files"`
	MaxUploadSize int64    `json:"max_upload_size"`

	TLSCertFile string `json:"certfile"`
	TLSKeyFile  string `json:"keyfile"`
}

func LoadConfig(path string) (err error) {

	d, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	Cfg = nil
	err = json.Unmarshal(d, &Cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Clean
	Cfg.ContentPath = filepath.Clean(Cfg.ContentPath)
	Cfg.TemplatePath = filepath.Clean(Cfg.TemplatePath)
	Cfg.PublicPath = filepath.Clean(Cfg.PublicPath)
	Cfg.UploadPath = filepath.Clean(Cfg.UploadPath)

	if Cfg.DocURLFormat == "" {
		Cfg.DocURLFormat = "%s"
	}

	return
}
