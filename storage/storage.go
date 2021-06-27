package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/jackc/pgx/v4/pgxpool"

	"sample-project/structs"
)

type (
	Storager interface {
		UpDBVersion(dsn string, source source.Driver) error
		SetDBVersion(dsn string, version uint, source source.Driver) error

		CreateDrone(d *structs.Drone) error
		GetUserUUID(u *structs.User) (string, error)
		GetAllDrones() ([]structs.Drone, error)
	}
	Storage struct {
		Pool *pgxpool.Pool
	}
)

func New(dsn string) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("open db pool: %w", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}

	return &Storage{Pool: pool}, err
}

func (s *Storage) UpDBVersion(dsn string, source source.Driver) error {
	m, err := migrate.NewWithSourceInstance("httpfs", source, dsn)
	if err != nil {
		return fmt.Errorf("initialize migrations: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (s *Storage) SetDBVersion(dsn string, version uint, source source.Driver) error {
	m, err := migrate.NewWithSourceInstance("httpfs", source, dsn)
	if err != nil {
		return fmt.Errorf("initialize migrations: %w", err)
	}

	err = m.Migrate(version)
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
