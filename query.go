package csvql

import (
	"database/sql"
)

func Query(db *sql.DB, query string, args ...any) (columns []string, rows [][]any, err error) {
	srows, err := db.Query(query, args...)
	if err != nil {
		return nil, nil, err
	}

	columns, err = srows.Columns()
	if err != nil {
		return nil, nil, err
	}

	for srows.Next() {
		row := make([]any, len(columns))
		pointers := make([]any, len(columns))
		for i := range columns {
			pointers[i] = &row[i]
		}

		if err := srows.Scan(pointers...); err != nil {
			return nil, nil, err
		}

		rows = append(rows, row)
	}

	return columns, rows, nil
}
