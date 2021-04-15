package gohttpserver

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"gitlab.com/safesurfer/go-http-server/pkg/common"
	"gitlab.com/safesurfer/go-http-server/pkg/handlers"
	"gitlab.com/safesurfer/go-http-server/pkg/metrics"
)

// HandleWebserver ...
// manages app initialisation
func HandleWebserver() {
	// bring up the API
	forever := make(chan bool)

	envFile := common.GetAppEnvFile()
	_ = godotenv.Load(envFile)

	go metrics.Handle()

	port := common.GetAppPort()
	router := mux.NewRouter().StrictSlash(false)
	router.Use(common.Logging)

	serveFolder := common.GetServeFolder()
	fullServePath, _ := filepath.Abs(serveFolder)
	log.Printf("Serving folder '%v'\n", fullServePath)
	router.PathPrefix("/").Handler(handlers.ServeHandler(serveFolder))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET"},
		AllowCredentials: true,
	})

	// Serve regular HTTP
	srv := &http.Server{
		Handler:      c.Handler(router),
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go func() {
		log.Println("Listening on", port)
		log.Fatal(srv.ListenAndServe())
	}()

	// Optionally, serve HTTPS
	useTLS, err := strconv.ParseBool(common.GetAppEnableHTTPS())
	if err != nil {
		log.Panicf("[fatal] Error parsing APP_ENABLE_HTTPS: %v\n", err)
	}
	if useTLS {
		// Load certs
		tlsConfig := &tls.Config{}
		tlsConfig.Certificates = make([]tls.Certificate, 1)
		tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(os.Getenv("APP_HTTPS_CRT_PATH"), os.Getenv("APP_HTTPS_KEY_PATH"))
		if err != nil {
			log.Panicf("[fatal] Error loading certs: %v\n", err)
		}
		tlsConfig.BuildNameToCertificate()
		// Create TLS listener and server
		tlsPort := common.GetAppHTTPSPort()
		srvTLS := &http.Server{
			Handler:      c.Handler(router),
			Addr:         tlsPort,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
			TLSConfig:    tlsConfig,
		}
		listener, err := tls.Listen("tcp", tlsPort, tlsConfig)
		if err != nil {
			log.Panicf("[fatal] Error creating TLS listener: %v\n", err)
		}
		go func() {
			log.Println("Listening on", tlsPort)
			log.Println(srvTLS.Serve(listener))
		}()
	}
	<-forever
}
