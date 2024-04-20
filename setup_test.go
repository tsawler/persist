package persist

import (
	"os"
	"testing"
)

var pgConnString = "host=localhost port=5433 user=postgres password=password dbname=foo sslmode=disable timezone=UTC connect_timeout=5"
var mariadbConnString = "mariadb:password@tcp(localhost:3307)/foo"

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
