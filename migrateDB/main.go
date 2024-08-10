package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	connStr := "postgres://myuser:mypassword@localhost:5432/mydb"

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		fmt.Println("Error creating connection pool:", err)
		return
	}
	defer pool.Close()

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS transactions (
        id SERIAL PRIMARY KEY,
        type VARCHAR(50) NOT NULL,
        amount DOUBLE PRECISION NOT NULL,
        parent_id INT,
        CONSTRAINT fk_parent
            FOREIGN KEY (parent_id) 
            REFERENCES transactions (id) 
            ON DELETE SET NULL
    );
    `

	_, err = pool.Exec(context.Background(), createTableSQL)
	if err != nil {
		fmt.Println("Error executing SQL command:", err)
		return
	}

	fmt.Println("Table created successfully with foreign key constraint")
}
