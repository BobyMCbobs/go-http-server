package httpserver

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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
	EnvFile               string
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
	handler               *handlers.Handler
}

// NewWebServer returns a default WebServer, as per environment configuration
func NewWebServer() *WebServer {
	w := &WebServer{
		AppPort:               common.GetAppPort(),
		EnvFile:               common.GetAppEnvFile(),
		Error404FilePath:      common.Get404PageFileName(),
		GzipEnabled:           common.GetEnableGZIP(),
		HTTPPort:              common.GetAppPort(),
		HTTPSPort:             common.GetAppHTTPSPort(),
		HTTPSPortEnabled:      common.GetAppEnableHTTPS(),
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
	if cfg, err := common.LoadDotfileConfig(w.ServeFolder); err == nil {
		w.VueJSHistoryMode = cfg.HistoryMode
		w.RedirectRoutes = cfg.RedirectRoutes
		w.HeaderMap = cfg.HeaderMap
		w.TemplateMap = cfg.TemplateMap

		if w.HeaderMap != nil {
			w.HeaderMapEnabled = true
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
func (w *WebServer) LoadTLS() *WebServer {
	w.TLSConfig = &tls.Config{}
	w.TLSConfig.Certificates = make([]tls.Certificate, 1)
	loadedCert, err := tls.LoadX509KeyPair(w.TLSCertPath, w.TLSKeyPath)
	if err != nil {
		log.Panicf("[fatal] Error loading certs: %v\n", err)
	}
	w.TLSConfig.Certificates[0] = loadedCert
	w.TLSConfig.BuildNameToCertificate()
	return w
}

// LoadTemplateMap loads the template map from the path
func (w *WebServer) LoadTemplateMap() *WebServer {
	if w.VueJSHistoryMode != true {
		return w
	}
	if w.TemplateMap == nil {
		if _, err := os.Stat(w.TemplateMapPath); os.IsNotExist(err) {
			log.Printf("[notice] history mode is enabled, template maps (currently set to '%v') can also be used\n", w.TemplateMapPath)
			return w
		}
		configMap, err := common.LoadMapConfig(w.TemplateMapPath)
		if err != nil {
			log.Panicf("[fatal] Error template map: %v\n", err)
		}
		w.TemplateMap = configMap
	}
	w.TemplateMap = common.EvaluateEnvFromMap(w.TemplateMap)
	w.handler.TemplateMap = w.TemplateMap
	return w
}

// SetTemplateMap set the template map
func (w *WebServer) SetTemplateMap(input map[string]string) *WebServer {
	w.handler.TemplateMap = common.EvaluateEnvFromMap(input)
	return w
}

// LoadHeaderMap loads the header map from the path
func (w *WebServer) LoadHeaderMap() *WebServer {
	if w.HeaderMapEnabled == false {
		return w
	}
	if w.HeaderMap == nil {
		if _, err := os.Stat(w.HeaderMapPath); os.IsNotExist(err) {
			log.Printf("[notice] history mode is enabled, template maps (currently set to '%v') can also be used\n", w.TemplateMapPath)
			return w
		}
		headerMap, err := common.LoadHeaderMapConfig(w.HeaderMapPath)
		if err != nil {
			panic(err)
		}
		w.HeaderMap = headerMap
	}
	w.HeaderMap = common.EvaluateEnvFromHeaderMap(w.HeaderMap)
	w.handler.HeaderMap = w.HeaderMap
	return w
}

// SetHeaderMap sets the header map
func (w *WebServer) SetHeaderMap(input map[string][]string) *WebServer {
	w.handler.HeaderMap = common.EvaluateEnvFromHeaderMap(input)
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
func (w *WebServer) Listen() {
	// bring up the API
	forever := make(chan bool)

	_ = godotenv.Load(w.EnvFile)

	go w.NewMetricsFromWebServer().Handle()

	router := mux.NewRouter().StrictSlash(false)
	router.Use(common.Logging)

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

	w.LoadHeaderMap()
	w.LoadTemplateMap()

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
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET"},
		AllowCredentials: true,
	})

	// Serve regular HTTP
	srv := &http.Server{
		Handler:      c.Handler(router),
		Addr:         w.AppPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go func() {
		log.Println("Listening on", w.AppPort)
		log.Fatal(srv.ListenAndServe())
	}()

	// Optionally, serve HTTPS
	if w.HTTPSPortEnabled {
		// Load certs
		w.LoadTLS()
		// Create TLS listener and server
		srvTLS := &http.Server{
			Handler:      c.Handler(router),
			Addr:         w.HTTPSPort,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
			TLSConfig:    w.TLSConfig,
		}
		listener, err := tls.Listen("tcp", w.HTTPSPort, w.TLSConfig)
		if err != nil {
			log.Panicf("[fatal] Error creating TLS listener: %v\n", err)
		}
		go func() {
			log.Println("Listening on", w.HTTPSPort)
			log.Println(srvTLS.Serve(listener))
		}()
	}
	<-forever
}
