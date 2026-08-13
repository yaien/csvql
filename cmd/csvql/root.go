package main

import (
	"database/sql"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yaien/csvql"
	_ "modernc.org/sqlite"
)

func root() *cobra.Command {
	var comma string
	cmd := &cobra.Command{
		Use:   "csvql [flags] <file> <tablename> <db>",
		Short: `csvql is a command-line tool that exports a csv into a sqlite database`,
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			filename, tablename, dbpath := args[0], args[1], args[2]

			db, err := sql.Open("sqlite", dbpath)
			if err != nil {
				return fmt.Errorf("failed opening db: %w", err)
			}

			err = csvql.ImportOnDB(db, filename, tablename, csvql.WithComma(comma))
			if err != nil {
				return fmt.Errorf("failed importing file: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&comma, "--comma", "-c", ",", "set the field separator to a comma")

	return cmd
}
