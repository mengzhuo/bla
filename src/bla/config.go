package bla

type Config struct {
	BaseURL  string
	HostName string

	RootPath string
	LinkPath []string

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

		RootPath: "root",
		LinkPath: defaultLinks,

		HomeDocCount: 5,
		Title:        "笔记本",
	}
}
