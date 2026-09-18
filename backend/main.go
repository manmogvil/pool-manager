package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"lottery-pool-manager/routes"
	"lottery-pool-manager/scheduler"
	"lottery-pool-manager/services"
	"lottery-pool-manager/store"
	"lottery-pool-manager/utils"
)

func startScheduler(ctx context.Context, store scheduler.Store) {
	enabled := utils.GetEnv("SCHEDULER_ENABLED", "false")
	if enabled != "true" {
		log.Println("Scheduler disabled")
		return
	}

	cronExpr := utils.GetEnv("SCHEDULER_CRON", "0 6 * * 1")
	api := services.NewLoteriaAPIClient()
	sched := scheduler.New(store, api, cronExpr)
	go sched.Start(ctx)
	log.Printf("Scheduler enabled with cron: %s", cronExpr)
}

func main() {
	godotenv.Load()

	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbPort := utils.GetEnv("DB_PORT", "5432")
	dbUser := utils.GetEnv("DB_USER", "lottery")
	dbPassword := utils.GetEnv("DB_PASSWORD", "lottery123")
	dbName := utils.GetEnv("DB_NAME", "lottery_pool")

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	s, err := store.NewPostgreSQLStore(connString)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer s.Close()

	port := utils.GetEnv("PORT", "8080")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startScheduler(ctx, s)

	r := routes.Setup(s)

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on :%s", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
