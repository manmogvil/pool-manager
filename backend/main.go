package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"lottery-pool-manager/routes"
	"lottery-pool-manager/store"
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

	s, err := store.NewPostgreSQLStore(connString)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer s.Close()

	port := utils.GetEnv("PORT", "8080")

	r := routes.Setup(s)

	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
