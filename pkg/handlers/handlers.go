package handlers

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/NYTimes/gziphandler"

	"gitlab.com/BobyMCbobs/go-http-server/pkg/common"
)

// Handler holds the information needed to create handlers
type Handler struct {
	Error404FilePath   string
	HeaderMap          map[string][]string
	GzipEnabled        bool
	HeaderMapEnabled   bool
	TemplateMap        map[string]string
	TemplateMapEnabled bool
	VueJSHistoryMode   bool
	ServeFolder        string
}

// serveHandlerVuejsHistoryMode ...
// handles sending the serve folder with Vuejs history mode
func (h *Handler) serveHandlerVuejsHistoryMode() http.Handler {
	handler := http.FileServer(http.Dir(h.ServeFolder))

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// static files
		if strings.Contains(req.URL.Path, ".") {
			handler.ServeHTTP(w, req)
			return
		}

		// frontend views
		indexPath := path.Join(h.ServeFolder, "/index.html")
		tmpl, err := template.ParseFiles(indexPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if h.HeaderMapEnabled == true {
			w = common.WriteHeadersToResponse(w, h.HeaderMap)
		}
		if err := tmpl.Execute(w, h.TemplateMap); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

// serveHandlerStandard ...
// handles sending the serve folder
func (h *Handler) serveHandlerStandard() http.Handler {
	handler := http.FileServer(http.Dir(h.ServeFolder))

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if h.HeaderMapEnabled == true {
			w = common.WriteHeadersToResponse(w, h.HeaderMap)
		}
		if _, err := os.Stat(path.Join(h.ServeFolder, req.URL.Path)); err != nil {
			req.URL.Path = h.Error404FilePath
			req.RequestURI = req.URL.Path
		}
		handler.ServeHTTP(w, req)
	})
}

// ServeHandler ...
// serves a folder
func (h *Handler) ServeHandler() (handler http.Handler) {
	switch {
	case h.VueJSHistoryMode == true:
		handler = h.serveHandlerVuejsHistoryMode()
	default:
		handler = h.serveHandlerStandard()
	}
	if h.GzipEnabled {
		handler = gziphandler.GzipHandler(handler)
	}
	return handler
}
