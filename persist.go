package persist

import (
	"database/sql"
	"errors"
	_ "github.com/glebarez/go-sqlite"
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

// NewPostgres is a convenience function for getting a pool of Postgres connections.
func NewPostgres(dsn string, ops *Options) (*sql.DB, error) {
	return New("pg", dsn, ops)
}

// NewMariaDB is a convenience function for getting a pool of MariaDB/MySQL connections.
func NewMariaDB(dsn string, ops *Options) (*sql.DB, error) {
	return New("mariadb", dsn, ops)
}

// NewSQLite is a convenience function for getting a pool of sqlite connections.
func NewSQLite(dsn string, ops *Options) (*sql.DB, error) {
	return New("sqlite", dsn, ops)
}

// New is a factory method which takes a db type, a DSN and options and attempts
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
		return connect("mysql", dsn)
	case "postgres", "pg", "postgresql":
		return connect("pgx", dsn)
	case "sqlite":
		return connect("sqlite", dsn)
	default:
		return nil, errors.New("invalid database engine supplied")
	}
}

// connect attempts to get a pool of connections for a given database
// with the DSN supplied as dsn.
// Note that for sqlite, dsn can be one of ":memory:" for in memory, or
// "path/to/some.db" for disk based storage. In order to avoid "database
// is locked errors" you must set MaxOpenConns to 1 so that only 1
// connection is ever used by the DB, allowing concurrent access to
// DB without making the writes concurrent
func connect(driver, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
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
