package persist

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"os"
	"testing"
)

var (
	host                  = "localhost"
	port                  = "5435"
	mysqlPort             = "3307"
	user                  = "username"
	password              = "password"
	dbName                = "foo"
	pgConnectionString    = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
	mysqlConnectionString = "%s:%s@tcp(%s:%s)/%s?parseTime=true&tls=false"
)

var pgTestDatabaseConnectionString string
var pgResource *dockertest.Resource
var pgPool *dockertest.Pool

var mysqlResource *dockertest.Resource
var mysqlPool *dockertest.Pool

func TestMain(m *testing.M) {

	spinUpPostgres()

	//spinUpMySQL()

	// run tests
	code := m.Run()

	// clean up
	if err := pgPool.Purge(pgResource); err != nil {
		log.Fatalf("could not purge pgResource: %s", err)
	}
	//if err := mysqlPool.Purge(mysqlResource); err != nil {
	//	log.Fatalf("could not purge mysqlResource: %s", err)
	//}

	os.Exit(code)
}

func spinUpMySQL() {

}

func spinUpPostgres() {
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker; is it running? %s", err)
	}
	pgPool = p
	// Set up Postgres container.
	pgTestDatabaseConnectionString = fmt.Sprintf(pgConnectionString, host, port, user, password, dbName)

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

	// get a pgResource (docker image)
	pgResource, err = pgPool.RunWithOptions(&opts, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		_ = pgPool.Purge(pgResource)
		log.Fatalf("could not start pgResource: %s", err)
	}

	// start the image and wait until it's ready
	if err = pgPool.Retry(func() error {
		var err error
		testDB, err := sql.Open("pgx", fmt.Sprintf(pgConnectionString, host, port, user, password, dbName))
		if err != nil {
			log.Println("Error:", err)
			return err
		}
		return testDB.Ping()
	}); err != nil {
		_ = pgPool.Purge(pgResource)
		log.Fatalf("could not connect to database: %s", err)
	}
}
