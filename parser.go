package csvql

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
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

type ParseOptions func(*csv.Reader)

func WithComma(comma string) ParseOptions {
	return func(r *csv.Reader) {
		r.Comma = []rune(comma)[0]
	}
}

// Parse reads CSV data from the provided reader and infers column types.
func Parse(reader io.Reader, name string, opts ...ParseOptions) (*Schema, [][]string, error) {
	r := csv.NewReader(reader)
	for _, opt := range opts {
		opt(r)
	}

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

	return &Schema{Name: name, Columns: columns}, records[1:], nil
}

// ParseFile reads CSV data from the specified file and infers column types.
func ParseFile(filename string, tablename string, opts ...ParseOptions) (*Schema, [][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}

	defer func() { _ = file.Close() }()

	return Parse(file, tablename, opts...)
}

// ImportOnDB creates an db from the provided Table structure.
func ImportOnDB(db *sql.DB, filename, tablename string, opts ...ParseOptions) error {
	table, rows, err := ParseFile(filename, tablename, opts...)
	if err != nil {
		return err
	}

	// Create table with columns in one statement
	var create strings.Builder
	fmt.Fprintf(&create, "CREATE TABLE %s (", table.Name)
	for index, column := range table.Columns {
		if index > 0 {
			fmt.Fprint(&create, ", ")
		}
		fmt.Fprintf(&create, "%s %s", column.Name, column.Type)
	}
	fmt.Fprint(&create, ")")

	_, err = db.Exec(create.String())
	if err != nil {
		return fmt.Errorf("create table error: %w", err)
	}

	// Insert rows
	for _, row := range rows {
		var insert strings.Builder
		fmt.Fprintf(&insert, "INSERT INTO %s VALUES (", table.Name)

		values := make([]any, len(row))
		for index, value := range row {
			fmt.Fprint(&insert, "?")
			if index < len(row)-1 {
				fmt.Fprint(&insert, ",")
			}
			values[index] = value
		}
		fmt.Fprint(&insert, ")")

		_, err = db.Exec(insert.String(), values...)
		if err != nil {
			return err
		}
	}

	return nil
}
