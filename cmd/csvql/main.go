package main

import "os"

func main() {
	cmd := root()
	cmd.AddCommand(serve())
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
