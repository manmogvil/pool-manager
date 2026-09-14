package store

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"lottery-pool-manager/models"
	"lottery-pool-manager/utils"
)

var testStore *PostgreSQLStore

const testDBName = "lottery_pool_test"

func TestMain(m *testing.M) {
	godotenv.Load("../.env")

	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbPort := utils.GetEnv("DB_PORT", "5432")
	dbUser := utils.GetEnv("DB_USER", "lottery")
	dbPassword := utils.GetEnv("DB_PASSWORD", "lottery123")

	baseConn := os.Getenv("TEST_DATABASE_URL")
	if baseConn == "" {
		baseConn = fmt.Sprintf("postgres://%s:%s@%s:%s/lottery_pool", dbUser, dbPassword, dbHost, dbPort)
	}

	// Connect to default database to create/drop test DB
	basePool, err := pgxpool.New(context.Background(), baseConn)
	if err != nil {
		fmt.Printf("failed to connect to base database: %v\n", err)
		os.Exit(1)
	}

	// Drop test DB if exists, then create it
	_, _ = basePool.Exec(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName))
	_, err = basePool.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", testDBName))
	if err != nil {
		fmt.Printf("failed to create test database: %v\n", err)
		os.Exit(1)
	}
	basePool.Close()

	// Connect to test database
	testConn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, testDBName)
	testStore, err = NewPostgreSQLStore(testConn)
	if err != nil {
		fmt.Printf("failed to connect to test database: %v\n", err)
		os.Exit(1)
	}

	// Run migrations
	if err := runMigrations(); err != nil {
		fmt.Printf("failed to run migrations: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	// Cleanup: close connection and drop test database
	testStore.Close()

	basePool, err = pgxpool.New(context.Background(), baseConn)
	if err == nil {
		_, _ = basePool.Exec(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName))
		basePool.Close()
	}

	os.Exit(code)
}

func runMigrations() error {
	migrationsDir := filepath.Join("..", "migrations")

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
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
			return fmt.Errorf("failed to read migration %s: %w", file, err)
		}

		if _, err := testStore.pool.Exec(context.Background(), string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}
	}

	return nil
}

func cleanup(t *testing.T) {
	t.Helper()
	_, _ = testStore.pool.Exec(context.Background(),
		`TRUNCATE tickets, draws, contributions, lottery_games, participants RESTART IDENTITY CASCADE;`)
}

func seedParticipant(t *testing.T, name, email string) models.Participant {
	t.Helper()
	p, err := testStore.CreateParticipant(name, email)
	if err != nil {
		t.Fatalf("seedParticipant failed: %v", err)
	}
	return p
}

func seedGame(t *testing.T, name, drawDays string, price float64) models.LotteryGame {
	t.Helper()
	g, err := testStore.CreateLotteryGame(name, drawDays, price)
	if err != nil {
		t.Fatalf("seedGame failed: %v", err)
	}
	return g
}

func seedDraw(t *testing.T, gameID int, drawDate time.Time) models.Draw {
	t.Helper()
	d, err := testStore.CreateDraw(gameID, drawDate)
	if err != nil {
		t.Fatalf("seedDraw failed: %v", err)
	}
	return d
}
