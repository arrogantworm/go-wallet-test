package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	postgresDB struct {
		db *pgxpool.Pool
	}

	Config struct {
		Host    string
		Port    string
		User    string
		Pass    string
		DBName  string
		SSLMode string
	}
)

var (
	pgInstance *postgresDB
	pgOnce     sync.Once
)

func NewPG(ctx context.Context, cfg Config) (*postgresDB, error) {
	var err error
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
	pgOnce.Do(func() {
		var db *pgxpool.Pool
		db, err = pgxpool.New(ctx, connString)
		if err == nil {
			pgInstance = &postgresDB{db}
		}
	})

	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	if err = pgInstance.ping(ctx); err != nil {
		return nil, err
	}

	return pgInstance, nil
}

func (pg *postgresDB) ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *postgresDB) Close() {
	pg.db.Close()
}
