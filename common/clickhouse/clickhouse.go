package clickhouse

import (
	"context"
	"errors"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

type Config struct {
	Host         string `yaml:"host" env:"HOST" env-default:"clickhouse"`
	Port         uint16 `yaml:"port" env:"PORT" env-default:"9000"`
	Username     string `yaml:"user" env:"USER" env-default:"default"`
	Password     string `yaml:"password" env:"PASSWORD"`
	Database     string `yaml:"db" env:"DB" env-default:"default"`
	MaxOpenConns int    `yaml:"max_open_conns" env:"MAX_OPEN_CONNS" env-default:"10"`
	MaxIdleConns int    `yaml:"max_idle_conns" env:"MAX_IDLE_CONNS" env-default:"10"`
}

func New(ctx context.Context, cfg Config) (driver.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.Username,
			Password: cfg.Password,
		},
		Debug:        true,
		MaxOpenConns: cfg.MaxOpenConns,
		MaxIdleConns: cfg.MaxIdleConns,
	})
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(ctx); err != nil {
		log.Fatal("clickhouse not available:", err)
	}
	log.Printf("connected to ClickHouse at %s:%d 🎉", cfg.Host, cfg.Port)
	return conn, nil
}

func Migrate(config Config, migrationsPath string) error {
	connString := config.GetConnString()
	m, err := migrate.New(
		migrationsPath,
		connString,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %v", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate database: %v", err)
	}
	log.Println("migrated successfully")
	return nil
}

func (c *Config) GetConnString() string {
	connString := fmt.Sprintf(
		"clickhouse://%s:%s@%s:%d/%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
	return connString
}
