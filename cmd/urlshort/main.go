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

func loadDefault() (urlshort.HandlerStacker, error) {
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	return urlshort.MapHandler(pathsToUrls)
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

func loadDB() (urlshort.HandlerStacker, error) {
	return urlshort.DBHandler("my.db")
}

func loadJSON() (urlshort.HandlerStacker, error) {
	if jsonFilename == "" {
		return nil, nil
	}

	json, err := ioutil.ReadFile(jsonFilename)
	if err != nil {
		return nil, err
	}

	return urlshort.JSONHandler(json)
}

func load(server *urlshort.Server, loaders ...func() (urlshort.HandlerStacker, error)) {
	for _, l := range loaders {
		hs, err := l()
		if err != nil {
			fmt.Printf("Error loading handler: %v\n", err)
		}
		if hs != nil {
			server.StackHandler(hs)
		}
	}
}

func main() {
	flags()
	server := urlshort.NewServer()

	load(server, loadDefault, loadYml, loadJSON, loadDB)

	fmt.Println("Starting the server on :8080")
	server.Start(":8080")
}
