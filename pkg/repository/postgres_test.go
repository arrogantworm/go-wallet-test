package repository

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	custom_errors "wallet-app/pkg/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
	pgContainer *postgres.PostgresContainer
	testPG      *postgresDB
	ctx         = context.Background()
)

func spinUpPostgres() (*postgres.PostgresContainer, error) {
	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		dbImage,
		postgres.WithInitScripts("../testutils/init-wallet-db.sh"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)

	if err != nil {
		log.Printf("failed to start container: %s", err)
		return nil, err
	}

	return postgresContainer, nil
}

func TestMain(m *testing.M) {
	var err error

	pgContainer, err = spinUpPostgres()
	if err != nil {
		log.Fatalf("failed to spin up postgres: %v", err)
	}

	containerIP, err := pgContainer.ContainerIP(ctx)
	if err != nil {
		log.Fatalf("failed to get container IP: %v", err)
	}

	cfg := Config{
		Host:    containerIP,
		Port:    "5432",
		User:    dbUser,
		Pass:    dbPassword,
		DBName:  dbName,
		SSLMode: "disable",
	}

	testPG, err = NewPG(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to connect to PG: %v", err)
	}

	code := m.Run()

	if err := testcontainers.TerminateContainer(pgContainer); err != nil {
		log.Printf("failed to terminate container: %v", err)
	}

	os.Exit(code)
}

func TestWalletNotFound(t *testing.T) {
	err := testPG.Deposit(ctx, uuid.New(), 1000)
	assert.Error(t, err)

	assert.True(t, errors.Is(err, custom_errors.ErrWalletNotFound))
}

func TestNewWallet(t *testing.T) {

	newWalletUUID := uuid.New()

	err := testPG.NewWallet(ctx, newWalletUUID, 1000)
	assert.NoError(t, err)

	newWalletBalance, err := testPG.GetBalance(ctx, newWalletUUID)
	assert.NoError(t, err)

	assert.Equal(t, newWalletBalance, 1000)
}

func TestDeposit(t *testing.T) {
	newWalletUUID := uuid.New()

	err := testPG.NewWallet(ctx, newWalletUUID, 1000)
	assert.NoError(t, err)

	err = testPG.Deposit(ctx, newWalletUUID, 500)
	assert.NoError(t, err)

	newWalletBalance, err := testPG.GetBalance(ctx, newWalletUUID)
	assert.NoError(t, err)

	assert.Equal(t, newWalletBalance, 1500)
}

func TestWithdraw(t *testing.T) {
	newWalletUUID := uuid.New()

	err := testPG.NewWallet(ctx, newWalletUUID, 1000)
	assert.NoError(t, err)

	err = testPG.Withdraw(ctx, newWalletUUID, 500)
	assert.NoError(t, err)

	newWalletBalance, err := testPG.GetBalance(ctx, newWalletUUID)
	assert.NoError(t, err)

	assert.Equal(t, newWalletBalance, 500)
}

func TestWithdrawError(t *testing.T) {
	newWalletUUID := uuid.New()

	err := testPG.NewWallet(ctx, newWalletUUID, 1000)
	assert.NoError(t, err)

	err = testPG.Withdraw(ctx, newWalletUUID, 1500)
	assert.Error(t, err)

	assert.True(t, errors.Is(err, custom_errors.ErrNotEnoughFunds))
}
