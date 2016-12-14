package bla

type Config struct {
	BaseURL string

	DocPath      string
	AssetPath    string
	TemplatePath string
	HomeDocCount int
}

func DefaultConfig() *Config {

	return &Config{
		BaseURL:      "",
		DocPath:      "docs",
		AssetPath:    "",
		TemplatePath: "template",

		HomeDocCount: 5,
	}
}
