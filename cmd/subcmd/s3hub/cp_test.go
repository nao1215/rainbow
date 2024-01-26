package s3hub

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
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

func Test_newCopyPathPair(t *testing.T) {
	t.Parallel()

	type args struct {
		from string
		to   string
	}
	tests := []struct {
		name string
		args args
		want *copyPathPair
	}{
		{
			name: "copyTypeLocalToS3",
			args: args{
				from: "/path/to/file.txt",
				to:   "s3://mybucket/path/to/file.txt",
			},
			want: &copyPathPair{
				From: "/path/to/file.txt",
				To:   "s3://mybucket/path/to/file.txt",
				Type: copyTypeLocalToS3,
			},
		},
		{
			name: "copyTypeS3ToLocal",
			args: args{
				from: "s3://mybucket/path/to/file.txt",
				to:   "/path/to/file.txt",
			},
			want: &copyPathPair{
				From: "s3://mybucket/path/to/file.txt",
				To:   "/path/to/file.txt",
				Type: copyTypeS3ToLocal,
			},
		},
		{
			name: "copyTypeS3ToS3",
			args: args{
				from: "s3://mybucket1/path/to/file.txt",
				to:   "s3://mybucket2/path/to/file.txt",
			},
			want: &copyPathPair{
				From: "s3://mybucket1/path/to/file.txt",
				To:   "s3://mybucket2/path/to/file.txt",
				Type: copyTypeS3ToS3,
			},
		},
		{
			name: "copyTypeUnknown: from local to local",
			args: args{
				from: "/path/to/file.txt",
				to:   "/path/to/file.txt",
			},
			want: &copyPathPair{
				From: "/path/to/file.txt",
				To:   "/path/to/file.txt",
				Type: copyTypeUnknown,
			},
		},
		{
			name: "copyTypeUnknown: from is empty",
			args: args{
				from: "",
				to:   "/path/to/file.txt",
			},
			want: &copyPathPair{
				From: "",
				To:   "/path/to/file.txt",
				Type: copyTypeUnknown,
			},
		},
		{
			name: "copyTypeUnknown: to is empty",
			args: args{
				from: "/path/to/file.txt",
				to:   "",
			},
			want: &copyPathPair{
				From: "/path/to/file.txt",
				To:   "",
				Type: copyTypeUnknown,
			},
		},
		{
			name: "copyTypeUnknown: from and to are empty",
			args: args{
				from: "",
				to:   "",
			},
			want: &copyPathPair{
				From: "",
				To:   "",
				Type: copyTypeUnknown,
			},
		},
		{
			name: "copyTypeUnknown: use file:// protocol",
			args: args{
				from: "file:///path/to/file.txt",
				to:   "file:///path/to/file.txt",
			},
			want: &copyPathPair{
				From: "file:///path/to/file.txt",
				To:   "file:///path/to/file.txt",
				Type: copyTypeUnknown,
			},
		},
		{
			name: "copyTypeUnknown: use bad s3:// protocol",
			args: args{
				from: "s3:/mybucket/path/to/file.txt",
				to:   "s3:/mybucket/path/to/file.txt",
			},
			want: &copyPathPair{
				From: "s3:/mybucket/path/to/file.txt",
				To:   "s3:/mybucket/path/to/file.txt",
				Type: copyTypeUnknown,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := newCopyPathPair(tt.args.from, tt.args.to)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("newCopyPathPair() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
