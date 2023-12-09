package http_server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"pex-backend/counter_service"
)

func InstantiateServer(serverAddress string) {
	server := initServer(serverAddress)

	// handle graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to server: %s\n", err)
		}
	}()
	log.Println("started server")

	<-done
	log.Println("stopping server gracefully")

	// kill the server if its shutdown bogs for more than 5 seconds for whatever reason
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("failed to stop the server cleanly: %s", err)
	}
	log.Println("server is stopped")
}

func initServer(serverAddress string) *http.Server {
	timeout := time.Minute

	mainRouter := chi.NewRouter()
	// attach some useful middlewares
	mainRouter.Use(middleware.Logger)
	mainRouter.Use(middleware.RequestID)
	mainRouter.Use(middleware.Recoverer)
	mainRouter.Use(middleware.RealIP)
	mainRouter.Use(middleware.AllowContentType("application/json"))
	mainRouter.Use(middleware.Heartbeat("/ping"))
	mainRouter.Use(middleware.Timeout(timeout))

	r := chi.NewRouter()
	counterService := &counter_service.CounterService{}
	counterService.InitService(r)

	mainRouter.Mount("/api/v1", r)

	server := &http.Server{
		Addr:              serverAddress,
		ReadTimeout:       timeout,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      timeout,
		Handler:           mainRouter,
	}

	return server
}
