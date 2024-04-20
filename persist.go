package persist

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"strings"
	"time"
)

var (
	MaxOpenConns    = 12              // Maximum number of open connections in our pool.
	MaxIdleConns    = 6               // Maximum idle connections in our pool.
	ConnMaxLifetime = 0 * time.Second // Max lifetime for a connection (how long before it expires). 0 is forever.
)

// Options holds useful options for a pool of connections.
type Options struct {
	MaxOpen     int
	MaxIdle     int
	MaxLifetime time.Duration
}

// New is a factory method which takes a db type, a dsn and options and attempts
// to open a connection to the database and return a pool of connections.
func New(db, dsn string, ops *Options) (*sql.DB, error) {
	engine := strings.ToLower(db)

	if ops != nil {
		if ops.MaxLifetime != 0 {
			ConnMaxLifetime = ops.MaxLifetime
		}
		if ops.MaxOpen != 0 {
			MaxOpenConns = ops.MaxOpen
		}
		if ops.MaxIdle != 0 {
			MaxIdleConns = ops.MaxIdle
		}
	}

	switch engine {
	case "mysql", "mariadb":
		return connectToMySQL(dsn)
	case "postgres", "pg", "postgresql":
		return connectToPostgres(dsn)
	default:
		return nil, errors.New("invalid database engine supplied")
	}
}

// connectToMySQL attempts to get a pool of connections for a MySQL/MariaDB database.
func connectToMySQL(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// test ping
	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(MaxOpenConns)
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetConnMaxLifetime(ConnMaxLifetime)

	return db, nil
}

// connectToPostgres attempts to get a pool of connections for a postgres database.
func connectToPostgres(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(MaxOpenConns)
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetConnMaxLifetime(ConnMaxLifetime)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
