package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/raonifn/gophercises/internal/urlshort"
)

var (
	ymlFilename  string
	jsonFilename string
)

func flags() {
	flag.StringVar(&ymlFilename, "ymlConfig", "", "YML to config redirects.")
	flag.StringVar(&jsonFilename, "jsonConfig", "", "YML to config redirects.")
	flag.Parse()
}

func loadYml() (urlshort.HandlerStacker, error) {
	if ymlFilename == "" {
		return nil, nil
	}

	yml, err := ioutil.ReadFile(ymlFilename)
	if err != nil {
		return nil, err
	}

	return urlshort.YAMLHandler(yml)
}

func loadJson() (urlshort.HandlerStacker, error) {
	if jsonFilename == "" {
		return nil, nil
	}

	json, err := ioutil.ReadFile(jsonFilename)
	if err != nil {
		return nil, err
	}

	return urlshort.JSONHandler(json)
}

func main() {
	flags()
	server := urlshort.NewServer()

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	hs := urlshort.MapHandler(pathsToUrls)
	server.StackHandler(hs)

	hs, err := loadYml()
	if err != nil {
		panic(err)
	}
	if hs != nil {
		server.StackHandler(hs)
	}

	hs, err = loadJson()
	if err != nil {
		panic(err)
	}
	if hs != nil {
		server.StackHandler(hs)
	}

	fmt.Println("Starting the server on :8080")
	server.Start(":8080")
}
