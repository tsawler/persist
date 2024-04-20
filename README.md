# Persist

Persist is a simple package that reduces boilerplate code when getting a pool of connections for 
a database.

## Example usage

~~~go
package main

import (
	"github.com/tsawler/persist"
	"log"
)

func main() {
	dsn := "host=localhost port=5433 user=postgres password=password dbname=foo sslmode=disable timezone=UTC connect_timeout=5"

	conn, err := persist.NewPostgres(dsn, nil)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		log.Panic(err)
	}

	log.Println("Connected to db successfully.")
}
~~~