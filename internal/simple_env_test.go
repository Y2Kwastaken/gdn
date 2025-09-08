package internal

import (
	"os"
	"testing"
)

func TestLoadEnvFrom(t *testing.T) {
	err := LoadEnvFrom("simple_env_test.env")
	if err != nil {
		t.Errorf("Failed to load file simple_env_test.env")
	}

	got := os.Getenv("TEST_SECRET")
	expected := "very_secure_secret"

	if got != expected {
		t.Errorf("Expected %s, but got %s", expected, got)
	}
}
