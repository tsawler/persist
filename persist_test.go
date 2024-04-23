package persist

import (
	"os"
	"testing"
	"time"
)

func Test_New(t *testing.T) {
	tests := []struct {
		name          string
		db            string
		dsn           string
		ops           *Options
		expectSuccess bool
	}{
		{"mariadb", "mariadb", mariadbConnString, nil, true},
		{"mariadb_fail", "mysql", "foo", nil, false},
		{"mariadb_ops", "mysql", "foo", &Options{MaxOpen: 10, MaxIdle: 4, MaxLifetime: 1 * time.Hour}, false},
		{"pg_fail", "postgres", "foo", nil, false},
		{"pg", "postgres", pgConnString, nil, true},
		{"pg_ops", "postgres", pgConnString, &Options{MaxOpen: 10, MaxIdle: 4, MaxLifetime: 1 * time.Hour}, true},
		{"invalid", "banana", "foo", nil, false},
	}

	for _, tt := range tests {
		_, err := New(tt.db, tt.dsn, tt.ops)
		if err != nil && tt.expectSuccess {
			t.Errorf("%s: expected no error but got one: %s", tt.name, err.Error())
		}
		if err == nil && !tt.expectSuccess {
			t.Errorf("%s: expected error but did not get one", tt.name)
		}

		if tt.ops != nil {
			if tt.ops.MaxIdle != 4 {
				t.Errorf("%s: wrong value for MaxIdle; expected 4 but got %d", tt.name, tt.ops.MaxIdle)
			}
			if tt.ops.MaxOpen != 10 {
				t.Errorf("%s: wrong value for MaxOpen; expected 10 but got %d", tt.name, tt.ops.MaxOpen)
			}
			if tt.ops.MaxLifetime != 1*time.Hour {
				t.Errorf("%s: wrong value for MaxLifetime; expected 1 hour but got %d", tt.name, tt.ops.MaxLifetime)
			}
		}
	}
}

func Test_NewPostgres(t *testing.T) {
	tests := []struct {
		name          string
		dsn           string
		ops           *Options
		expectSuccess bool
	}{
		{"pg_fail", "foo", nil, false},
		{"pg", pgConnString, nil, true},
		{"pg_ops", pgConnString, &Options{MaxOpen: 10, MaxIdle: 4, MaxLifetime: 1 * time.Hour}, true}}

	for _, tt := range tests {
		_, err := NewPostgres(tt.dsn, tt.ops)
		if err != nil && tt.expectSuccess {
			t.Errorf("%s: expected no error but got one: %s", tt.name, err.Error())
		}
		if err == nil && !tt.expectSuccess {
			t.Errorf("%s: expected error but did not get one", tt.name)
		}
	}
}

func Test_NewMariaDB(t *testing.T) {
	tests := []struct {
		name          string
		dsn           string
		ops           *Options
		expectSuccess bool
	}{
		{"mariadb", mariadbConnString, nil, true},
		{"mariadb_fail", "foo", nil, false},
		{"mariadb_ops", "foo", &Options{MaxOpen: 10, MaxIdle: 4, MaxLifetime: 1 * time.Hour}, false},
	}

	for _, tt := range tests {
		_, err := NewMariaDB(tt.dsn, tt.ops)
		if err != nil && tt.expectSuccess {
			t.Errorf("%s: expected no error but got one: %s", tt.name, err.Error())
		}
		if err == nil && !tt.expectSuccess {
			t.Errorf("%s: expected error but did not get one", tt.name)
		}
	}
}

func Test_NewSQLite(t *testing.T) {
	tests := []struct {
		name          string
		dsn           string
		ops           *Options
		expectSuccess bool
	}{
		{"sqlite_memory", ":memory:", nil, true},
		{"sqlite_disk", "./testdata/test.db", nil, true},
	}

	for _, tt := range tests {
		_, err := NewSQLite(tt.dsn, nil)
		if err != nil && tt.expectSuccess {
			t.Errorf("%s: expected no error but got one: %s", tt.name, err.Error())
		}
		if err == nil && !tt.expectSuccess {
			t.Errorf("%s: expected error but did not get one", tt.name)
		}
	}

	defer func() {
		err := os.Remove("./testdata/test.db")
		if err != nil {
			// do nothing
		}
	}()
}

func TestBuildConnectionString(t *testing.T) {
	tests := []struct {
		name          string
		dbType        string
		host          string
		user          string
		pass          string
		ssl           string
		db            string
		port          int
		expected      string
		expectSuccess bool
	}{
		{"pg", "pg", "localhost", "user", "pass", "foo", "disable", 5432, "host=localhost port=5432 user=user password=pass dbname=foo sslmode=disable", true},
		{"mysql", "mysql", "localhost", "user", "pass", "foo", "false", 3306, "user:pass@tcp(localhost:3306)/foo?parseTime=true&tls=false&collation=utf8_unicode_ci&timeout=5s&readTimeout=5s", true},
		{"bad", "fish", "localhost", "user", "pass", "foo", "false", 3306, "", false},
	}

	for _, tt := range tests {
		cd := ConnectionData{
			DBType:   tt.dbType,
			UserName: tt.user,
			Password: tt.pass,
			Host:     tt.host,
			Database: tt.db,
			SSL:      tt.ssl,
			Port:     tt.port,
		}
		result, err := BuildConnectionString(cd)
		if err != nil && tt.expectSuccess {
			t.Errorf("%s: expected no error but got one: %s", tt.name, err.Error())
		}
		if err == nil && !tt.expectSuccess {
			t.Errorf("%s: expected error but did not get one", tt.name)
		}
		if tt.expectSuccess {
			if tt.expected != tt.expected {
				t.Errorf("%s: expected %s but got %s", tt.name, tt.expected, result)
			}
		}
	}
}
