package testutil

import (
	"os"
	"strings"
)

// ClearArcticEnvVars removes all ARCTIC_* environment variables
// for complete test isolation. Use this in TestMain to ensure
// tests start with a clean environment.
func ClearArcticEnvVars() {
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "ARCTIC_") {
			key := strings.Split(env, "=")[0]
			os.Unsetenv(key)
		}
	}
}

// SetupTestDatasource sets up default datasource environment variables
// for testing. These values should work with test containers or mock databases.
func SetupTestDatasource() {
	os.Setenv("ARCTIC_DATASOURCE_USERNAME", "testuser")
	os.Setenv("ARCTIC_DATASOURCE_PASSWORD", "testpass")
	os.Setenv("ARCTIC_DATASOURCE_HOST", "localhost")
	os.Setenv("ARCTIC_DATASOURCE_PORT", "5432")
	os.Setenv("ARCTIC_DATASOURCE_DBNAME", "testdb")
	os.Setenv("ARCTIC_DATASOURCE_SSLMODE", "disable")
}
