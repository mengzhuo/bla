package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/mengzhuo/bla"
)

func main() {

	cd, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config *bla.Config
	json.Unmarshal(cd, &config)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Config:%+v", config)
	bla.New(config).Run()
}
