package s3hub

import (
	"testing"
)

func Test_interactive(t *testing.T) {
	t.Skip("TODO: fix this test")
	t.Run("Interactive mode", func(t *testing.T) {
		if err := interactive(); err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})
}
