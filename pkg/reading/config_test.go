package reading

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	configPath, ok := os.LookupEnv("CONFIG_PATH")
	if !ok {
		t.Fatal("Missing 'CONFIG_PATH' environment variable")
	}

	defer Cleanup(t)
	_, err := Config(configPath)
	if err != nil {
		t.Fatal("Could not read config file", err)
	}
}
