package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/aewens/anda-server/pkg/core"
)

func OpenPostgreSQL(config *core.Config) (*sql.DB, error) {
	var conn []string

	conn = append(conn, fmt.Sprintf("host=%s", config.DBHost))
	conn = append(conn, fmt.Sprintf("port=%s", config.DBPort))
	conn = append(conn, fmt.Sprintf("user=%s", config.DBUser))
	conn = append(conn, fmt.Sprintf("password=%s", config.DBPswd))
	conn = append(conn, fmt.Sprintf("dbname=%s", config.DBName))
	conn = append(conn, "sslmode=disable")

	db, err := sql.Open("postgres", strings.Join(conn, " "))
	return db, err
}
