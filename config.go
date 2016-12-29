package bla

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
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

func loadData(s *Handler) {
	log.Print("Loading docs from:", s.docPath)

	s.mu.Lock()
	s.sortDocs = []*Doc{}
	s.docs = map[string]*Doc{}
	s.tags = map[string][]*Doc{}
	s.mu.Unlock()

	f, err := os.Open(s.docPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = filepath.Walk(s.docPath, s.docWalker)
	if err != nil {
		log.Print(err)
	}
	sort.Sort(docsByTime(s.sortDocs))
	log.Print("End Loading docs from:", s.docPath)
}
