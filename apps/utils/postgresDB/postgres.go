package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Db *pgxpool.Pool
}

func NewPostgres() (*Postgres, error) {
	connStr := "postgres://myuser:mypassword@localhost:5432/mydb" // ideally should get from env variable or secret store
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		fmt.Println("Error creating connection pool:", err)
		return nil, err
	}
	return &Postgres{Db: pool}, nil
}
