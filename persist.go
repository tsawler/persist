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
	MaxOpenDbConn = 25
	MaxIdleDbConn = 25
	MaxDbLifetime = 5 * time.Minute
)

type Options struct {
	MaxOpen     int
	MaxIdle     int
	MaxLifetime time.Duration
}

type Database struct {
	Engine *sql.DB
}

func New(db, dsn string, ops *Options) (*Database, error) {
	engine := strings.ToLower(db)

	if ops != nil {
		if ops.MaxLifetime != 0 {
			MaxDbLifetime = ops.MaxLifetime
		}
		if ops.MaxOpen != 0 {
			MaxOpenDbConn = ops.MaxOpen
		}
		if ops.MaxIdle != 0 {
			MaxIdleDbConn = ops.MaxIdle
		}
	}

	switch engine {
	case "mysql", "mariadb":
		return initMySQLDB(dsn)
	case "postgres", "pg", "postgresql":
		return initPostgresDB(dsn)
	default:
		return nil, errors.New("invalid database engine supplied")
	}
}

func initMySQLDB(dsn string) (*Database, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// test ping
	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(MaxOpenDbConn)
	db.SetMaxIdleConns(MaxIdleDbConn)
	db.SetConnMaxLifetime(MaxDbLifetime)

	return &Database{Engine: db}, nil
}

func initPostgresDB(dsn string) (*Database, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(MaxOpenDbConn)
	db.SetMaxIdleConns(MaxIdleDbConn)
	db.SetConnMaxLifetime(MaxDbLifetime)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Database{Engine: db}, nil
}
