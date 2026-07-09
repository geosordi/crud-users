package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"crud-users/internal/config"
	"crud-users/internal/database"
	"crud-users/internal/repository"
	"crud-users/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	client, err := database.Connect(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}
	defer database.Close(ctx, client)

	collection := client.Database(cfg.DBName).Collection(cfg.CollectionName)

	userRepo := repository.NewMongoUserRepository(collection)
	userSvc := service.NewUserService(userRepo)

	srv := &http.Server{Addr: ":8080", Handler: newRouter(userSvc)}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown error: %v", err)
		}
	}()

	log.Println("server listening on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
