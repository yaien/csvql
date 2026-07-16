package csvql

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func TestParse(t *testing.T) {

	tests := []struct {
		name      string
		tablename string
		data      string
		columns   []Column
		rows      int
	}{
		{
			name:      "Basic types",
			tablename: "test",
			data: `id,name,age,active,signup_date,balance
				    1,Alice,30,true,2023-01-15,1000.50
                    2,Bob,25,false,2022-12-20,750.00
                    3,Charlie,35,true,2021-11-05,1250.75`,
			columns: []Column{
				{Name: "id", Type: "VARCHAR"},
				{Name: "name", Type: "VARCHAR"},
				{Name: "age", Type: "VARCHAR"},
				{Name: "active", Type: "VARCHAR"},
				{Name: "signup_date", Type: "VARCHAR"},
				{Name: "balance", Type: "VARCHAR"},
			},
			rows: 3,
		},
		{
			name: "Time and Timestamp",
			data: `event_id,event_name,event_time,event_timestamp
				    1,Login,14:30:00,2023-01-15T14:30:00Z
					2,Logout,15:45:00,2023-01-15T15:45:00Z`,
			tablename: "events",
			columns: []Column{
				{Name: "event_id", Type: "VARCHAR"},
				{Name: "event_name", Type: "VARCHAR"},
				{Name: "event_time", Type: "VARCHAR"},
				{Name: "event_timestamp", Type: "VARCHAR"},
			},
			rows: 2,
		},
		{
			name:      "Mixed types with NULLs",
			tablename: "mixed",
			data: `col1,col2,col3
				    123,Hello,2023-01-01
					NULL,World,2023-02-01
					456,,2023-03-01`,
			columns: []Column{
				{Name: "col1", Type: "VARCHAR"},
				{Name: "col2", Type: "VARCHAR"},
				{Name: "col3", Type: "VARCHAR"},
			},
			rows: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			table, rows, err := Parse(strings.NewReader(test.data), test.tablename)
			if err != nil {
				t.Fatalf("parse error: %s", err.Error())
			}

			if table.Name != test.tablename {
				t.Errorf("expected table name %s, got %s", test.tablename, table.Name)
			}

			if len(table.Columns) != len(test.columns) {
				t.Fatalf("expected %d columns, got %d", len(test.columns), len(table.Columns))
			}

			for index, expected := range test.columns {
				column := table.Columns[index]

				if column.Name != expected.Name {
					t.Errorf("column %d: expected name %s, got %s", index, expected.Name, column.Name)
				}

				if column.Type != expected.Type {
					t.Errorf("column %d (%s): expected type %s, got %s", index, column.Name, expected.Type, column.Type)
				}

			}

			if len(rows) != test.rows {
				t.Errorf("expected %d rows, got %d", test.rows, len(rows))
			}

		})
	}
}

func TestDB(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		tablename string
		rows      int
		columns   int
	}{
		{
			name:      "Simple",
			filename:  "testdata/simple.csv",
			tablename: "sample",
			rows:      4,
			columns:   7,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, err := sql.Open("sqlite", ":memory:")
			if err != nil {
				t.Fatalf("open error: %s", err.Error())
			}

			t.Cleanup(func() {
				db.Close()
			})

			err = ImportOnDB(db, test.filename, test.tablename)
			if err != nil {
				t.Fatalf("read error: %s", err.Error())
			}

			rows, err := db.Query(fmt.Sprintf("select * from %s", test.tablename))
			if err != nil {
				t.Fatalf("query error: %s", err.Error())
			}

			cols, err := rows.Columns()
			if err != nil {
				t.Fatalf("columns error: %s", err.Error())
			}

			if len(cols) != test.columns {
				t.Errorf("expected %d columns, got %d", test.columns, len(cols))
			}

			count := 0
			for rows.Next() {
				record := make([]any, test.columns)
				for i := range record {
					var v any
					record[i] = &v
				}
				err := rows.Scan(record...)
				if err != nil {
					t.Fatalf("scan error: %s", err.Error())
				}

				t.Logf("row %d: %v", count, record)
				count++
			}

			if count != test.rows {
				t.Errorf("expected %d rows, got %d", test.rows, count)
			}
		})
	}
}
