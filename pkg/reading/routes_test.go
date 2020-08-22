package reading

import (
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/aewens/anda/internal/storage"
)

func TestEntity(t *testing.T) {
	configPath, ok := os.LookupEnv("CONFIG_PATH")
	if !ok {
		t.Fatal("Missing 'CONFIG_PATH' environment variable")
	}

	defer Cleanup(t)
	config, err := Config(configPath)
	if err != nil {
		t.Fatal("Could not read config file", err)
	}

	db, err := storage.OpenPostgreSQL(config)
	if err != nil {
		t.Fatal("Could not open database", err)
	}

	defer db.Close()
	entities, err := Entities(db)
	if err != nil {
		t.Fatal("Could not obtain entities", err)
	}

	if len(entities) == 0 {
		t.Error("Did not find any entities")
	}
}
