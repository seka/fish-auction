package testing

import (
	"os"
	"testing"
)

func requireIntegrationTests(t *testing.T) {
	t.Helper()

	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test unless RUN_INTEGRATION_TESTS=true")
	}
}

func requireStressTests(t *testing.T) {
	t.Helper()

	if os.Getenv("RUN_STRESS_TESTS") != "true" {
		t.Skip("Skipping stress test unless RUN_STRESS_TESTS=true")
	}
}
