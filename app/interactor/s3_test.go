// Package interactor contains the implementations of usecases.
package interactor

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/domain/service"
	"github.com/nao1215/rainbow/app/external/mock"
	"github.com/nao1215/rainbow/app/usecase"
)

func TestS3BucketCreator_CreateS3Bucket(t *testing.T) {
	t.Parallel()
	t.Run("success to create S3 bucket", func(t *testing.T) {
		t.Parallel()

		s3BucketCreatorMock := mock.S3BucketCreator(func(ctx context.Context, input *service.S3BucketCreatorInput) (*service.S3BucketCreatorOutput, error) {
			if input.Bucket != "bucket-name" {
				t.Errorf("input.Bucket = %s, want %s", input.Bucket, "bucket-name")
			}
			if input.Region != model.RegionAPEast1 {
				t.Errorf("input.Region = %s, want %s", input.Region, model.RegionAPEast1)
			}
			return &service.S3BucketCreatorOutput{}, nil
		})

		s3BucketCreator := NewS3BucketCreator(s3BucketCreatorMock)
		input := &usecase.S3BucketCreatorInput{
			Bucket: "bucket-name",
			Region: model.RegionAPEast1,
		}

		if _, err := s3BucketCreator.CreateS3Bucket(context.Background(), input); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("If bucket name is too short, failed to create bucket", func(t *testing.T) {
		t.Parallel()

		s3BucketCreatorMock := mock.S3BucketCreator(func(ctx context.Context, input *service.S3BucketCreatorInput) (*service.S3BucketCreatorOutput, error) {
			return &service.S3BucketCreatorOutput{}, nil
		})

		s3BucketCreator := NewS3BucketCreator(s3BucketCreatorMock)
		input := &usecase.S3BucketCreatorInput{
			Bucket: "b", // too short
			Region: model.RegionAPEast1,
		}

		if _, err := s3BucketCreator.CreateS3Bucket(context.Background(), input); err == nil {
			t.Fatal("should be failed to create bucket, however err is nil")
		}
	})

	t.Run("If bucket name is too long, failed to create bucket", func(t *testing.T) {
		t.Parallel()

		s3BucketCreatorMock := mock.S3BucketCreator(func(ctx context.Context, input *service.S3BucketCreatorInput) (*service.S3BucketCreatorOutput, error) {
			return &service.S3BucketCreatorOutput{}, nil
		})

		s3BucketCreator := NewS3BucketCreator(s3BucketCreatorMock)
		input := &usecase.S3BucketCreatorInput{
			Bucket: model.Bucket(strings.Repeat("a", model.MaxBucketNameLength+1)), // too long
			Region: model.RegionAPEast1,
		}

		if _, err := s3BucketCreator.CreateS3Bucket(context.Background(), input); err == nil {
			t.Fatal("should be failed to create bucket, however err is nil")
		}
	})

	t.Run("If region is invalid, failed to create bucket", func(t *testing.T) {
		t.Parallel()

		s3BucketCreatorMock := mock.S3BucketCreator(func(ctx context.Context, input *service.S3BucketCreatorInput) (*service.S3BucketCreatorOutput, error) {
			return &service.S3BucketCreatorOutput{}, nil
		})

		s3BucketCreator := NewS3BucketCreator(s3BucketCreatorMock)
		input := &usecase.S3BucketCreatorInput{
			Bucket: model.Bucket("bucket-name"),
			Region: model.Region("invalid-region"),
		}

		if _, err := s3BucketCreator.CreateS3Bucket(context.Background(), input); err == nil {
			t.Fatal("should be failed to create bucket, however err is nil")
		}
	})

	t.Run("An error occurs when calling CreateS3Bucket()", func(t *testing.T) {
		t.Parallel()

		s3BucketCreatorMock := mock.S3BucketCreator(func(ctx context.Context, input *service.S3BucketCreatorInput) (*service.S3BucketCreatorOutput, error) {
			if input.Bucket != "bucket-name" {
				t.Errorf("input.Bucket = %s, want %s", input.Bucket, "bucket-name")
			}
			if input.Region != model.RegionAPEast1 {
				t.Errorf("input.Region = %s, want %s", input.Region, model.RegionAPEast1)
			}
			return nil, errors.New("some error")
		})

		s3BucketCreator := NewS3BucketCreator(s3BucketCreatorMock)
		input := &usecase.S3BucketCreatorInput{
			Bucket: "bucket-name",
			Region: model.RegionAPEast1,
		}

		if _, err := s3BucketCreator.CreateS3Bucket(context.Background(), input); err == nil {
			t.Fatal("should be failed to create bucket, however err is nil")
		}
	})
}

func TestS3BucketLister_ListS3Buckets(t *testing.T) {
	t.Parallel()

	t.Run("success to list S3 buckets", func(t *testing.T) {
		t.Parallel()

		s3BucketListerMock := mock.S3BucketLister(func(ctx context.Context, input *service.S3BucketListerInput) (*service.S3BucketListerOutput, error) {
			return &service.S3BucketListerOutput{
				Buckets: model.BucketSets{
					{
						Bucket:       model.Bucket("bucket-name-A"),
						CreationDate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						Bucket:       model.Bucket("bucket-name-B"),
						CreationDate: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					},
					{
						Bucket:       model.Bucket("bucket-name-C"),
						CreationDate: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
					},
				},
			}, nil
		})

		s3BucketLocationGetter := mock.S3BucketLocationGetter(func(ctx context.Context, input *service.S3BucketLocationGetterInput) (*service.S3BucketLocationGetterOutput, error) {
			switch input.Bucket {
			case "bucket-name-A":
				return &service.S3BucketLocationGetterOutput{
					Region: model.RegionAPEast1,
				}, nil
			case "bucket-name-B":
				return &service.S3BucketLocationGetterOutput{
					Region: model.RegionAPNortheast1,
				}, nil
			case "bucket-name-C":
				return &service.S3BucketLocationGetterOutput{
					Region: model.RegionAPNortheast2,
				}, nil
			default:
				return nil, errors.New("some error")
			}
		})

		s3BucketLister := NewS3BucketLister(s3BucketListerMock, s3BucketLocationGetter)
		got, err := s3BucketLister.ListS3Buckets(context.Background(), nil)
		if err != nil {
			t.Fatal(err)
		}

		want := &usecase.S3BucketListerOutput{
			Buckets: model.BucketSets{
				{
					Bucket:       model.Bucket("bucket-name-A"),
					CreationDate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					Region:       model.RegionAPEast1,
				},
				{
					Bucket:       model.Bucket("bucket-name-B"),
					CreationDate: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					Region:       model.RegionAPNortheast1,
				},
				{
					Bucket:       model.Bucket("bucket-name-C"),
					CreationDate: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
					Region:       model.RegionAPNortheast2,
				},
			},
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("differs: (-want +got)\n%s", diff)
		}
	})

	t.Run("An error occurs when calling ListS3Buckets()", func(t *testing.T) {
		t.Parallel()

		s3BucketListerMock := mock.S3BucketLister(func(ctx context.Context, input *service.S3BucketListerInput) (*service.S3BucketListerOutput, error) {
			return nil, errors.New("some error")
		})

		s3BucketLister := NewS3BucketLister(s3BucketListerMock, nil)
		if _, err := s3BucketLister.ListS3Buckets(context.Background(), nil); err == nil {
			t.Fatal("should be failed to list buckets, however err is nil")
		}
	})

	t.Run("An error occurs when calling GetS3BucketLocation()", func(t *testing.T) {
		t.Parallel()

		s3BucketListerMock := mock.S3BucketLister(func(ctx context.Context, input *service.S3BucketListerInput) (*service.S3BucketListerOutput, error) {
			return &service.S3BucketListerOutput{
				Buckets: model.BucketSets{
					{
						Bucket:       model.Bucket("bucket-name-A"),
						CreationDate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						Bucket:       model.Bucket("bucket-name-B"),
						CreationDate: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					},
					{
						Bucket:       model.Bucket("bucket-name-C"),
						CreationDate: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
					},
				},
			}, nil
		})

		s3BucketLocationGetter := mock.S3BucketLocationGetter(func(ctx context.Context, input *service.S3BucketLocationGetterInput) (*service.S3BucketLocationGetterOutput, error) {
			return nil, errors.New("some error")
		})

		s3BucketLister := NewS3BucketLister(s3BucketListerMock, s3BucketLocationGetter)
		if _, err := s3BucketLister.ListS3Buckets(context.Background(), nil); err == nil {
			t.Fatal("should be failed to list buckets, however err is nil")
		}
	})
}

func TestS3ObjectsLister_ListS3Objects(t *testing.T) {
	t.Parallel()

	t.Run("success to list S3 objects", func(t *testing.T) {
		t.Parallel()

		s3ObjectsListerMock := mock.S3ObjectsLister(func(ctx context.Context, input *service.S3ObjectsListerInput) (*service.S3ObjectsListerOutput, error) {
			if input.Bucket != "bucket-name" {
				t.Errorf("input.Bucket = %s, want %s", input.Bucket, "bucket-name")
			}
			return &service.S3ObjectsListerOutput{
				Objects: model.S3ObjectIdentifiers{
					{
						S3Key:     model.S3Key("object-key-A"),
						VersionID: model.VersionID("version-id-A"),
					},
					{
						S3Key:     model.S3Key("object-key-B"),
						VersionID: model.VersionID("version-id-B"),
					},
					{
						S3Key:     model.S3Key("object-key-C"),
						VersionID: model.VersionID("version-id-C"),
					},
				},
			}, nil
		})

		s3ObjectsLister := NewS3ObjectsLister(s3ObjectsListerMock)
		got, err := s3ObjectsLister.ListS3Objects(context.Background(), &usecase.S3ObjectsListerInput{
			Bucket: model.Bucket("bucket-name"),
		})
		if err != nil {
			t.Fatal(err)
		}

		want := &usecase.S3ObjectsListerOutput{
			Objects: model.S3ObjectIdentifiers{
				{
					S3Key:     model.S3Key("object-key-A"),
					VersionID: model.VersionID("version-id-A"),
				},
				{
					S3Key:     model.S3Key("object-key-B"),
					VersionID: model.VersionID("version-id-B"),
				},
				{
					S3Key:     model.S3Key("object-key-C"),
					VersionID: model.VersionID("version-id-C"),
				},
			},
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("differs: (-want +got)\n%s", diff)
		}
	})

	t.Run("An error occurs when calling ListS3Objects()", func(t *testing.T) {
		t.Parallel()

		s3ObjectsListerMock := mock.S3ObjectsLister(func(ctx context.Context, input *service.S3ObjectsListerInput) (*service.S3ObjectsListerOutput, error) {
			if input.Bucket != "bucket-name" {
				t.Errorf("input.Bucket = %s, want %s", input.Bucket, "bucket-name")
			}
			return nil, errors.New("some error")
		})

		s3ObjectsLister := NewS3ObjectsLister(s3ObjectsListerMock)
		if _, err := s3ObjectsLister.ListS3Objects(context.Background(), &usecase.S3ObjectsListerInput{
			Bucket: model.Bucket("bucket-name"),
		}); err == nil {
			t.Fatal("should be failed to list objects, however err is nil")
		}
	})

	t.Run("If bucket name is too short, failed to list objects", func(t *testing.T) {
		t.Parallel()

		s3ObjectsLister := NewS3ObjectsLister(nil)
		if _, err := s3ObjectsLister.ListS3Objects(context.Background(), &usecase.S3ObjectsListerInput{
			Bucket: "b", // too short
		}); err == nil {
			t.Fatal("should be failed to list objects, however err is nil")
		}
	})
}

func TestS3ObjectsDeleter_DeleteS3Objects(t *testing.T) {
	t.Parallel()

	t.Run("success to delete S3 objects", func(t *testing.T) {
		t.Parallel()

		s3BucketLocationGetter := mock.S3BucketLocationGetter(func(ctx context.Context, input *service.S3BucketLocationGetterInput) (*service.S3BucketLocationGetterOutput, error) {
			switch input.Bucket {
			case "bucket-name":
				return &service.S3BucketLocationGetterOutput{
					Region: model.RegionAPEast1,
				}, nil
			default:
				return nil, errors.New("some error")
			}
		})

		s3ObjectsDeleterMock := mock.S3ObjectsDeleter(func(ctx context.Context, input *service.S3ObjectsDeleterInput) (*service.S3ObjectsDeleterOutput, error) {
			want := &service.S3ObjectsDeleterInput{
				Bucket: model.Bucket("bucket-name"),
				Region: model.RegionAPEast1,
				S3ObjectSets: model.S3ObjectIdentifiers{
					{
						S3Key:     model.S3Key("object-key-A"),
						VersionID: model.VersionID("version-id-A"),
					},
					{
						S3Key:     model.S3Key("object-key-B"),
						VersionID: model.VersionID("version-id-B"),
					},
					{
						S3Key:     model.S3Key("object-key-C"),
						VersionID: model.VersionID("version-id-C"),
					},
				},
			}
			if diff := cmp.Diff(want, input); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
			return &service.S3ObjectsDeleterOutput{}, nil
		})

		s3ObjectVersionLister := mock.S3ObjectVersionsLister(func(ctx context.Context, input *service.S3ObjectVersionsListerInput) (*service.S3ObjectVersionsListerOutput, error) {
			return &service.S3ObjectVersionsListerOutput{
				Objects: model.S3ObjectIdentifiers{},
			}, nil
		})

		s3ObjectsDeleter := NewS3ObjectsDeleter(s3ObjectsDeleterMock, s3BucketLocationGetter, s3ObjectVersionLister)
		if _, err := s3ObjectsDeleter.DeleteS3Objects(context.Background(), &usecase.S3ObjectsDeleterInput{
			Bucket: model.Bucket("bucket-name"),
			S3ObjectIdentifiers: model.S3ObjectIdentifiers{
				{
					S3Key:     model.S3Key("object-key-A"),
					VersionID: model.VersionID("version-id-A"),
				},
				{
					S3Key:     model.S3Key("object-key-B"),
					VersionID: model.VersionID("version-id-B"),
				},
				{
					S3Key:     model.S3Key("object-key-C"),
					VersionID: model.VersionID("version-id-C"),
				},
			},
		}); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("An error occurs when calling DeleteS3Objects()", func(t *testing.T) {
		t.Parallel()

		s3BucketLocationGetter := mock.S3BucketLocationGetter(func(ctx context.Context, input *service.S3BucketLocationGetterInput) (*service.S3BucketLocationGetterOutput, error) {
			switch input.Bucket {
			case "bucket-name":
				return &service.S3BucketLocationGetterOutput{
					Region: model.RegionAPEast1,
				}, nil
			default:
				return nil, errors.New("some error")
			}
		})

		s3ObjectsDeleterMock := mock.S3ObjectsDeleter(func(ctx context.Context, input *service.S3ObjectsDeleterInput) (*service.S3ObjectsDeleterOutput, error) {
			return nil, errors.New("some error")
		})

		s3ObjectVersionLister := mock.S3ObjectVersionsLister(func(ctx context.Context, input *service.S3ObjectVersionsListerInput) (*service.S3ObjectVersionsListerOutput, error) {
			return &service.S3ObjectVersionsListerOutput{
				Objects: model.S3ObjectIdentifiers{},
			}, nil
		})

		s3ObjectsDeleter := NewS3ObjectsDeleter(s3ObjectsDeleterMock, s3BucketLocationGetter, s3ObjectVersionLister)
		if _, err := s3ObjectsDeleter.DeleteS3Objects(context.Background(), &usecase.S3ObjectsDeleterInput{
			Bucket: model.Bucket("bucket-name"),
			S3ObjectIdentifiers: model.S3ObjectIdentifiers{
				{
					S3Key:     model.S3Key("object-key-A"),
					VersionID: model.VersionID("version-id-A"),
				},
				{
					S3Key:     model.S3Key("object-key-B"),
					VersionID: model.VersionID("version-id-B"),
				},
				{
					S3Key:     model.S3Key("object-key-C"),
					VersionID: model.VersionID("version-id-C"),
				},
			},
		}); err == nil {
			t.Fatal("should be failed to delete objects, however err is nil")
		}
	})

	t.Run("An error occurs when calling GetS3BucketLocation()", func(t *testing.T) {
		t.Parallel()

		s3BucketLocationGetter := mock.S3BucketLocationGetter(func(ctx context.Context, input *service.S3BucketLocationGetterInput) (*service.S3BucketLocationGetterOutput, error) {
			return nil, errors.New("some error")
		})

		s3ObjectsDeleter := NewS3ObjectsDeleter(nil, s3BucketLocationGetter, nil)
		if _, err := s3ObjectsDeleter.DeleteS3Objects(context.Background(), &usecase.S3ObjectsDeleterInput{
			Bucket: model.Bucket("bucket-name"),
			S3ObjectIdentifiers: model.S3ObjectIdentifiers{
				{
					S3Key:     model.S3Key("object-key-A"),
					VersionID: model.VersionID("version-id-A"),
				},
				{
					S3Key:     model.S3Key("object-key-B"),
					VersionID: model.VersionID("version-id-B"),
				},
				{
					S3Key:     model.S3Key("object-key-C"),
					VersionID: model.VersionID("version-id-C"),
				},
			},
		}); err == nil {
			t.Fatal("should be failed to delete objects, however err is nil")
		}
	})

	t.Run("If bucket name is too short, failed to delete objects", func(t *testing.T) {
		t.Parallel()

		s3ObjectsDeleter := NewS3ObjectsDeleter(nil, nil, nil)
		if _, err := s3ObjectsDeleter.DeleteS3Objects(context.Background(), &usecase.S3ObjectsDeleterInput{
			Bucket: "b", // too short
		}); err == nil {
			t.Fatal("should be failed to delete objects, however err is nil")
		}
	})
}

func TestS3BucketDeleter_DeleteS3Bucket(t *testing.T) {
	t.Parallel()

	t.Run("success to delete S3 bucket", func(t *testing.T) {
		t.Parallel()

		s3BucketLocationGetter := mock.S3BucketLocationGetter(func(ctx context.Context, input *service.S3BucketLocationGetterInput) (*service.S3BucketLocationGetterOutput, error) {
			switch input.Bucket {
			case "bucket-name":
				return &service.S3BucketLocationGetterOutput{
					Region: model.RegionAPEast1,
				}, nil
			default:
				return nil, errors.New("some error")
			}
		})

		s3BucketDeleterMock := mock.S3BucketDeleter(func(ctx context.Context, input *service.S3BucketDeleterInput) (*service.S3BucketDeleterOutput, error) {
			if input.Bucket != "bucket-name" {
				t.Errorf("input.Bucket = %s, want %s", input.Bucket, "bucket-name")
			}
			if input.Region != model.RegionAPEast1 {
				t.Errorf("input.Region = %s, want %s", input.Region, model.RegionAPEast1)
			}
			return &service.S3BucketDeleterOutput{}, nil
		})

		s3BucketDeleter := NewS3BucketDeleter(s3BucketDeleterMock, s3BucketLocationGetter)
		if _, err := s3BucketDeleter.DeleteS3Bucket(context.Background(), &usecase.S3BucketDeleterInput{
			Bucket: "bucket-name",
		}); err != nil {
			t.Errorf("err = %s, want nil", err)
		}
	})

	t.Run("An error occurs when calling DeleteS3Bucket()", func(t *testing.T) {
		t.Parallel()

		s3BucketLocationGetter := mock.S3BucketLocationGetter(func(ctx context.Context, input *service.S3BucketLocationGetterInput) (*service.S3BucketLocationGetterOutput, error) {
			switch input.Bucket {
			case "bucket-name":
				return &service.S3BucketLocationGetterOutput{
					Region: model.RegionAPEast1,
				}, nil
			default:
				return nil, errors.New("some error")
			}
		})

		s3BucketDeleterMock := mock.S3BucketDeleter(func(ctx context.Context, input *service.S3BucketDeleterInput) (*service.S3BucketDeleterOutput, error) {
			if input.Bucket != "bucket-name" {
				t.Errorf("input.Bucket = %s, want %s", input.Bucket, "bucket-name")
			}
			if input.Region != model.RegionAPEast1 {
				t.Errorf("input.Region = %s, want %s", input.Region, model.RegionAPEast1)
			}
			return nil, errors.New("some error")
		})

		s3BucketDeleter := NewS3BucketDeleter(s3BucketDeleterMock, s3BucketLocationGetter)
		if _, err := s3BucketDeleter.DeleteS3Bucket(context.Background(), &usecase.S3BucketDeleterInput{
			Bucket: "bucket-name",
		}); err == nil {
			t.Fatal("should be failed to delete bucket, however err is nil")
		}
	})

	t.Run("An error occurs when calling GetS3BucketLocation()", func(t *testing.T) {
		t.Parallel()

		s3BucketLocationGetter := mock.S3BucketLocationGetter(func(ctx context.Context, input *service.S3BucketLocationGetterInput) (*service.S3BucketLocationGetterOutput, error) {
			if input.Bucket != "bucket-name" {
				t.Errorf("input.Bucket = %s, want %s", input.Bucket, "bucket-name")
			}
			return nil, errors.New("some error")
		})

		s3BucketDeleter := NewS3BucketDeleter(nil, s3BucketLocationGetter)
		if _, err := s3BucketDeleter.DeleteS3Bucket(context.Background(), &usecase.S3BucketDeleterInput{
			Bucket: "bucket-name",
		}); err == nil {
			t.Fatal("should be failed to delete bucket, however err is nil")
		}
	})

	t.Run("An error occurs when calling GetS3BucketLocation()", func(t *testing.T) {
		t.Parallel()

		s3BucketLocationGetter := mock.S3BucketLocationGetter(func(ctx context.Context, input *service.S3BucketLocationGetterInput) (*service.S3BucketLocationGetterOutput, error) {
			return nil, errors.New("some error")
		})

		s3BucketDeleter := NewS3BucketDeleter(nil, s3BucketLocationGetter)
		if _, err := s3BucketDeleter.DeleteS3Bucket(context.Background(), &usecase.S3BucketDeleterInput{
			Bucket: "bucket-name",
		}); err == nil {
			t.Fatal("should be failed to delete bucket, however err is nil")
		}
	})

	t.Run("If bucket name is too short, failed to delete bucket", func(t *testing.T) {
		t.Parallel()

		s3BucketDeleter := NewS3BucketDeleter(nil, nil)
		if _, err := s3BucketDeleter.DeleteS3Bucket(context.Background(), &usecase.S3BucketDeleterInput{
			Bucket: "b", // too short
		}); err == nil {
			t.Fatal("should be failed to delete bucket, however err is nil")
		}
	})
}

func TestFileUploader_UploadFile(t *testing.T) {
	t.Parallel()

	t.Run("success to upload file", func(t *testing.T) {
		t.Parallel()

		s3ObjectUploader := mock.S3ObjectUploader(func(ctx context.Context, input *service.S3ObjectUploaderInput) (*service.S3ObjectUploaderOutput, error) {
			want := &service.S3ObjectUploaderInput{
				Bucket:   model.Bucket("bucket-name"),
				Region:   model.RegionAFSouth1,
				S3Key:    model.S3Key("object-key"),
				S3Object: model.NewS3Object([]byte("some data")),
			}

			opt := cmpopts.IgnoreUnexported(model.S3Object{}, bytes.Buffer{})
			if diff := cmp.Diff(want, input, opt); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}

			return &service.S3ObjectUploaderOutput{
				ContentType:   "text/plain",
				ContentLength: 100,
			}, nil
		})

		fileUploader := NewFileUploader(s3ObjectUploader)
		got, err := fileUploader.UploadFile(context.Background(), &usecase.FileUploaderInput{
			Bucket: "bucket-name",
			Region: model.RegionAFSouth1,
			Key:    "object-key",
			Data:   []byte("some data"),
		})
		if err != nil {
			t.Fatal(err)
		}

		want := &usecase.FileUploaderOutput{
			ContentType:   "text/plain",
			ContentLength: 100,
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("differs: (-want +got)\n%s", diff)
		}
	})

	t.Run("An error occurs when calling UploadFile()", func(t *testing.T) {
		t.Parallel()

		s3ObjectUploader := mock.S3ObjectUploader(func(ctx context.Context, input *service.S3ObjectUploaderInput) (*service.S3ObjectUploaderOutput, error) {
			return nil, errors.New("some error")
		})

		fileUploader := NewFileUploader(s3ObjectUploader)
		if _, err := fileUploader.UploadFile(context.Background(), &usecase.FileUploaderInput{
			Bucket: "bucket-name",
			Region: model.RegionAFSouth1,
			Key:    "object-key",
			Data:   []byte("some data"),
		}); err == nil {
			t.Fatal("should be failed to upload file, however err is nil")
		}
	})

	t.Run("If bucket name is too short, failed to upload file", func(t *testing.T) {
		t.Parallel()

		fileUploader := NewFileUploader(nil)
		if _, err := fileUploader.UploadFile(context.Background(), &usecase.FileUploaderInput{
			Bucket: "b", // too short
			Region: model.RegionAFSouth1,
			Key:    "object-key",
			Data:   []byte("some data"),
		}); err == nil {
			t.Fatal("should be failed to upload file, however err is nil")
		}
	})

	t.Run("If region is invalid, failed to upload file", func(t *testing.T) {
		t.Parallel()

		fileUploader := NewFileUploader(nil)
		if _, err := fileUploader.UploadFile(context.Background(), &usecase.FileUploaderInput{
			Bucket: "bucket-name",
			Region: model.Region("invalid-region"),
			Key:    "object-key",
			Data:   []byte("some data"),
		}); err == nil {
			t.Fatal("should be failed to upload file, however err is nil")
		}
	})
}

func TestS3BucketPublicAccessBlocker_BlockS3BucketPublicAccess(t *testing.T) {
	t.Parallel()

	t.Run("success to block S3 bucket public access", func(t *testing.T) {
		t.Parallel()

		s3BucketPublicAccessBlockerMock := mock.S3BucketPublicAccessBlocker(func(ctx context.Context, input *service.S3BucketPublicAccessBlockerInput) (*service.S3BucketPublicAccessBlockerOutput, error) {
			want := &service.S3BucketPublicAccessBlockerInput{
				Bucket: model.Bucket("bucket-name"),
				Region: model.RegionAFSouth1,
			}

			if diff := cmp.Diff(want, input); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}

			return &service.S3BucketPublicAccessBlockerOutput{}, nil
		})

		s3BucketPublicAccessBlocker := NewS3BucketPublicAccessBlocker(s3BucketPublicAccessBlockerMock)
		if _, err := s3BucketPublicAccessBlocker.BlockS3BucketPublicAccess(context.Background(), &usecase.S3BucketPublicAccessBlockerInput{
			Bucket: "bucket-name",
			Region: model.RegionAFSouth1,
		}); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("An error occurs when calling BlockS3BucketPublicAccess()", func(t *testing.T) {
		t.Parallel()

		s3BucketPublicAccessBlockerMock := mock.S3BucketPublicAccessBlocker(func(ctx context.Context, input *service.S3BucketPublicAccessBlockerInput) (*service.S3BucketPublicAccessBlockerOutput, error) {
			if input.Bucket != "bucket-name" {
				t.Errorf("input.Bucket = %s, want %s", input.Bucket, "bucket-name")
			}
			if input.Region != model.RegionAFSouth1 {
				t.Errorf("input.Region = %s, want %s", input.Region, model.RegionAFSouth1)
			}
			return nil, errors.New("some error")
		})

		s3BucketPublicAccessBlocker := NewS3BucketPublicAccessBlocker(s3BucketPublicAccessBlockerMock)
		if _, err := s3BucketPublicAccessBlocker.BlockS3BucketPublicAccess(context.Background(), &usecase.S3BucketPublicAccessBlockerInput{
			Bucket: "bucket-name",
			Region: model.RegionAFSouth1,
		}); err == nil {
			t.Fatal("should be failed to block public access, however err is nil")
		}
	})

	t.Run("If bucket name is too short, failed to block public access", func(t *testing.T) {
		t.Parallel()

		s3BucketPublicAccessBlocker := NewS3BucketPublicAccessBlocker(nil)
		if _, err := s3BucketPublicAccessBlocker.BlockS3BucketPublicAccess(context.Background(), &usecase.S3BucketPublicAccessBlockerInput{
			Bucket: "b", // too short
			Region: model.RegionAFSouth1,
		}); err == nil {
			t.Fatal("should be failed to block public access, however err is nil")
		}
	})

	t.Run("If region is invalid, failed to block public access", func(t *testing.T) {
		t.Parallel()

		s3BucketPublicAccessBlocker := NewS3BucketPublicAccessBlocker(nil)
		if _, err := s3BucketPublicAccessBlocker.BlockS3BucketPublicAccess(context.Background(), &usecase.S3BucketPublicAccessBlockerInput{
			Bucket: "bucket-name",
			Region: model.Region("invalid-region"),
		}); err == nil {
			t.Fatal("should be failed to block public access, however err is nil")
		}
	})
}

func TestS3BucketPolicySetter_SetS3BucketPolicy(t *testing.T) {
	t.Parallel()

	t.Run("success to set S3 bucket policy", func(t *testing.T) {
		t.Parallel()

		s3BucketPolicySetterMock := mock.S3BucketPolicySetter(func(ctx context.Context, input *service.S3BucketPolicySetterInput) (*service.S3BucketPolicySetterOutput, error) {
			want := &service.S3BucketPolicySetterInput{
				Bucket: model.Bucket("bucket-name"),
				Policy: model.NewAllowCloudFrontS3BucketPolicy(model.Bucket("bucket-name")),
			}
			if diff := cmp.Diff(want, input); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
			return &service.S3BucketPolicySetterOutput{}, nil
		})

		s3BucketPolicySetter := NewS3BucketPolicySetter(s3BucketPolicySetterMock)
		if _, err := s3BucketPolicySetter.SetS3BucketPolicy(context.Background(), &usecase.S3BucketPolicySetterInput{
			Bucket: "bucket-name",
			Policy: model.NewAllowCloudFrontS3BucketPolicy(model.Bucket("bucket-name")),
		}); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("An error occurs when calling SetS3BucketPolicy()", func(t *testing.T) {
		t.Parallel()

		s3BucketPolicySetterMock := mock.S3BucketPolicySetter(func(ctx context.Context, input *service.S3BucketPolicySetterInput) (*service.S3BucketPolicySetterOutput, error) {
			want := &service.S3BucketPolicySetterInput{
				Bucket: model.Bucket("bucket-name"),
				Policy: model.NewAllowCloudFrontS3BucketPolicy(model.Bucket("bucket-name")),
			}
			if diff := cmp.Diff(want, input); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
			return nil, errors.New("some error")
		})

		s3BucketPolicySetter := NewS3BucketPolicySetter(s3BucketPolicySetterMock)
		if _, err := s3BucketPolicySetter.SetS3BucketPolicy(context.Background(), &usecase.S3BucketPolicySetterInput{
			Bucket: "bucket-name",
			Policy: model.NewAllowCloudFrontS3BucketPolicy(model.Bucket("bucket-name")),
		}); err == nil {
			t.Fatal("should be failed to set bucket policy, however err is nil")
		}
	})

	t.Run("If bucket name is too short, failed to set bucket policy", func(t *testing.T) {
		t.Parallel()

		s3BucketPolicySetter := NewS3BucketPolicySetter(nil)
		if _, err := s3BucketPolicySetter.SetS3BucketPolicy(context.Background(), &usecase.S3BucketPolicySetterInput{
			Bucket: "b", // too short
			Policy: model.NewAllowCloudFrontS3BucketPolicy(model.Bucket("bucket-name")),
		}); err == nil {
			t.Fatal("should be failed to set bucket policy, however err is nil")
		}
	})
}

func TestS3ObjectDownloader_DownloadS3Object(t *testing.T) {
	t.Parallel()

	t.Run("success to download S3 object", func(t *testing.T) {
		t.Parallel()

		s3ObjectDownloaderMock := mock.S3ObjectDownloader(func(ctx context.Context, input *service.S3ObjectDownloaderInput) (*service.S3ObjectDownloaderOutput, error) {
			want := &service.S3ObjectDownloaderInput{
				Bucket: model.Bucket("bucket-name"),
				Key:    model.S3Key("object-key"),
			}
			if diff := cmp.Diff(want, input); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
			return &service.S3ObjectDownloaderOutput{
				Bucket:        model.Bucket("bucket-name"),
				Key:           model.S3Key("object-key"),
				ContentType:   "text/plain",
				ContentLength: 100,
				S3Object:      model.NewS3Object([]byte("some data")),
			}, nil
		})

		s3ObjectDownloader := NewS3ObjectDownloader(s3ObjectDownloaderMock)
		got, err := s3ObjectDownloader.DownloadS3Object(context.Background(), &usecase.S3ObjectDownloaderInput{
			Bucket: "bucket-name",
			Key:    "object-key",
		})
		if err != nil {
			t.Fatal(err)
		}

		want := &usecase.S3ObjectDownloaderOutput{
			Bucket:        model.Bucket("bucket-name"),
			Key:           model.S3Key("object-key"),
			ContentType:   "text/plain",
			ContentLength: 100,
			S3Object:      model.NewS3Object([]byte("some data")),
		}

		opt := cmpopts.IgnoreUnexported(model.S3Object{}, bytes.Buffer{})
		if diff := cmp.Diff(want, got, opt); diff != "" {
			t.Errorf("differs: (-want +got)\n%s", diff)
		}
	})

	t.Run("An error occurs when calling DownloadS3Object()", func(t *testing.T) {
		t.Parallel()

		s3ObjectDownloaderMock := mock.S3ObjectDownloader(func(ctx context.Context, input *service.S3ObjectDownloaderInput) (*service.S3ObjectDownloaderOutput, error) {
			want := &service.S3ObjectDownloaderInput{
				Bucket: model.Bucket("bucket-name"),
				Key:    model.S3Key("object-key"),
			}
			if diff := cmp.Diff(want, input); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
			return nil, errors.New("some error")
		})

		s3ObjectDownloader := NewS3ObjectDownloader(s3ObjectDownloaderMock)
		if _, err := s3ObjectDownloader.DownloadS3Object(context.Background(), &usecase.S3ObjectDownloaderInput{
			Bucket: "bucket-name",
			Key:    "object-key",
		}); err == nil {
			t.Fatal("should be failed to download object, however err is nil")
		}
	})

	t.Run("If bucket name is too short, failed to download object", func(t *testing.T) {
		t.Parallel()

		s3ObjectDownloader := NewS3ObjectDownloader(nil)
		if _, err := s3ObjectDownloader.DownloadS3Object(context.Background(), &usecase.S3ObjectDownloaderInput{
			Bucket: "b", // too short
			Key:    "object-key",
		}); err == nil {
			t.Fatal("should be failed to download object, however err is nil")
		}
	})
}

func TestS3ObjectCopier_CopyS3Object(t *testing.T) {
	t.Parallel()

	t.Run("success to copy S3 object", func(t *testing.T) {
		t.Parallel()

		s3ObjectCopierMock := mock.S3ObjectCopier(func(ctx context.Context, input *service.S3ObjectCopierInput) (*service.S3ObjectCopierOutput, error) {
			want := &service.S3ObjectCopierInput{
				SourceBucket:      model.Bucket("source-bucket-name"),
				SourceKey:         model.S3Key("source-object-key"),
				DestinationBucket: model.Bucket("dest-bucket-name"),
				DestinationKey:    model.S3Key("dest-object-key"),
			}
			if diff := cmp.Diff(want, input); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
			return &service.S3ObjectCopierOutput{}, nil
		})

		s3ObjectCopier := NewS3ObjectCopier(s3ObjectCopierMock)
		if _, err := s3ObjectCopier.CopyS3Object(context.Background(), &usecase.S3ObjectCopierInput{
			SourceBucket:      "source-bucket-name",
			SourceKey:         "source-object-key",
			DestinationBucket: "dest-bucket-name",
			DestinationKey:    "dest-object-key",
		}); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("An error occurs when calling CopyS3Object()", func(t *testing.T) {
		t.Parallel()

		s3ObjectCopierMock := mock.S3ObjectCopier(func(ctx context.Context, input *service.S3ObjectCopierInput) (*service.S3ObjectCopierOutput, error) {
			want := &service.S3ObjectCopierInput{
				SourceBucket:      model.Bucket("source-bucket-name"),
				SourceKey:         model.S3Key("source-object-key"),
				DestinationBucket: model.Bucket("dest-bucket-name"),
				DestinationKey:    model.S3Key("dest-object-key"),
			}
			if diff := cmp.Diff(want, input); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
			return nil, errors.New("some error")
		})

		s3ObjectCopier := NewS3ObjectCopier(s3ObjectCopierMock)
		if _, err := s3ObjectCopier.CopyS3Object(context.Background(), &usecase.S3ObjectCopierInput{
			SourceBucket:      "source-bucket-name",
			SourceKey:         "source-object-key",
			DestinationBucket: "dest-bucket-name",
			DestinationKey:    "dest-object-key",
		}); err == nil {
			t.Fatal("should be failed to copy object, however err is nil")
		}
	})

	t.Run("If source bucket name is too short, failed to copy object", func(t *testing.T) {
		t.Parallel()

		s3ObjectCopier := NewS3ObjectCopier(nil)
		if _, err := s3ObjectCopier.CopyS3Object(context.Background(), &usecase.S3ObjectCopierInput{
			SourceBucket:      "b", // too short
			SourceKey:         "source-object-key",
			DestinationBucket: "dest-bucket-name",
			DestinationKey:    "dest-object-key",
		}); err == nil {
			t.Fatal("should be failed to copy object, however err is nil")
		}
	})

	t.Run("If destination bucket name is too short, failed to copy object", func(t *testing.T) {
		t.Parallel()

		s3ObjectCopier := NewS3ObjectCopier(nil)
		if _, err := s3ObjectCopier.CopyS3Object(context.Background(), &usecase.S3ObjectCopierInput{
			SourceBucket:      "source-bucket-name",
			SourceKey:         "source-object-key",
			DestinationBucket: "b", // too short
			DestinationKey:    "dest-object-key",
		}); err == nil {
			t.Fatal("should be failed to copy object, however err is nil")
		}
	})
}
