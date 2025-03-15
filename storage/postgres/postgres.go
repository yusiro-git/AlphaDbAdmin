package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

func New(postgresUrl string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), postgresUrl)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %v", err)
	}
	return conn, nil
}
