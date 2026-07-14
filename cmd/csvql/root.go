package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yaien/csvql"
)

func root() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "csvql [flags] <file> <tablename> <db>",
		Short: `csvql is a command-line tool that exports a csv into a duckdb database`,
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			filename, tablename, dbpath := args[0], args[1], args[2]

			db, err := csvql.CreateDB(dbpath)
			if err != nil {
				return fmt.Errorf("failed creating db: %w", err)
			}

			schema, records, err := csvql.ParseFile(filename, tablename)
			if err != nil {
				return fmt.Errorf("failed parsing file: %w", err)
			}

			err = csvql.ImportOnDB(db, schema, records)
			if err != nil {
				return fmt.Errorf("failed importing file: %w", err)
			}

			return nil
		},
	}

	return cmd
}
