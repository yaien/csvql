package main

import (
	"net"
	"net/http"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"github.com/yaien/csvql"
	"github.com/yaien/csvql/csvqlserver"
	"github.com/yaien/csvql/csvqlsite"
)

func serve() *cobra.Command {
	var port string
	var maxMemory int64
	var dbpath string
	var tablename string

	cmd := &cobra.Command{
		Use:  "serve [file.csv]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			db, schema, err := csvql.Read(args[0], tablename)
			if err != nil {
				return err
			}

			schemas := []*csvql.Schema{schema}

			server := csvqlserver.New(db, schemas)
			server.SetMaxMemory(maxMemory)
			server.Route(http.DefaultServeMux)

			csvqlsite.Route(http.DefaultServeMux)

			addr := net.JoinHostPort("localhost", port)

			go func() {
				err := browser.OpenURL("http://" + addr)
				if err != nil {
					cmd.PrintErrf("failed to open browser: %v\n", err)
				}
			}()

			return http.ListenAndServe(addr, nil)
		},
	}

	cmd.Flags().StringVarP(&port, "port", "p", "8047", "port to listen on")
	cmd.Flags().Int64Var(&maxMemory, "max-memory", 1024*1024*1024, "Maximum memory for multipart form parsing (in bytes)")
	cmd.Flags().StringVar(&dbpath, "dbpath", "", "path to DuckDB database file (empty for in-memory)")
	cmd.Flags().StringVar(&tablename, "table", "data", "table name to use for the CSV data")

	return cmd
}
