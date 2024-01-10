package cfn

import (
	"bytes"
	"testing"

	ver "github.com/nao1215/rainbow/version"
)

func Test_version(t *testing.T) {
	t.Run("Get version information", func(t *testing.T) {
		cmd := newVersionCmd()
		stdout := bytes.NewBufferString("")
		cmd.SetOutput(stdout)

		orgVersion := ver.Version
		ver.Version = "v0.0.0"
		t.Cleanup(func() {
			ver.Version = orgVersion
		})

		cmd.Run(cmd, []string{})

		want := "s3hub v0.0.0 (under MIT LICENSE)\n"
		got := stdout.String()
		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})
}
