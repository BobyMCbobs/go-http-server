package httpserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"gitlab.com/BobyMCbobs/go-http-server/pkg/common"
	"gitlab.com/BobyMCbobs/go-http-server/pkg/handlers"
	"gitlab.com/BobyMCbobs/go-http-server/pkg/metrics"
)

// ExtraHandler ...
// add extra endpoints to listen on
type ExtraHandler struct {
	Path        string
	HandlerFunc http.HandlerFunc
	HTTPMethods []string
}

// WebServer configures the runtime
type WebServer struct {
	AppPort               string
	HTTPAllowedOrigins    []string
	Error404FilePath      string
	ExtraHandlers         []*ExtraHandler
	ExtraMiddleware       []func(http.Handler) http.Handler
	GzipEnabled           bool
	HTTPPort              string
	HTTPSPort             string
	HTTPSPortEnabled      bool
	HeaderMap             map[string][]string
	HeaderMapEnabled      bool
	HeaderMapPath         string
	HealthPort            string
	HealthPortEnabled     bool
	MetricsPort           string
	MetricsPortEnabled    bool
	RealIPHeader          string
	RedirectRoutes        map[string]string
	RedirectRoutesEnabled bool
	RedirectRoutesPath    string
	ServeFolder           string
	TLSCertPath           string
	TLSConfig             *tls.Config
	TLSKeyPath            string
	TemplateMap           map[string]string
	TemplateMapEnabled    bool
	TemplateMapPath       string
	VueJSHistoryMode      bool

	handler       *handlers.Handler
	server        *http.Server
	serverTLS     *http.Server
	dotfileLoaded bool
}

// NewWebServer returns a default WebServer, as per environment configuration
//
// TODO clean up function
func NewWebServer() *WebServer {
	httpOrigins, err := common.GetHTTPAllowedOrigins()
	if err != nil {
		log.Printf("error: failed to load allowed http origins: %v\n", err)
	}
	w := &WebServer{
		AppPort:               common.GetAppPort(),
		Error404FilePath:      common.Get404PageFileName(),
		GzipEnabled:           common.GetEnableGZIP(),
		HTTPPort:              common.GetAppPort(),
		HTTPSPort:             common.GetAppHTTPSPort(),
		HTTPSPortEnabled:      common.GetAppEnableHTTPS(),
		HTTPAllowedOrigins:    httpOrigins,
		HeaderMapEnabled:      common.GetHeaderSetEnable(),
		HeaderMapPath:         common.GetHeaderMapPath(),
		HealthPort:            common.GetAppHealthPort(),
		HealthPortEnabled:     common.GetAppHealthPortEnabled(),
		MetricsPort:           common.GetAppMetricsPort(),
		MetricsPortEnabled:    common.GetAppMetricsEnabled(),
		RealIPHeader:          common.GetAppRealIPHeader(),
		RedirectRoutesEnabled: common.GetRedirectRoutesEnabled(),
		RedirectRoutesPath:    common.GetRedirectRoutesPath(),
		ServeFolder:           common.GetServeFolder(),
		TLSCertPath:           common.GetAppHTTPSCrtPath(),
		TLSKeyPath:            common.GetAppHTTPSKeyPath(),
		TemplateMapEnabled:    true,
		TemplateMapPath:       common.GetTemplateMapPath(),
		VueJSHistoryMode:      common.GetVuejsHistoryMode(),
		handler:               &handlers.Handler{},
	}
	cfg, err := common.LoadDotfileConfig(w.ServeFolder)
	if err != nil {
		log.Printf("error loading dotfile config: %v\n", err)
	} else if cfg != nil {
		w.dotfileLoaded = true
		w.VueJSHistoryMode = cfg.HistoryMode
		if cfg.RedirectRoutes != nil {
			w.RedirectRoutes = cfg.RedirectRoutes
		}
		if cfg.HeaderMap != nil {
			w.HeaderMap = cfg.HeaderMap
		}
		if cfg.TemplateMap != nil {
			w.TemplateMap = cfg.TemplateMap
		}
		if w.HeaderMap != nil {
			w.HeaderMapEnabled = true
		}
		w.Error404FilePath = cfg.Error404FilePath
		if w.Error404FilePath == "" {
			w.Error404FilePath = common.Get404PageFileName()
		}
	}
	router := mux.NewRouter().StrictSlash(false)
	router.Use(common.Logging)
	// NOTE is this ever called?
	for _, m := range w.ExtraMiddleware {
		router.Use(m)
	}
	if w.RedirectRoutesEnabled {
		if w.RedirectRoutes == nil {
			redirectRoutes, err := common.LoadRedirectRoutesConfig(w.RedirectRoutesPath)
			if err != nil {
				log.Println("Warning: failed to load redirect routes")
			}
			w.RedirectRoutes = redirectRoutes
		}
		for from, to := range w.RedirectRoutes {
			router.HandleFunc(from, w.handler.ServeStandardRedirect(from, to)).Methods(http.MethodGet)
		}
	}

	if _, err := w.LoadHeaderMap(); err != nil {
		log.Printf("error: failed to load header map: %v\n", err)
	}
	if _, err := w.LoadTemplateMap(); err != nil {
		log.Printf("error: failed to load template map: %v\n", err)
	}

	for _, h := range w.ExtraHandlers {
		if h.Path == "/" {
			log.Println("Warning: path / not allowed for extra handlers")
			continue
		}
		router.HandleFunc(h.Path, h.HandlerFunc).Methods(h.HTTPMethods...)
	}
	w.handler = w.newHandlerForWebServer()

	fullServePath, _ := filepath.Abs(w.ServeFolder)
	log.Printf("Serving folder '%v'\n", fullServePath)
	router.PathPrefix("/").Handler(w.handler.ServeHandler())

	c := cors.New(cors.Options{
		AllowedOrigins:   w.HTTPAllowedOrigins,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET"},
		AllowCredentials: true,
	})

	// Serve regular HTTP
	w.server = &http.Server{
		Handler:      c.Handler(router),
		Addr:         w.AppPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	if w.HTTPSPortEnabled {
		if _, err := w.LoadTLS(); err != nil {
			log.Printf("error: failed to load TLS: %v\n", err)
		}
		w.serverTLS = &http.Server{
			Handler:      c.Handler(router),
			Addr:         w.HTTPSPort,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
			TLSConfig:    w.TLSConfig,
		}
	}

	return w
}

// SetServeFolder sets the path to the ServeFolder
func (w *WebServer) SetServeFolder(path string) *WebServer {
	w.ServeFolder = path
	return w
}

// SetExtraHandlers sets extra http handlers
func (w *WebServer) SetExtraHandlers(hs ...*ExtraHandler) *WebServer {
	w.ExtraHandlers = hs
	return w
}

// SetExtraMiddleware sets extra http middleware
func (w *WebServer) SetExtraMiddleware(m ...func(http.Handler) http.Handler) *WebServer {
	w.ExtraMiddleware = m
	return w
}

// LoadTLS loads in the TLS certs
func (w *WebServer) LoadTLS() (*WebServer, error) {
	w.TLSConfig = &tls.Config{}
	w.TLSConfig.Certificates = make([]tls.Certificate, 1)
	loadedCert, err := tls.LoadX509KeyPair(w.TLSCertPath, w.TLSKeyPath)
	if err != nil {
		return w, err
	}
	w.TLSConfig.Certificates[0] = loadedCert
	return w, nil
}

// LoadTemplateMap loads the template map from the path
func (w *WebServer) LoadTemplateMap() (*WebServer, error) {
	if w.TemplateMap == nil && !w.dotfileLoaded {
		if _, err := os.Stat(w.TemplateMapPath); os.IsNotExist(err) {
			log.Printf("[notice] history mode templating is enabled, template maps (currently set to '%v') can also be used\n", w.TemplateMapPath)
			return w, fmt.Errorf("error: template map file not found")
		}
		templateMap, err := common.LoadTemplateMapConfig(w.TemplateMapPath)
		if err != nil {
			return w, err
		}
		w.TemplateMap = common.EvaluateEnvFromMap(templateMap)
	}
	w.handler.TemplateMap = w.TemplateMap
	return w, nil
}

// SetTemplateMap set the template map
func (w *WebServer) SetTemplateMap(input map[string]string) *WebServer {
	if !w.dotfileLoaded {
		input = common.EvaluateEnvFromMap(input)
	}
	w.TemplateMap = input
	w.handler.TemplateMap = input
	return w
}

// LoadHeaderMap loads the header map from the path
func (w *WebServer) LoadHeaderMap() (*WebServer, error) {
	if w.HeaderMap == nil && !w.dotfileLoaded {
		if _, err := os.Stat(w.HeaderMapPath); os.IsNotExist(err) {
			log.Printf("[notice] header templating is enabled, header template maps (currently set to '%v') can also be used\n", w.HeaderMapPath)
			return w, fmt.Errorf("error: header template not found")
		}
		headerMap, err := common.LoadHeaderMapConfig(w.HeaderMapPath)
		if err != nil {
			return w, err
		}
		w.HeaderMap = common.EvaluateEnvFromHeaderMap(headerMap)
	}
	w.handler.HeaderMap = w.HeaderMap
	return w, nil
}

// SetHeaderMap sets the header map
func (w *WebServer) SetHeaderMap(input map[string][]string) *WebServer {
	if !w.dotfileLoaded {
		input = common.EvaluateEnvFromHeaderMap(input)
	}
	w.HeaderMap = input
	w.handler.HeaderMap = input
	return w
}

func (w *WebServer) newHandlerForWebServer() *handlers.Handler {
	return &handlers.Handler{
		ServeFolder:        w.ServeFolder,
		VueJSHistoryMode:   w.VueJSHistoryMode,
		HeaderMapEnabled:   w.HeaderMapEnabled,
		TemplateMapEnabled: w.TemplateMapEnabled,
		Error404FilePath:   w.Error404FilePath,
		GzipEnabled:        w.GzipEnabled,
		HeaderMap:          w.HeaderMap,
		TemplateMap:        w.TemplateMap,
	}
}

// SetHandler sets a new handler
func (w *WebServer) SetHandler(input *handlers.Handler) *WebServer {
	w.handler = input
	return w
}

// NewMetricsFromWebServer returns a new metrics from a webserver
func (w *WebServer) NewMetricsFromWebServer() *metrics.Metrics {
	return &metrics.Metrics{
		Enabled: w.MetricsPortEnabled,
		Port:    w.MetricsPort,
	}
}

// Listen starting listening according to the configuration
func (w *WebServer) Listen(ch ...<-chan bool) {
	go w.NewMetricsFromWebServer().Handle()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if len(ch) == 0 {
			return
		}
		for {
			c, ok := <-ch[0]
			log.Println("<- receieved event:", c, ok)
			if !ok {
				break
			}
			if c {
				done <- os.Interrupt
			}
		}
	}()
	go func() {
		log.Println("Listening on", w.AppPort)
		if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Optionally, serve HTTPS
	if w.HTTPSPortEnabled {
		listener, err := tls.Listen("tcp", w.HTTPSPort, w.TLSConfig)
		if err != nil {
			log.Panicf("[fatal] Error creating TLS listener: %v\n", err)
		}
		go func() {
			log.Println("Listening on", w.HTTPSPort)
			log.Println(w.serverTLS.Serve(listener))
		}()
	}

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := w.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server didn't exit gracefully %v", err)
	}
	if w.HTTPSPortEnabled {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := w.serverTLS.Shutdown(ctx); err != nil {
			log.Fatalf("Server didn't exit gracefully %v", err)
		}
	}
}
