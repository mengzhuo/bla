package bla

type Config struct {
	BaseURL  string
	HostName string

	DocPath      string
	StaticPath   string
	TemplatePath string

	HomeDocCount int
	Title        string
}

func DefaultConfig() *Config {

	return &Config{
		BaseURL:      "",
		DocPath:      "docs",
		StaticPath:   "",
		TemplatePath: "template",

		HomeDocCount: 5,
		Title:        "笔记本",
	}
}
