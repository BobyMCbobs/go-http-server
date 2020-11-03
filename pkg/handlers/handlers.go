package handlers

import (
	"html/template"
	"net/http"
	"path"
	"strings"
	"gitlab.com/safesurfer/go-http-server/pkg/common"
)

// serveHandlerVuejsHistoryMode ...
// handles sending the serve folder with Vuejs history mode
func serveHandlerVuejsHistoryMode(publicDir string) http.Handler {
	handler := http.FileServer(http.Dir(publicDir))

	tplMapPath := common.GetTemplateMapPath()
	configMap, err := common.LoadMapConfig(tplMapPath)
	if err != nil {
		panic(err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// static files
		if strings.Contains(req.URL.Path, ".") {
			handler.ServeHTTP(w, req)
			return
		}

		// frontend views
		indexPath := path.Join(publicDir, "/index.html")
		tmpl, err := template.ParseFiles(indexPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		htmlTemplateOptions := common.EvaluateEnvFromMap(configMap)
		if err := tmpl.Execute(w, htmlTemplateOptions); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

// serveHandlerStandard ...
// handles sending the serve folder
func serveHandlerStandard(publicDir string) http.Handler {
	handler := http.FileServer(http.Dir(publicDir))

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(w, req)
	})
}

// ServeHandler ...
// serves a folder
func ServeHandler(publicDir string) http.Handler {
	if common.GetVuejsHistoryMode() == "true" {
		return serveHandlerVuejsHistoryMode(publicDir)
	}
	return serveHandlerStandard(publicDir)
}
