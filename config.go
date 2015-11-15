// Package bla provides ...
package bla

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"path/filepath"
)

var (
	Cfg        *Config
	configPath = flag.String("c", "./config.json", "")
)

type Config struct {
	Addr     string `json:"address"`
	BaseURL  string `json:"base_url"`
	BaseRoot string `json:"base_root"`
	BasePath string `json:"base_path"`

	username string `json:"username"`
	password string `json:"password"`

	HomeArticles int    `json:"home_articles"`
	Title        string `json:"title"`
	Footer       string `json:"footer"`

	ContentPath  string `json:"content_path"`
	TemplatePath string `json:"template_path"`

	PublicPath string `json:"public_path"` // all parsed html/content etc...
}

func LoadConfig(path string) (err error) {

	d, err := ioutil.ReadFile(path)
	LFatal(err)
	err = json.Unmarshal(d, &Cfg)

	// Clean
	Cfg.ContentPath = filepath.Clean(Cfg.ContentPath)
	Cfg.TemplatePath = filepath.Clean(Cfg.TemplatePath)
	Cfg.PublicPath = filepath.Clean(Cfg.PublicPath)
	return
}
