package s3hub

import (
	"bytes"
	"testing"
)

func Test_ls(t *testing.T) {
	t.Skip("TODO: fix this test")
	t.Run("List S3 buckets", func(t *testing.T) {
		cmd := newLsCmd()
		stdout := bytes.NewBufferString("")
		cmd.SetOutput(stdout)

		if err := cmd.RunE(cmd, []string{}); err != nil {
			t.Errorf("got %v, want nil", err)
		}

		want := "ls is not implemented yet\n"
		got := stdout.String()
		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})
}
