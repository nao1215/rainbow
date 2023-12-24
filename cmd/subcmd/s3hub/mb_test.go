package s3hub

import (
	"bytes"
	"testing"
)

func Test_mb(t *testing.T) {
	t.Run("Make S3 bucket", func(t *testing.T) {
		cmd := newMbCmd()
		stdout := bytes.NewBufferString("")
		cmd.SetOutput(stdout)

		if err := cmd.RunE(cmd, []string{}); err != nil {
			t.Errorf("got %v, want nil", err)
		}

		want := "mb is not implemented yet\n"
		got := stdout.String()
		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})
}
