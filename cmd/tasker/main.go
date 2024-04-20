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
	var (
		srv      = server.New()
		sigTerm  = make(chan os.Signal, 1)
		waitStop = make(chan struct{}, 1)
	)
	go func() {
		signal.Notify(sigTerm, syscall.SIGTERM, syscall.SIGINT)
		log.Printf("Stop signal received: %v, exiting...", <-sigTerm)
		srv.Stop(context.TODO())
		close(waitStop)
	}()

	if err := srv.Start(); err != nil {
		log.Panicf("Failed to start server: %v", err)
	}
	<-waitStop
	log.Printf("Gracefully stopped")
}
