package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KOJIMEISTER/it_russian_stat/internal/domain"
	"github.com/KOJIMEISTER/it_russian_stat/internal/messaging"
)

func main() {
	mux := http.NewServeMux()

	rabbitService, err := messaging.NewRabbitMQService()
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitService.Close()

	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/update", updateHandler(rabbitService))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	shutdown(srv)
}

func shutdown(srv *http.Server) {
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server forced to shutdown:", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func updateHandler(rabbitService *messaging.RabbitMQService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req domain.UpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		if err := rabbitService.Publish(req); err != nil {
			http.Error(w, "Failed to publish message", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
	}
}
