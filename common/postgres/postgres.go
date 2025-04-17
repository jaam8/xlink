package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host           string `yaml:"POSTGRES_HOST" env:"POSTGRES_HOST" env-default:"localhost"`
	Port           uint16 `yaml:"POSTGRES_PORT" env:"POSTGRES_PORT" env-default:"5432"`
	Username       string `yaml:"POSTGRES_USER" env:"POSTGRES_USER" env-default:"postgres"`
	Password       string `yaml:"POSTGRES_PASSWORD" env:"POSTGRES_PASSWORD" env-default:"1234"`
	Database       string `yaml:"POSTGRES_DB" env:"POSTGRES_DB" env-default:"postgres"`
	MaxConns       int32  `yaml:"POSTGRES_MAX_CONNS" env:"POSTGRES_MAX_CONNS" env-default:"10"`
	MinConns       int32  `yaml:"POSTGRES_MIN_CONNS" env:"POSTGRES_MIN_CONNS" env-default:"5"`
	MigrationsPath string `yaml:"POSTGRES_ROOT_MIGRATIONS_PATH" env:"POSTGRES_ROOT_MIGRATIONS_PATH" env-default:"file:///app/migrations"`
}

func New(ctx context.Context, config Config) (*pgxpool.Pool, error) {
	connString := config.GetConnString()
	connString += fmt.Sprintf("&pool_max_conns=%d&pool_min_conns=%d",
		config.MaxConns,
		config.MinConns,
	)

	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return conn, nil
}

func Migrate(ctx context.Context, config Config, migrationsPath string) error {
	connString := config.GetConnString()

	m, err := migrate.New(
		migrationsPath, // "file:///app/db/migrations"
		connString,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %v", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	return nil
}

func (c *Config) GetConnString() string {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
	return connString
}
