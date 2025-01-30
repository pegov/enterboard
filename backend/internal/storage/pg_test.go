package storage_test

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/pegov/enterboard/backend/internal/storage"
)

func TestWithPostgres(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dbDatabase := "testdb"
	dbUsername := "user"
	dbPassword := "password"
	req := testcontainers.ContainerRequest{
		Image:        "postgres:14.4",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       dbDatabase,
			"POSTGRES_USER":     dbUsername,
			"POSTGRES_PASSWORD": dbPassword,
		},
		Tmpfs: map[string]string{
			"/var/lib/postgresql/data": "rw",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(5 * time.Second),
	}
	ctr, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}

	defer ctr.Terminate(ctx, testcontainers.StopContext(ctx))

	host, err := ctr.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get host: %v", err)
	}

	port, err := ctr.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get mapped port: %v", err)
	}

	pgURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbUsername,
		dbPassword,
		host,
		port.Int(),
		dbDatabase,
	)
	db, err := storage.NewPG(ctx, slog.Default(), pgURL, 10, 10, 10*time.Minute)
	if err != nil {
		t.Fatalf("failed to open pg connection: %v", err)
	}
	defer db.Close()

	var four int
	if err := db.GetContext(ctx, &four, "SELECT 4"); err != nil {
		t.Fatalf("failed to get four: %v", err)
	}

	if got, want := four, 4; got != want {
		t.Fatalf("not 4")
	}
}
