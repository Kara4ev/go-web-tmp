package app

import (
	"fmt"

	"github.com/Kara4ev/go-web-tmp/internal/config"
	"github.com/Kara4ev/go-web-tmp/pkg/logger"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type dataSource struct {
	DB    *sqlx.DB
	Radis *redis.Client
}

func initDS(cfg *config.Config) (*dataSource, error) {
	logger.Debug("initializing data sorce")

	// postgres sql
	pgConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PGHost,
		cfg.PGPort,
		cfg.PGUser,
		cfg.PGPassword,
		cfg.PGDB,
		cfg.PGSSL)

	logger.Debug("connect to postgres sql")
	db, err := sqlx.Open("postgres", pgConnString)

	if err != nil {
		logger.Debug("error open db: %w", err)
		return nil, fmt.Errorf("error open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		logger.Debug("error connecting to db: %w", err)
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	logger.Debug("connect to redis")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RDHost, cfg.RDPort),
		Password: cfg.RDPassword,
		DB:       cfg.RDdb,
	})

	logger.Debug("data sorce initializing")
	return &dataSource{
		DB:    db,
		Radis: rdb,
	}, nil
}

func (d *dataSource) close() error {
	if err := d.DB.Close(); err != nil {
		return fmt.Errorf("error closed postgres connect: %w", err)
	}

	if err := d.Radis.Close(); err != nil {
		return fmt.Errorf("error closed radis connect: %w", err)
	}

	return nil
}
