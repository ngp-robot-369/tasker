package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tasker/server"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		log.Printf("Stop signal received: %v, exiting...", <-c)
		server.Stop(context.TODO())
	}()

	if err := server.Start(); err != nil {
		log.Panicf("Failed to start server: %v", err)
	}
	log.Printf("Gracefully stopped")
}
