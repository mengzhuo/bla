package bla

type Config struct {
	BaseURL  string
	HostName string

	DocPath      string
	TemplatePath string
	PublicPath   string
	LinkPath     []string

	HomeDocCount int
	Title        string
}

func DefaultConfig() *Config {

	defaultLinks := []string{
		"libs", "asset",
	}

	return &Config{
		BaseURL:  "",
		HostName: "meng.zhuo.blog",

		DocPath:      "docs",
		PublicPath:   "public",
		TemplatePath: "template",

		LinkPath: defaultLinks,

		HomeDocCount: 5,
		Title:        "笔记本",
	}
}
