package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/NYTimes/gziphandler"

	"gitlab.com/BobyMCbobs/go-http-server/pkg/common"
)

var (
	fileServeDisallowList = []string{
		// TODO add .git and sub directory listing to disallow list
		"/.ghs.yaml",
		"/.env",
	}
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
		if h.HeaderMapEnabled {
			w = common.WriteHeadersToResponse(w, h.HeaderMap)
		}
		isDisallowed := false
		for _, f := range fileServeDisallowList {
			if match, _ := path.Match(f, req.URL.Path); match {
				isDisallowed = true
				break
			}
		}

		// static files
		if strings.Contains(req.URL.Path, ".") && !isDisallowed {
			handler.ServeHTTP(w, req)
			return
		}

		// frontend views
		indexPath := path.Join(h.ServeFolder, "/index.html")
		tmpl, err := template.ParseFiles(indexPath)
		if err != nil {
			log.Println("warning: unable to parse template html:", err)
			http.Error(w, "500 internal error", http.StatusInternalServerError)
			return
		}
		var buf bytes.Buffer
		defer buf.Reset()
		if err := tmpl.ExecuteTemplate(&buf, tmpl.Name(), h.TemplateMap); err != nil {
			log.Println("warning: unable to execute template html:", err)
			http.Error(w, "500 internal error", http.StatusInternalServerError)
		}
		fmt.Fprint(w, &buf)
	})
}

// serveHandlerStandard ...
// handles sending the serve folder
func (h *Handler) serveHandlerStandard() http.Handler {
	handler := http.FileServer(http.Dir(h.ServeFolder))

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if h.HeaderMapEnabled {
			w = common.WriteHeadersToResponse(w, h.HeaderMap)
		}
		isDisallowed := false
		for _, f := range fileServeDisallowList {
			if match, _ := path.Match(f, req.URL.Path); match {
				isDisallowed = true
				break
			}
		}
		if _, err := os.Stat(path.Join(h.ServeFolder, req.URL.Path)); err != nil || isDisallowed {
			w.WriteHeader(http.StatusNotFound)
			http.ServeFile(w, req, path.Join(h.ServeFolder, h.Error404FilePath))
			return
		}
		handler.ServeHTTP(w, req)
	})
}

// ServeHandler ...
// serves a folder
func (h *Handler) ServeHandler() (handler http.Handler) {
	switch {
	case h.VueJSHistoryMode:
		handler = h.serveHandlerVuejsHistoryMode()
	default:
		handler = h.serveHandlerStandard()
	}
	if h.GzipEnabled {
		handler = gziphandler.GzipHandler(handler)
	}
	return handler
}

// ServeStandardRedirect ...
// handles a standard path redirect
func (h *Handler) ServeStandardRedirect(from string, to string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// TODO revisit disallowing certain paths like '/' or ''
		toURL, err := url.Parse(to)
		if err != nil {
			log.Printf("Unable to parse redirection destination URL '%v' for route '%v'\n", to, from)
			http.Error(w, "fatal: unable to redirect to destination URL", http.StatusInternalServerError)
			return
		}
		toURL.RawQuery = req.URL.Query().Encode()
		log.Printf("redirecting '%v' -> '%v'\n", from, to)
		http.Redirect(w, req, toURL.String(), http.StatusTemporaryRedirect)
	})
}
