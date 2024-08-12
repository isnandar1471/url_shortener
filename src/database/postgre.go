package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

func MakeConnection() *pgx.Conn {
	connection, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database %v\n", err)
		os.Exit(1)
	}
	return connection
}

func CloseConnection(conn *pgx.Conn) {
	conn.Close(context.Background())
}

func Select(conn *pgx.Conn) (error, string) {
	var result string
	rows, err := conn.Query(context.Background(), "SELECT * FROM short")
	if err != nil {
		return err, ""
	}

	rows.Scan()

	return nil, result
}
