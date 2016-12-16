package bla

type Config struct {
	BaseURL  string
	HostName string

	DocPath         string
	ExternalLibPath string
	TemplatePath    string

	HomeDocCount int
	Title        string
}

func DefaultConfig() *Config {

	return &Config{
		BaseURL:  "",
		HostName: "meng.zhuo.blog",

		DocPath:         "docs",
		ExternalLibPath: "libs",
		TemplatePath:    "template",

		HomeDocCount: 5,
		Title:        "笔记本",
	}
}
