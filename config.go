// Package bla provides ...
package bla

type Config struct {
	Addr     string `json:"address"`
	BaseURL  string `json:"base_url"`
	BaseRoot string `json:"base_root"`
	BasePath string `json:"base_path"`

	username string `json:"username"`
	password string `json:"password"`

	HomeArticles int    `json:"home_articles"`
	Title        string `json:"title"`

	ContentPath  string `json:"content_path"`
	TemplatePath string `json:"template_path"`

	PublicPath string `json:"public_path"` // all parsed html/content etc...
}
