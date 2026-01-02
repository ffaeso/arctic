//go:build unit

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDSN(t *testing.T) {
	tests := []struct {
		name     string
		config   DatasourceConfig
		expected string
	}{
		{
			name: "generic config",
			config: DatasourceConfig{
				Username: "arctic",
				Password: "changeme",
				Host:     "localhost",
				Port:     5554,
				DbName:   "arctic_dev_db",
				Sslmode:  "require",
			},
			expected: "postgres://arctic:changeme@localhost:5554/arctic_dev_db?sslmode=require",
		},
		{
			name: "different ssl mode",
			config: DatasourceConfig{
				Username: "user",
				Password: "pass",
				Host:     "db.example.com",
				Port:     5433,
				DbName:   "mydb",
				Sslmode:  "disable",
			},
			expected: "postgres://user:pass@db.example.com:5433/mydb?sslmode=disable",
		},
		{
			name: "special characters in password",
			config: DatasourceConfig{
				Username: "admin",
				Password: "p@ssw0rd!",
				Host:     "127.0.0.1",
				Port:     5432,
				DbName:   "prod",
				Sslmode:  "verify-full",
			},
			expected: "postgres://admin:p@ssw0rd!@127.0.0.1:5432/prod?sslmode=verify-full",
		},
		{
			name: "password with URL-unsafe characters",
			config: DatasourceConfig{
				Username: "app_user",
				Password: "p@ss:w/rd#123",
				Host:     "db.internal",
				Port:     5432,
				DbName:   "app_db",
				Sslmode:  "require",
			},
			expected: "postgres://app_user:p@ss:w/rd#123@db.internal:5432/app_db?sslmode=require",
		},
		{
			name: "ipv6 host",
			config: DatasourceConfig{
				Username: "user",
				Password: "pass",
				Host:     "::1",
				Port:     5432,
				DbName:   "testdb",
				Sslmode:  "disable",
			},
			expected: "postgres://user:pass@::1:5432/testdb?sslmode=disable",
		},
		{
			name: "ipv6 host with brackets",
			config: DatasourceConfig{
				Username: "user",
				Password: "pass",
				Host:     "[2001:db8::1]",
				Port:     5432,
				DbName:   "testdb",
				Sslmode:  "require",
			},
			expected: "postgres://user:pass@[2001:db8::1]:5432/testdb?sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.config.DSN()
			expected := tt.expected
			assert.Equal(t, expected, actual)
		})
	}
}
