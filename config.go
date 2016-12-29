package bla

import (
	"fmt"
	"log"
	"os"
	"os/user"
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
