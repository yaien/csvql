package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yaien/csvql"
)

func root() *cobra.Command {
	var tablename string
	var query string

	cmd := &cobra.Command{
		Use:   "csvql [flags] <file>",
		Short: "A high-performance CSV query tool powered by DuckDB",
		Long: `csvql is a command-line tool that reads CSV files, infers data types for each column,
			   and allows querying with SQL. It supports various data types including integers,
               floats, booleans, dates, times, and timestamps.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			db, _, err := csvql.Read(args[0], tablename)
			if err != nil {
				return err
			}

			columns, rows, err := csvql.Query(db, query)
			if err != nil {
				return fmt.Errorf("query error: %w", err)
			}

			results := make([]map[string]any, len(rows))
			for rindex, row := range rows {
				entry := make(map[string]any, len(columns))
				for cindex, column := range columns {
					entry[column] = row[cindex]
				}
				results[rindex] = entry
			}

			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(results)
		},
	}

	cmd.Flags().StringVarP(&tablename, "table", "t", "data", "Name of the table for SQL queries")
	cmd.Flags().StringVarP(&query, "query", "q", "SELECT * FROM data LIMIT 20;", "SQL query to execute on the CSV data")

	return cmd
}
