package app

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"os"
	"time"
)

func connectToDB(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		),
	)

	if err != nil {
		return nil, errors.Wrap(err, "ConnectToDB #1")
	}

	timeout, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	if err = db.PingContext(timeout); err != nil {
		return nil, errors.Wrap(err, "ConnectToDB #2")
	}

	return db, nil
}
