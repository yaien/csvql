package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	csvql "github.com/yaien/csvviewer"
)

func root() *cobra.Command {
	var tablename string
	var query string

	cmd := &cobra.Command{
		Use:   "csvviewer [flags] <file>",
		Short: "A simple CSV viewer with type inference",
		Long: `csvviewer is a command-line tool that reads a CSV file, infers data types for each column,
			   and displays the content in a formatted table. It supports various data types including integers,
               floats, booleans, dates, times, and datetimes.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			db, _, err := csvql.Read(args[0], tablename)
			if err != nil {
				return err
			}

			writter := tabwriter.NewWriter(os.Stdout, 0, 4, 1, ' ', 0)

			rows, err := db.Query(query)
			if err != nil {
				return fmt.Errorf("query error: %w", err)
			}

			cols, err := rows.Columns()
			if err != nil {
				return err
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
			for rows.Next() {

				columns := make([]any, len(cols))
				pointers := make([]any, len(cols))
				for i := range columns {
					pointers[i] = &columns[i]
				}

				if err := rows.Scan(pointers...); err != nil {
					return err
				}

				for i, col := range columns {
					if i > 0 {
						_, _ = fmt.Fprint(writter, "\t")
					}
					_, _ = fmt.Fprint(writter, col)
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
