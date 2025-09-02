package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"

	"github.com/yourname/eventdrivenorders/internal/order"
)

func main() {
	db, err := sql.Open("postgres", "postgres://user:pass@localhost:5432/ordersdb?sslmode=disable")
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	// Pooling (tune as needed)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	repo := order.NewRepository(db)
	handler := order.NewHandler(repo)

	r := chi.NewRouter()
	r.Post("/orders", handler.CreateOrder)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// start server
	go func() {
		log.Printf("Server running on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
	log.Println("Server stopped")
}
