package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type DB struct {
	Pool *pgxpool.Pool
}

func New(config Config) (*DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.User, config.Password, config.Host, config.Port, config.Database)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PostgreSQL connection string: %w", err)
	}

	poolConfig.ConnConfig.Logger = &qLogger{}

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to TimescaleDB: %w", err)
	}

	db := &DB{Pool: pool}
	return db, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}

type qLogger struct {
}

func (l *qLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	fmt.Printf("SQL:\n%s\nARGS:%v\n", data["sql"], data["args"])
}
