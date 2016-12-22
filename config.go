package bla

import (
	"fmt"
	"os"
	"os/user"
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
		"libs", "asset",
	}
	name := user.Current().Name
	return &Config{
		BaseURL:  "",
		HostName: os.Hostname(),

		RootPath: "root",
		LinkPath: defaultLinks,

		HomeDocCount: 5,
		Title:        fmt.Sprintf("%s's blog", name),
		UserName:     name,
		Password:     "PLEASE_UPDATE_PASSWORD!",
	}
}
