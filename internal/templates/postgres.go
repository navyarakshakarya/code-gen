package templates

// PostgreSQL Database Template
const PostgresTemplate = `package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewConnection(dsn string) (*Postgres, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("Postgres ping failed: %w", err)
	}

	log.Println("Postgres pool connected successfully")
	return &Postgres{pool: pool}, nil
}

func (pg *Postgres) DB() *pgxpool.Pool {
	return pg.pool
}

func (pg *Postgres) Close() {
	log.Println("Closing Postgres pool...")
	pg.pool.Close()
}


`

// WithConnSchema acquires a connection, sets the search_path, runs fn, and resets search_path
const WithConnSchema = `package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

// WithSchema runs a function using a connection with tenant-specific search_path
func (pg *Postgres) WithSchema(ctx context.Context, schema string, fn func(conn *pgx.Conn) error) error {
	conn, err := pg.pool.Acquire(ctx)
	if err != nil {
		log.Printf("error happened while trying to acquire pool connection, err: %v", err)
		return errors.New("failed to acquire pool connection")
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, fmt.Sprintf("SET search_path TO \"%s\", public", schema))
	if err != nil {
		log.Printf("error happened while trying to set search path, err: %v", err)
		return errors.New("failed to set search path")
	}

	return fn(conn.Conn())
}


`
