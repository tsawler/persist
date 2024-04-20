package persist

import (
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
		{"mysql_fail", "mysql", "foo", nil, false},
		{"mysql_ops", "mysql", "foo", &Options{MaxOpen: 10, MaxIdle: 4, MaxLifetime: 1 * time.Hour}, false},
		{"pg_fail", "postgres", "foo", nil, false},
		{"pg", "postgres", pgTestDatabaseConnectionString, nil, true},
		{"pg_ops", "postgres", pgTestDatabaseConnectionString, &Options{MaxOpen: 10, MaxIdle: 4, MaxLifetime: 1 * time.Hour}, true},
		{"invalid", "banana", "foo", nil, false},
	}

	for _, tt := range tests {
		_, err := New(tt.db, tt.dsn, tt.ops)
		if err != nil && tt.expectSuccess {
			t.Errorf("%s: expected no error but got one", tt.name)
		}
		if err == nil && !tt.expectSuccess {
			t.Errorf("%s: expected error but did not get one", tt.name)
		}
	}
}
