package s3hub

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/rainbow/app/di"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/interactor/mock"
	"github.com/nao1215/rainbow/app/usecase"
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

func Test_cpCmd_filterS3Objects(t *testing.T) {
	t.Parallel()

	t.Run("filterS3Objects ", func(t *testing.T) {
		mockLister := mock.S3ObjectsLister(func(ctx context.Context, input *usecase.S3ObjectsListerInput) (*usecase.S3ObjectsListerOutput, error) {
			want := &usecase.S3ObjectsListerInput{
				Bucket: model.NewBucketWithoutProtocol("mybucket"),
			}
			if diff := cmp.Diff(input, want); diff != "" {
				t.Errorf("got %v, want %v", input, want)
			}
			return &usecase.S3ObjectsListerOutput{
				Objects: model.S3ObjectIdentifiers{
					{S3Key: model.S3Key("path/to/file1.txt")},
					{S3Key: model.S3Key("path/to/file2.txt")},
					{S3Key: model.S3Key("path/to/file3.txt")},
				},
			}, nil
		})

		cpCmd := &cpCmd{
			s3hub: &s3hub{
				S3App: &di.S3App{
					S3ObjectsLister: mockLister,
				},
				ctx: context.Background(),
			},
		}

		got, err := cpCmd.filterS3Objects(model.NewBucketWithoutProtocol("mybucket"), model.S3Key("path/to"))
		if err != nil {
			t.Errorf("got %v, want nil", err)
		}

		want := []model.S3Key{
			model.S3Key("path/to/file1.txt"),
			model.S3Key("path/to/file2.txt"),
			model.S3Key("path/to/file3.txt"),
		}

		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("filterS3Objects: no objects found", func(t *testing.T) {
		mockLister := mock.S3ObjectsLister(func(ctx context.Context, input *usecase.S3ObjectsListerInput) (*usecase.S3ObjectsListerOutput, error) {
			want := &usecase.S3ObjectsListerInput{
				Bucket: model.NewBucketWithoutProtocol("mybucket"),
			}
			if diff := cmp.Diff(input, want); diff != "" {
				t.Errorf("got %v, want %v", input, want)
			}
			return &usecase.S3ObjectsListerOutput{
				Objects: model.S3ObjectIdentifiers{},
			}, nil
		})

		cpCmd := &cpCmd{
			s3hub: &s3hub{
				S3App: &di.S3App{
					S3ObjectsLister: mockLister,
				},
				ctx: context.Background(),
			},
		}

		_, err := cpCmd.filterS3Objects(model.NewBucketWithoutProtocol("mybucket"), model.S3Key("path/to"))
		if err == nil {
			t.Errorf("got nil, want error")
		}

		want := "no objects found. bucket=mybucket, key=path/to"
		got := err.Error()
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("filterS3Objects: ListS3Objects returns error", func(t *testing.T) {
		mockLister := mock.S3ObjectsLister(func(ctx context.Context, input *usecase.S3ObjectsListerInput) (*usecase.S3ObjectsListerOutput, error) {
			want := &usecase.S3ObjectsListerInput{
				Bucket: model.NewBucketWithoutProtocol("mybucket"),
			}
			if diff := cmp.Diff(input, want); diff != "" {
				t.Errorf("got %v, want %v", input, want)
			}
			return nil, errors.New("dummy error")
		})

		cpCmd := &cpCmd{
			s3hub: &s3hub{
				S3App: &di.S3App{
					S3ObjectsLister: mockLister,
				},
				ctx: context.Background(),
			},
		}

		_, err := cpCmd.filterS3Objects(model.NewBucketWithoutProtocol("mybucket"), model.S3Key("path/to"))
		if err == nil {
			t.Errorf("got nil, want error")
		}

		want := "dummy error: bucket=mybucket"
		got := err.Error()
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
