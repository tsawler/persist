package persist

import (
	"database/sql"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"os"
	"testing"
)

var pgConnString = "host=localhost port=5433 user=postgres password=password dbname=foo sslmode=disable"
var mariadbConnString = "mariadb:password@tcp(localhost:3307)/foo?parseTime=true&tls=false&collation=utf8_unicode_ci&timeout=5s&readTimeout5"

var (
	user        = "postgres"
	mariadbUser = "mariadb"
	password    = "password"
	dbName      = "foo"
	pgPort      = "5433"
	mariadbPort = "3307"
)

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
				{HostIP: "0.0.0.0", HostPort: pgPort},
			},
		},
	}

	// get a resource (docker image)
	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	// start the image and wait until it's ready
	if err := pool.Retry(func() error {
		testDB, retryErr := sql.Open("pgx", pgConnString)
		if retryErr != nil {
			log.Println("Error:", retryErr)
			return retryErr
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
			"MYSQL_USER=" + mariadbUser,
			"MYSQL_PASSWORD=" + password,
			"MYSQL_DATABASE=" + dbName,
		},
		ExposedPorts: []string{"3306"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"3306": {
				{HostIP: "0.0.0.0", HostPort: mariadbPort},
			},
		},
	}

	// get a resource (docker image)
	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	// start the image and wait until it's ready
	if err := pool.Retry(func() error {
		testDB, retryErr := sql.Open("mysql", mariadbConnString)
		if retryErr != nil {
			return retryErr
		}
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to mariadb database: %s", err)
	}

	return resource, pool
}
