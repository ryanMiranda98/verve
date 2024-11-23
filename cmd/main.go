package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ryanMiranda98/verve/api"
)

func main() {
	apiServer := api.NewApiServer(":8080")
	apiServer.SetupRouter()

	go func() {
		if err := apiServer.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error during ListenAndServe:%s\n", err)
		}
	}()

	go func() {
		if err := api.LogUniqueRequests(); err != nil {
			log.Fatalf("Error during TrackUniqueRequests")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutdown server...\n")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
