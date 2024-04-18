package server

import (
	"context"
	"log"
	"net/http"
)

var server = &http.Server{
	Addr: ":9090",
}

func init() {
	for _, r := range GetRoutes() {
		http.HandleFunc(r.Pattern, r.Handler)
	}
}

func Start() error {
	log.Printf("Starting HTTP server on %s...", server.Addr)
	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func Stop(ctx context.Context) {
	log.Printf("Stopping HTTP server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Failed to stop HTTP server: %v", err)
	}
}
