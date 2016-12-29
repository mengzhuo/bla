package bla

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	ini "gopkg.in/ini.v1"
)

type Config struct {
	BaseURL  string
	HostName string

	RootPath string
	LinkPath []string

	HomeDocCount int
	Title        string
	UserName     string
	Password     string
}

func DefaultConfig() *Config {

	defaultLinks := []string{
		"libs",
	}
	current, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	name := current.Username
	host, _ := os.Hostname()

	return &Config{
		BaseURL:  "",
		HostName: host,

		RootPath: "root",
		LinkPath: defaultLinks,

		HomeDocCount: 5,
		Title:        fmt.Sprintf("%s's blog", strings.Title(name)),
		UserName:     name,
		Password:     "PLEASE_UPDATE_PASSWORD!",
	}
}

func loadConfig(h *Handler) {

	rawCfg, err := ini.Load(h.cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("loading config")
	cfg := DefaultConfig()

	rawCfg.MapTo(cfg)

	h.templatePath = filepath.Join(cfg.RootPath, "template")
	h.docPath = filepath.Join(cfg.RootPath, "docs")

	h.Cfg = cfg
	log.Printf("%#v", *cfg)

}
