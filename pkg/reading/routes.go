package reading

import (
	"database/sql"

	"github.com/aewens/anda/pkg/core"
)

func Entities(db *sql.DB) ([]core.SQLEntry, error) {
	var entries []core.SQLEntry

	query := `
		SELECT e.uuid, a.name, vt.name, convert_from(v.value, 'utf-8'), v.flag
		FROM entity e
		INNER JOIN value v ON v.entity_id = e.id
		INNER JOIN attribute a ON a.id = v.attribute_id
		INNER JOIN value_type vt ON vt.id = a.value_type_id;
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			uuid  string
			name  string
			vtype string
			value interface{}
			flag  int
		)
		err := rows.Scan(&uuid, &name, &vtype, &value, &flag)
		if err != nil {
			return nil, err
		}

		entry := core.SQLEntry{
			UUID:  uuid,
			Name:  name,
			Type:  vtype,
			Value: value,
			Flag:  flag,
		}
		entries = append(entries, entry)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return entries, nil
}
