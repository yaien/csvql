package main

import (
	"fmt"
	"os"
	"text/tabwriter"

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

			writter := tabwriter.NewWriter(os.Stdout, 0, 4, 1, ' ', 0)

			rows, cols, err := csvql.Query(db, query)
			if err != nil {
				return fmt.Errorf("query error: %w", err)
			}

			// Print header
			for i, col := range cols {
				if i > 0 {
					_, _ = fmt.Fprint(writter, "\t")
				}
				_, _ = fmt.Fprint(writter, col)
			}

			_, _ = fmt.Fprintln(writter)

			// Print rows
			for _, row := range rows {
				for i, value := range row {
					if i > 0 {
						_, _ = fmt.Fprint(writter, "\t")
					}
					_, _ = fmt.Fprint(writter, value)
				}
				_, _ = fmt.Fprintln(writter)
			}

			_ = writter.Flush()

			return nil
		},
	}

	cmd.Flags().StringVarP(&tablename, "table", "t", "data", "Name of the table for SQL queries")
	cmd.Flags().StringVarP(&query, "query", "q", "SELECT * FROM data LIMIT 20;", "SQL query to execute on the CSV data")

	return cmd
}
