package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	postgres struct {
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
	pgInstance *postgres
	pgOnce     sync.Once
)

func NewPG(ctx context.Context, cfg Config) (*postgres, error) {
	var err error
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
	pgOnce.Do(func() {
		var db *pgxpool.Pool
		db, err = pgxpool.New(ctx, connString)
		if err == nil {
			pgInstance = &postgres{db}
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

func (pg *postgres) ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *postgres) Close() {
	pg.db.Close()
}
