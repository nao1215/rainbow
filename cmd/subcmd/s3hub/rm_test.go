package s3hub

import (
	"bytes"
	"testing"
)

func Test_rm(t *testing.T) {
	t.Skip("TODO: fix this test")
	t.Run("Remove objects from S3 bucket (or remove S3 bucket)", func(t *testing.T) {
		cmd := newRmCmd()
		stdout := bytes.NewBufferString("")
		cmd.SetOutput(stdout)

		if err := cmd.RunE(cmd, []string{}); err != nil {
			t.Errorf("got %v, want nil", err)
		}

		want := "rm is not implemented yet\n"
		got := stdout.String()
		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})
}
