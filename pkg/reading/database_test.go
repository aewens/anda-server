package reading

import (
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/aewens/anda-server/internal/storage"
)

func TestEntries(t *testing.T) {
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
	entries, err := Entries(db)
	if err != nil {
		t.Fatal("Could not obtain entries", err)
	}

	if len(entries) == 0 {
		t.Error("Did not find any entries")
	}
}
