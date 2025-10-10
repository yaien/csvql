package csvql

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/marcboeker/go-duckdb"
)

var DateFormats = []string{
	"2006-01-02",
	"02/01/2006",
	"01/02/2006",
	"2006/01/02",
	"02-01-2006",
	"01-02-2006",
}

var TimestampFormats = []string{
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02 15:04:05",
	"02/01/2006 15:04:05",
	"01/02/2006 15:04:05",
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05",
	"2006/01/02 15:04:05",
}

var TimeFormats = []string{
	"15:04:05",
	"3:04PM",
	"3:04:05PM",
	"15:04",
	"3:04",
}

type Column struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Format string `json:"format,omitempty"`
}

type Schema struct {
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`
}

type Parser struct {
	Name     string
	Validate func(v string) (value any, format string, err error)
}

type Parsers []*Parser

func (p Parsers) Get(name string) (*Parser, bool) {
	for _, parser := range p {
		if parser.Name == name {
			return parser, true
		}
	}
	return nil, false
}

var parsers = Parsers{
	{
		Name: "INTEGER",
		Validate: func(v string) (any, string, error) {
			n, err := strconv.ParseInt(v, 10, 64)
			return n, "", err
		},
	},
	{
		Name: "DOUBLE",
		Validate: func(v string) (any, string, error) {
			n, err := strconv.ParseFloat(v, 64)
			return n, "", err
		},
	},
	{
		Name: "BOOLEAN",
		Validate: func(v string) (any, string, error) {
			switch strings.ToLower(v) {
			case "true":
				return true, "", nil
			case "false":
				return false, "", nil
			default:
				return nil, "", fmt.Errorf("invalid boolean value: %s", v)
			}
		},
	},
	{
		Name: "DATE",
		Validate: func(v string) (any, string, error) {
			for _, format := range DateFormats {
				if t, err := time.Parse(format, v); err == nil {
					return t, format, nil
				}
			}
			return nil, "", fmt.Errorf("invalid DATE value: %s", v)
		},
	},
	{
		Name: "TIME",
		Validate: func(v string) (any, string, error) {
			for _, format := range TimeFormats {
				if t, err := time.Parse(format, v); err == nil {
					return t, format, nil
				}
			}
			return nil, "", fmt.Errorf("invalid TIME value: %s", v)
		},
	},
	{
		Name: "TIMESTAMP",
		Validate: func(v string) (any, string, error) {
			for _, format := range TimestampFormats {
				if t, err := time.Parse(format, v); err == nil {
					return t, format, nil
				}
			}
			return nil, "", fmt.Errorf("invalid TIMESTAMP value: %s", v)
		},
	},
	{
		Name: "VARCHAR",
		Validate: func(v string) (any, string, error) {
			return v, "", nil
		},
	},
}

// Parse reads CSV data from the provided reader and infers column types.
func Parse(reader io.Reader, name string) (*Schema, [][]string, error) {
	r := csv.NewReader(reader)
	records, err := r.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	if len(records) < 2 {
		return nil, nil, fmt.Errorf("not enough records")
	}

	columns := make([]Column, len(records[0]))
	for index, name := range records[0] {
		name = strings.TrimSpace(name)
		name = strings.ReplaceAll(name, " ", "_")
		name = strings.ToLower(name)
		columns[index] = Column{Name: name, Type: "VARCHAR"}
	}

	// Analyze first few rows to determine column types
	for index, v := range records[1] {
		v = strings.TrimSpace(v)
		if v == "" || v == "NULL" {
			continue
		}

		for _, parser := range parsers {
			if _, format, err := parser.Validate(v); err == nil {
				columns[index].Type = parser.Name
				columns[index].Format = format
				break
			}
		}
	}

	return &Schema{Name: name, Columns: columns}, records[1:], nil
}

// ParseFile reads CSV data from the specified file and infers column types.
func ParseFile(filename string, tablename string) (*Schema, [][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}

	defer func() { _ = file.Close() }()

	return Parse(file, tablename)
}

func CreateDB(path string) (*sql.DB, error) {
	db, err := sql.Open("duckdb", path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// ImportOnDB creates an in-memory DuckDB database from the provided Table structure.
func ImportOnDB(db *sql.DB, table *Schema, rows [][]string) error {

	// Create table with columns in one statement (DuckDB style)
	sentece := fmt.Sprintf("CREATE TABLE %s (", table.Name)
	for index, column := range table.Columns {
		if index > 0 {
			sentece += ", "
		}
		sentece += fmt.Sprintf("%s %s", column.Name, column.Type)
	}
	sentece += ")"

	_, err := db.Exec(sentece)
	if err != nil {
		return err
	}

	// Insert rows
	for _, row := range rows {
		placeholders := ""
		values := make([]any, len(row))
		for index, value := range row {
			if index > 0 {
				placeholders += ", "
			}
			placeholders += "?"

			column := table.Columns[index]

			parser, ok := parsers.Get(column.Type)
			if !ok {
				return fmt.Errorf("unknown parser for type: %s", column.Type)
			}

			v, _, err := parser.Validate(value)
			if err != nil {
				return fmt.Errorf("failed to parse value '%s' for column '%s' of type '%s': %s", row[index], column.Name, column.Type, err.Error())
			}

			values[index] = v
		}
		_, err := db.Exec(fmt.Sprintf("INSERT INTO %s VALUES (%s);", table.Name, placeholders), values...)
		if err != nil {
			return err
		}
	}

	return nil
}

// Read reads a CSV file, infers its structure, and creates an in-memory DuckDB database.
func Read(filename, tablename string) (*sql.DB, *Schema, error) {
	schema, rows, err := ParseFile(filename, tablename)
	if err != nil {
		return nil, nil, fmt.Errorf("parse error: %s", err.Error())
	}

	db, err := CreateDB("")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create database: %s", err.Error())
	}

	err = ImportOnDB(db, schema, rows)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to import database: %s", err.Error())
	}

	return db, schema, nil
}
