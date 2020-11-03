package gohttpserver

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/joho/godotenv"

	"gitlab.com/safesurfer/go-http-server/pkg/common"
	"gitlab.com/safesurfer/go-http-server/pkg/handlers"
	"gitlab.com/safesurfer/go-http-server/pkg/metrics"
)

// HandleWebserver ...
// manages app initialisation
func HandleWebserver() {
	// bring up the API

	envFile := common.GetAppEnvFile()
	_ = godotenv.Load(envFile)

	go metrics.Handle()

	port := common.GetAppPort()
	router := mux.NewRouter().StrictSlash(false)
	router.Use(common.Logging)

	serveFolder := common.GetServeFolder()
	router.PathPrefix("/").Handler(handlers.ServeHandler(serveFolder))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET"},
		AllowCredentials: true,
	})

	srv := &http.Server{
		Handler:      c.Handler(router),
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Listening on", port)
	log.Fatal(srv.ListenAndServe())
}
