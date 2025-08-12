package utils

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

func OpenDB(dsn string, setLimits bool) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	if setLimits {
		log.Printf("Setting connection limits")
		db.Config().MaxConns = 50
		db.Config().MinConns = 1
		db.Config().MaxConnLifetime = 10 * time.Minute
		db.Config().MaxConnIdleTime = 1 * time.Minute
		db.Config().HealthCheckPeriod = 1 * time.Minute
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
