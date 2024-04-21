# Persist

Persist is a simple package that reduces boilerplate code when getting a pool of connections for 
a database.

This project does nothing particularly spectacular; it exists simply to save me writing the same 100 lines of code
every time I need to connect to a database.

Currently, this project supports Postgres, MariaDB/MySQL, and sqlite.

## Example usage

~~~go
package main

import (
	"github.com/tsawler/persist"
	"log"
)

func main() {
	// Specify your connection string.
	dsn := "host=localhost port=5433 user=postgres password=password dbname=foo sslmode=disable"

	// Get a pool of connections. The first parameter can be "postgres", "mariadb", "mysql", or "sqlite".
	conn, err := persist.New("postgres", dsn, nil)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	log.Println("Connected to db successfully.")
}
~~~