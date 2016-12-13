package bla

type Config struct {
	BaseURL string

	DocPath      string
	AssetPath    string
	HomeDocCount int
}

func DefaultConfig() *Config {

	return &Config{
		BaseURL:   "",
		DocPath:   "doc",
		AssetPath: "",

		HomeDocCount: 5,
	}
}
