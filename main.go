package main

import (
	"log"
	"net/http"
	"os"

	"lottery-pool-manager/routes"
	"lottery-pool-manager/store"
)

func main() {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		connString = "postgres://lottery:lottery123@localhost:5432/lottery_pool"
	}

	s, err := store.NewPostgreSQLStore(connString)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer s.Close()

	r := routes.Setup(s)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
