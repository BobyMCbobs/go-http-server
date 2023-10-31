package handlers

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/NYTimes/gziphandler"
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

	headerMap := map[string][]string{}
	err = nil
	if common.GetHeaderSetEnable() == "true" {
		tplHeaderMapPath := common.GetHeaderMapPath()
		headerMap, err = common.LoadHeaderMapConfig(tplHeaderMapPath)
		if err != nil {
			panic(err)
		}
		headerMap = common.EvaluateEnvFromHeaderMap(headerMap)
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
		if common.GetHeaderSetEnable() == "true" {
			w = common.WriteHeadersToResponse(w, headerMap)
		}
		if err := tmpl.Execute(w, htmlTemplateOptions); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

// serveHandlerStandard ...
// handles sending the serve folder
func serveHandlerStandard(publicDir string) http.Handler {
	handler := http.FileServer(http.Dir(publicDir))

	headerMap := map[string][]string{}
	var err error = nil
	if common.GetHeaderSetEnable() == "true" {
		tplHeaderMapPath := common.GetHeaderMapPath()
		headerMap, err = common.LoadHeaderMapConfig(tplHeaderMapPath)
		if err != nil {
			panic(err)
		}
		headerMap = common.EvaluateEnvFromHeaderMap(headerMap)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if common.GetHeaderSetEnable() == "true" {
			w = common.WriteHeadersToResponse(w, headerMap)
		}
		if _, err := os.Stat(path.Join(common.GetServeFolder(), path.Clean(req.URL.Path+"/"))); err != nil {
			w.WriteHeader(404)
			http.ServeFile(w, req, path.Join(common.GetServeFolder(), common.Get404PageFileName()))
			return
		}
		handler.ServeHTTP(w, req)
	})
}

// ServeHandler ...
// serves a folder
func ServeHandler(publicDir string) (handler http.Handler) {
	if common.GetVuejsHistoryMode() == "true" {
		handler = serveHandlerVuejsHistoryMode(publicDir)
	} else {
		handler = serveHandlerStandard(publicDir)
	}
	if common.GetEnableGZIP() {
		handler = gziphandler.GzipHandler(handler)
	}
	return handler
}
