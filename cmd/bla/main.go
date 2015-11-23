package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/mengzhuo/bla"
)

const (
	DefaultConfig  = "bla.config.json"
	DefaultContent = "bla_content"
)

func main() {

	flag.Parse()
	if len(flag.Args()) > 0 && flag.Args()[0] == "new" {
		New()
		FetchDefaultTemplate()
		CreateDefaultContent()
	} else {
		bla.New()
	}

}

func New() {

	if _, err := os.Stat(DefaultConfig); !os.IsNotExist(err) {
		log.Fatalf("Already has default config:%s", DefaultConfig)
	}

	cfg := &bla.Config{
		Addr:     ":8080",
		BaseURL:  "http://localhost:8080",
		BasePath: "/",

		Username: "admin",
		Password: "",

		HomeArticles: 10,
		Title:        "Some blah",
		Footer:       "My daily blog",
		DocURLFormat: "%s",

		ContentPath:  DefaultContent,
		UploadPath:   "./bla_uploads",
		TemplatePath: "./bla_default_template",

		PublicPath: "./bla_public",
		Favicon:    "",

		TLSCertFile: "",
		TLSKeyFile:  "",
	}

	f, err := os.Create("bla.config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	b, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	f.Write(b)
	log.Printf("new config :\n%s", b)
	return

}

func FetchDefaultTemplate() {
	cmd := exec.Command("git", "clone", "--depth=1", "https://github.com/mengzhuo/bla_default_template.git")
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	cmd.Wait()
}

func CreateDefaultContent() {

	if _, err := os.Stat(DefaultContent); !os.IsNotExist(err) {
		log.Fatalf("Already has default config:%s", DefaultContent)
	}

	err := os.MkdirAll(DefaultContent, 0755)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(path.Join(DefaultContent, "hello.html"))
	if err != nil {
		log.Fatal(err)
	}
	f.WriteString(`<h1>Hello from old time</h1><p class="date">1970-01-01</p><p> Blah!</p>`)
	defer f.Close()
}
