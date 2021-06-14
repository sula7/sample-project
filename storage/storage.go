package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	Storager interface {
	}
	Storage struct {
		Pool *pgxpool.Pool
	}
)

func New(connStr string) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	pool, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("open db pool error: ", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("db pool opened but ping err: ", err)
	}

	return &Storage{Pool: pool}, err
}
