package testutils

import (
	"context"
	"sync"
	"testing"
	"wallet-app/pkg/repository"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	dbName     = "database"
	dbUser     = "user"
	dbPassword = "password"
	dbImage    = "postgres:17.5"
)

var (
	TestPG      *repository.Repository
	testOnce    sync.Once
	pgContainer *postgres.PostgresContainer
	ctx         = context.Background()
)

func SetupTestPG(t *testing.T) *repository.Repository {
	testOnce.Do(func() {
		var err error
		// t.Fatalf(os.Executable())
		pgContainer, err = postgres.Run(
			ctx,
			dbImage,
			postgres.WithInitScripts("../testutils/init-wallet-db.sh"),
			postgres.WithDatabase(dbName),
			postgres.WithUsername(dbUser),
			postgres.WithPassword(dbPassword),
			postgres.BasicWaitStrategies(),
		)
		if err != nil {
			t.Fatalf("failed to start postgres container: %v", err)
		}

		ip, err := pgContainer.ContainerIP(ctx)
		if err != nil {
			t.Fatalf("failed to get container ip: %v", err)
		}

		cfg := repository.Config{
			Host:    ip,
			Port:    "5432",
			User:    "user",
			Pass:    "password",
			DBName:  "database",
			SSLMode: "disable",
		}

		pg, err := repository.NewPG(ctx, cfg)
		if err != nil {
			t.Fatalf("pg init error: %v", err)
		}

		TestPG = repository.NewRepository(pg)

		go func() {
			<-ctx.Done()
			_ = testcontainers.TerminateContainer(pgContainer)
		}()
	})

	return TestPG
}
