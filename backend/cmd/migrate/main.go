package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"lottery-pool-manager/utils"
)

func main() {
	godotenv.Load()

	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbPort := utils.GetEnv("DB_PORT", "5432")
	dbUser := utils.GetEnv("DB_USER", "lottery")
	dbPassword := utils.GetEnv("DB_PASSWORD", "lottery123")
	dbName := utils.GetEnv("DB_NAME", "lottery_pool")

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	migrationsDir := filepath.Join("migrations")

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("failed to read migrations directory: %v", err)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".sql" {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, file := range files {
		content, err := os.ReadFile(filepath.Join(migrationsDir, file))
		if err != nil {
			log.Fatalf("failed to read migration %s: %v", file, err)
		}

		if _, err := pool.Exec(context.Background(), string(content)); err != nil {
			log.Fatalf("failed to execute migration %s: %v", file, err)
		}

		log.Printf("applied: %s", file)
	}

	log.Println("migrations completed successfully")
}
