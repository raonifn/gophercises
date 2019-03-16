package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

type redirect struct {
	URL  string `json:"url", yaml:"url"`
	Path string `json:"path", yaml:"path"`
}

type mapper struct {
	config   map[string]string
	fallback http.Handler
}

type unmarshaller func([]byte, interface{}) error

func (m *mapper) handle(w http.ResponseWriter, r *http.Request) {
	url, ok := m.config[r.URL.Path]
	if !ok {
		m.fallback.ServeHTTP(w, r)
		return
	}
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func resultToConfig(rs []redirect) map[string]string {
	ret := map[string]string{}
	for _, r := range rs {
		ret[r.Path] = r.URL
	}
	return ret
}

func resultHandler(content []byte, fallback http.Handler, unmarshal unmarshaller) (http.Handler, error) {
	rs := []redirect{}
	err := unmarshal(content, &rs)

	if err != nil {
		return nil, err
	}

	config := resultToConfig(rs)
	return MapHandler(config, fallback), nil

}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.Handler {
	m := mapper{config: pathsToUrls, fallback: fallback}
	return http.HandlerFunc(m.handle)
}

func JSONHandler(content []byte, fallback http.Handler) (http.Handler, error) {
	return resultHandler(content, fallback, json.Unmarshal)
}

func YAMLHandler(content []byte, fallback http.Handler) (http.Handler, error) {
	return resultHandler(content, fallback, yaml.Unmarshal)
}
