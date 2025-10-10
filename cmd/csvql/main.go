package main

import "os"

func main() {
	cmd := root()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
