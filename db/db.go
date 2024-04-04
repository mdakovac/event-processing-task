package db

import (
	"context"
	"log"

	"github.com/Bitstarz-eng/event-processing-challenge/util/env_vars"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DbConnectionPool = pgxpool.Pool

func Connect() *DbConnectionPool {
	var connectionString = env_vars.EnvVariables.DB_CONNECTION_URL
	log.Println("Connecting to DB with connection string:", connectionString)

	conn, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}
