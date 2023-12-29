package subcmd

import (
	"errors"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestQuestion(t *testing.T) {
	type args struct {
		ask string
	}
	tests := []struct {
		name  string
		args  args
		input string
		want  bool
	}{
		{
			name:  "user input 'y'",
			args:  args{"no check"},
			input: "y",
			want:  true,
		},
		{
			name:  "user input 'yes'",
			args:  args{"no check"},
			input: "yes",
			want:  true,
		},
		{
			name:  "user input 'n'",
			args:  args{"no check"},
			input: "n",
			want:  false,
		},
		{
			name:  "user input 'no'",
			args:  args{"no check"},
			input: "no",
			want:  false,
		},
		{
			name:  "user input 'yes' after 'a'",
			args:  args{"no check"},
			input: "a\nyes",
			want:  true,
		},
		{
			name:  "user only input enter",
			args:  args{"no check"},
			input: "\nyes",
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			funcDefer, err := mockStdin(t, tt.input)
			if err != nil {
				t.Fatal(err)
			}
			defer funcDefer()

			if got := Question(os.Stdout, tt.args.ask); got != tt.want {
				t.Errorf("Question() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestion_FmtScanlnErr(t *testing.T) {
	t.Run("fmt.Scanln() return error", func(t *testing.T) {
		orgFmtScanln := FmtScanln
		FmtScanln = func(a ...any) (n int, err error) {
			return -1, errors.New("some error")
		}
		defer func() { FmtScanln = orgFmtScanln }()

		if got := Question(os.Stdout, "no check"); got != false {
			t.Errorf("Question() = %v, want %v", got, false)
		}
	})
}

// mockStdin is a helper function that lets the test pretend dummyInput as os.Stdin.
// It will return a function for `defer` to clean up after the test.
func mockStdin(t *testing.T, dummyInput string) (funcDefer func(), err error) {
	t.Helper()

	oldOsStdin := os.Stdin
	var tmpFile *os.File
	var e error
	if runtime.GOOS != "windows" {
		tmpFile, e = os.CreateTemp(t.TempDir(), strings.ReplaceAll(t.Name(), "/", ""))
	} else {
		// See https://github.com/golang/go/issues/51442
		tmpFile, e = os.CreateTemp(os.TempDir(), strings.ReplaceAll(t.Name(), "/", ""))
	}
	if e != nil {
		return nil, e
	}

	content := []byte(dummyInput)

	if _, err := tmpFile.Write(content); err != nil {
		return nil, err
	}

	if _, err := tmpFile.Seek(0, 0); err != nil {
		return nil, err
	}

	// Set stdin to the temp file
	os.Stdin = tmpFile

	return func() {
		// clean up
		os.Stdin = oldOsStdin
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Log(err)
		}
	}, nil
}
