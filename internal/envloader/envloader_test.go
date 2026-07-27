package envloader

import (
	"strings"
	"testing"
)

const envPath = "../../.env"

func TestUnrecognizedEnvironment(t *testing.T) {
	t.Setenv("ENV", "FAIL")
	_, err := LoadEnvironment(envPath)
	if err == nil {
		t.Errorf("expected error, didn't get one")
		return
	}
}

func TestDevEnvironment(t *testing.T) {
	t.Setenv("ENV", "DEV")
	envSettings, err := LoadEnvironment(envPath)
	if err != nil {
		t.Errorf("didn't expect error, got %v\n", err)
		return
	}
	if !strings.Contains(envSettings.GetDBString(), "localhost:5432") {
		t.Errorf("wrong dbString for dev environment: %s\n", envSettings.GetDBString())
	}
}
