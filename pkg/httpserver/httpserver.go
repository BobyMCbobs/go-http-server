package httpserver

import (
	"crypto/tls"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"gitlab.com/safesurfer/go-http-server/pkg/common"
	"gitlab.com/safesurfer/go-http-server/pkg/handlers"
	"gitlab.com/safesurfer/go-http-server/pkg/metrics"
)

type WebServer struct {
	AppPort            string
	EnvFile            string
	Error404FilePath   string
	GzipEnabled        bool
	HTTPPort           string
	HTTPSPort          string
	HTTPSPortEnabled   bool
	HeaderMapEnabled   bool
	HeaderMapPath      string
	HealthPort         string
	HealthPortEnabled  bool
	MetricsPort        string
	MetricsPortEnabled bool
	RealIPHeader       string
	ServeFolder        string
	TLSCertPath        string
	TLSConfig          *tls.Config
	TLSKeyPath         string
	TemplateMapPath    string
	TemplateMapEnabled bool
	VueJSHistoryMode   bool
	Handler            *handlers.Handler
}

func NewWebServer() *WebServer {
	return &WebServer{
		AppPort:            common.GetAppPort(),
		EnvFile:            common.GetAppEnvFile(),
		Error404FilePath:   common.Get404PageFileName(),
		GzipEnabled:        common.GetEnableGZIP(),
		HTTPPort:           common.GetAppPort(),
		HTTPSPort:          common.GetAppHTTPSPort(),
		HTTPSPortEnabled:   common.GetAppEnableHTTPS(),
		HeaderMapEnabled:   common.GetHeaderSetEnable(),
		HeaderMapPath:      common.GetHeaderMapPath(),
		HealthPort:         common.GetAppHealthPort(),
		HealthPortEnabled:  common.GetAppHealthPortEnabled(),
		MetricsPort:        common.GetAppMetricsPort(),
		MetricsPortEnabled: common.GetAppMetricsEnabled(),
		RealIPHeader:       common.GetAppRealIPHeader(),
		ServeFolder:        common.GetServeFolder(),
		TLSCertPath:        common.GetAppHTTPSCrtPath(),
		TLSKeyPath:         common.GetAppHTTPSKeyPath(),
		TemplateMapPath:    common.GetTemplateMapPath(),
		TemplateMapEnabled: true,
		VueJSHistoryMode:   common.GetVuejsHistoryMode(),
	}
}

func (w *WebServer) SetServeFolder(path string) *WebServer {
	w.ServeFolder = path
	return w
}

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

func (w *WebServer) LoadTemplateMap() *WebServer {
	if w.VueJSHistoryMode != true {
		return w
	}
	configMap, err := common.LoadMapConfig(w.TemplateMapPath)
	if err != nil {
		log.Panicf("[fatal] Error template map: %v\n", err)
	}
	w.Handler.TemplateMap = common.EvaluateEnvFromMap(configMap)
	return w
}

func (w *WebServer) LoadHeaderMap() *WebServer {
	if w.HeaderMapEnabled == false {
		return w
	}
	headerMap, err := common.LoadHeaderMapConfig(w.HeaderMapPath)
	if err != nil {
		panic(err)
	}
	headerMap = common.EvaluateEnvFromHeaderMap(headerMap)
	w.Handler.HeaderMap = headerMap
	return w
}

func (w *WebServer) NewHandlerForWebServer() *handlers.Handler {
	return &handlers.Handler{
		ServeFolder:        w.ServeFolder,
		VueJSHistoryMode:   w.VueJSHistoryMode,
		HeaderMapEnabled:   w.HeaderMapEnabled,
		TemplateMapEnabled: w.TemplateMapEnabled,
		Error404FilePath:   w.Error404FilePath,
		GzipEnabled:        w.GzipEnabled,
	}
}

func (w *WebServer) NewMetricsFromWebServer() *metrics.Metrics {
	return &metrics.Metrics{
		Enabled: w.MetricsPortEnabled,
		Port:    w.MetricsPort,
	}
}

// HandleWebserver ...
// manages app initialisation
func (w *WebServer) Listen() {
	// bring up the API
	forever := make(chan bool)

	_ = godotenv.Load(w.EnvFile)

	go w.NewMetricsFromWebServer().Handle()

	router := mux.NewRouter().StrictSlash(false)
	router.Use(common.Logging)

	w.Handler = w.NewHandlerForWebServer()
	w.LoadHeaderMap()
	w.LoadTemplateMap()

	fullServePath, _ := filepath.Abs(w.ServeFolder)
	log.Printf("Serving folder '%v'\n", fullServePath)
	router.PathPrefix("/").Handler(w.Handler.ServeHandler())

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
