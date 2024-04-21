package persist

import (
	"database/sql"
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"os"
	"testing"

	_ "github.com/glebarez/go-sqlite"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var pgConnString = "host=localhost port=5433 user=postgres password=password dbname=foo sslmode=disable"
var mariadbConnString = "mariadb:password@tcp(localhost:3307)/foo?parseTime=true&tls=false&collation=utf8_unicode_ci&timeout=5s&readTimeout5"

//func TestMain(m *testing.M) {
//	code := m.Run()
//	os.Exit(code)
//}

var (
	host     = "localhost"
	user     = "postgres"
	password = "password"
	dbName   = "foo"
	port     = "5433"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
)

var resource *dockertest.Resource

func TestMain(m *testing.M) {
	pgResource, pgPool := postgresUp()
	mariadbResource, mariadbPool := mariadbUp()

	// run tests
	code := m.Run()

	// clean up
	if err := pgPool.Purge(pgResource); err != nil {
		log.Fatalf("could not purge pg resource: %s", err)
	}

	if err := mariadbPool.Purge(mariadbResource); err != nil {
		log.Fatalf("could not purge mariadb resource: %s", err)
	}

	os.Exit(code)
}

func postgresUp() (*dockertest.Resource, *dockertest.Pool) {
	// connect to docker; fail if docker not running
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker; is it running? %s", err)
	}

	// set up our docker options, specifying the image and so forth
	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.5",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	// get a resource (docker image)
	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	// start the image and wait until it's ready
	if err := pool.Retry(func() error {
		var err error
		testDB, err := sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println("Error:", err)
			return err
		}
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to pg database: %s", err)
	}

	return resource, pool
}

func mariadbUp() (*dockertest.Resource, *dockertest.Pool) {
	// connect to docker; fail if docker not running
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker; is it running? %s", err)
	}

	// set up our docker options, specifying the image and so forth
	opts := dockertest.RunOptions{
		Repository: "mariadb",
		Tag:        "10.6",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=" + password,
			"MYSQL_USER=mariadb",
			"MYSQL_PASSWORD=" + password,
			"MYSQL_DATABASE=" + dbName,
		},
		ExposedPorts: []string{"3306"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"3306": {
				{HostIP: "0.0.0.0", HostPort: "3307"},
			},
		},
	}

	// get a resource (docker image)
	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	// start the image and wait until it's ready
	if err := pool.Retry(func() error {
		var err error
		testDB, err := sql.Open("mysql", mariadbConnString)
		if err != nil {
			return err
		}
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to mariadb database: %s", err)
	}

	return resource, pool
}
