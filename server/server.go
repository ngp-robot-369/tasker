package server

import (
	"context"
	"log"
	"net/http"
	"tasker/model"
	"tasker/tasks"
)

type Handlers struct {
	Tasker model.ITaskerService
}
type Server struct {
	http.Server
	Handlers
}

func New() *Server {
	h := Handlers{
		Tasker: tasks.NewTaskerService(tasks.MakeStorageRam()),
	}
	mux := http.NewServeMux()
	mux.Handle("POST /task", http.HandlerFunc(h.serveCreateTask))
	mux.Handle("GET /task/{id}", http.HandlerFunc(h.serveGetTaskStatus))

	return &Server{
		Server: http.Server{
			Addr:    ":9000",
			Handler: mux,
		},
		Handlers: h,
	}
}

func (s *Server) Start() error {
	log.Printf("Starting HTTP server on %s...", s.Addr)
	err := s.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (s *Server) Stop(ctx context.Context) {
	log.Printf("Stopping HTTP server...")
	if err := s.Server.Shutdown(ctx); err != nil {
		log.Printf("Failed to stop HTTP server: %v", err)
	}

	log.Printf("Stopping Tasker service...")
	if err := s.Tasker.Shutdown(ctx); err != nil {
		log.Printf("Failed to stop tasker service: %v", err)
	}
}
