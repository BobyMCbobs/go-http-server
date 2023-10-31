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

func pageExistsByPath(p string) bool {
	const indexPage = "/index.html"
	f := p
	if !strings.HasPrefix(f, "/") {
		f = "/" + f
	}
	f = path.Clean(f)
	if strings.HasSuffix(p, indexPage) {
		return true // Don't actually know if it exists, but this will be redirected by the file server
	}
	stat, err := os.Stat(path.Join(common.GetServeFolder(), f))
	if err != nil {
		return false
	}
	if stat.IsDir() {
		if p[len(p)-1] != '/' {
			return true // Don't actually know if it exists, but this will be redirected by the file server
		}
	} else {
		if p[len(p)-1] == '/' {
			return true // Don't actually know if it exists, but this will be redirected by the file server
		}
	}
	if stat.IsDir() {
		if p == "" || p[len(p)-1] != '/' {
			return true // Don't actually know if it exists, but this will be redirected by the file server
		}
		// Does index.html exist in the directory?
		index := strings.TrimSuffix(f, "/") + indexPage
		stat, err := os.Stat(path.Join(common.GetServeFolder(), index))
		return err == nil && !stat.IsDir()
	}
	return true
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
		if !pageExistsByPath(req.URL.Path) {
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
