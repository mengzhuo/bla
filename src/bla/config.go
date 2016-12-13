package bla

import (
	"log"
	"net/http"
	"os"
)

type Config struct {
	BaseURL string

	DocPath   string
	AssetPath string
}

func DefaultConfig() *Config {

	return &Config{
		BaseURL:   "/",
		DocPath:   "doc",
		AssetPath: "asset",
	}
}

type InitConfigSite struct {
	cfgPath string
	valid   bool
}

func (s *InitConfigSite) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if s.valid {

	}

	f, err := os.Open(s.cfgPath)
	if err != nil {
		log.Print(err)
	}

}

func NewInitConfigHandler(configPath string) http.Handler {
	return &InitConfigSite{configPath}
}
