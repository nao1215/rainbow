package s3hub

import (
	"bytes"
	"testing"
)

func Test_cp(t *testing.T) {
	t.Skip("TODO: fix this test")
	t.Run("Copy file from local(S3 bucket) to S3 bucket(local)", func(t *testing.T) {
		cmd := newCpCmd()
		stdout := bytes.NewBufferString("")
		cmd.SetOutput(stdout)

		if err := cmd.RunE(cmd, []string{}); err != nil {
			t.Errorf("got %v, want nil", err)
		}

		want := "cp is not implemented yet\n"
		got := stdout.String()
		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})
}
