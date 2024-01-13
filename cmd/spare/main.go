// Package main is the entrypoint of spare.
package main

import (
	"fmt"
	"os"

	"github.com/nao1215/rainbow/cmd/subcmd/spare"
)

func main() {
	if err := spare.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
