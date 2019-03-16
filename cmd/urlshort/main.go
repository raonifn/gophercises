package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

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

func loadYml(handler http.Handler) (http.Handler, error) {
	if ymlFilename == "" {
		return handler, nil
	}

	yml, err := ioutil.ReadFile(ymlFilename)
	if err != nil {
		return nil, err
	}

	return urlshort.YAMLHandler(yml, handler)
}

func loadJson(handler http.Handler) (http.Handler, error) {
	if jsonFilename == "" {
		return handler, nil
	}

	json, err := ioutil.ReadFile(jsonFilename)
	if err != nil {
		return nil, err
	}

	return urlshort.JSONHandler(json, handler)
}

func main() {
	flags()
	mux := defaultMux()

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	handler := urlshort.MapHandler(pathsToUrls, mux)

	handler, err := loadYml(handler)
	if err != nil {
		panic(err)
	}

	handler, err = loadJson(handler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
