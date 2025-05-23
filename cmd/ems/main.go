package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amito07/ems/internal/config"
	"github.com/amito07/ems/internal/http/routes/rootRouter"
)

func main() {
	fmt.Println("Welcome to the Education Management System (EMS)!")

	// Load the configuration
	cfg := config.MustLoadConfig()

	// database setup
	// setup router
	router := rootrouter.RouterInit()

	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Starting server...", slog.String("address", cfg.Addr))
	fmt.Printf("Server started at %s", cfg.Addr)

	// Graceful shutdown
	// Create a channel to listen for OS signals
	// This will allow us to gracefully shut down the server
	// when we receive a termination signal
	done := make(chan os.Signal, 1)

	// Listen for OS signals
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()

		if err != nil {
			fmt.Printf("Failed to start server: %s", err.Error())
		}
	}()

	// Wait for a signal to exit
	<-done

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", "error", slog.String("error", err.Error()))
	}

	slog.Info("Server shut down gracefully")
}
