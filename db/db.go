package db

import (
	"context"
	"fmt"
	"log"

	"github.com/Bitstarz-eng/event-processing-challenge/util/env_vars"
	"github.com/jackc/pgx/v5"
)

type DbConnection = pgx.Conn

func Connect() *DbConnection {
	var connectionString = env_vars.EnvVariables.DB_CONNECTION_URL
	fmt.Println("Connecting to DB with connection string:", connectionString)

	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}

func Disconnect(c *DbConnection) {
	c.Close(context.Background())
}
