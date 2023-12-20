package metrics

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics configures the metrics handler
type Metrics struct {
	Enabled bool
	Port    string
}

// Handle ...
// HTTP handler for metrics
func (m *Metrics) Handle(ch ...<-chan bool) {
	if !m.Enabled {
		return
	}

	router := mux.NewRouter().StrictSlash(true)
	r := router.Handle("/metrics", promhttp.Handler())
	server := &http.Server{
		Handler:           r.GetHandler(),
		Addr:              m.Port,
		WriteTimeout:      15 * time.Second,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Printf("Metrics listening on %v\n", server.Addr)
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
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server didn't exit gracefully %v", err)
	}
}
