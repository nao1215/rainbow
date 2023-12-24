// Package main is the entrypoint of s3hub.
package main

import (
	"fmt"
	"os"

	"github.com/nao1215/rainbow/cmd/subcmd/s3hub"
)

func main() {
	if err := s3hub.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
