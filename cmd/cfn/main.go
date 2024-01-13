// Package main is the entrypoint of s3hub.
package main

import (
	"fmt"
	"os"

	"github.com/nao1215/rainbow/cmd/subcmd/cfn"
)

func main() {
	if err := cfn.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
