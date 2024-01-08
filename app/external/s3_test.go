// Package external implements the external service.
package external

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/nao1215/rainbow/app/domain"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/domain/service"
)

func TestS3BucketCreator_CreateS3Bucket(t *testing.T) {
	t.Run("success to create bucket", func(t *testing.T) {
		DeleteAllS3BucketDelete(t, S3Client(t))

		config, err := model.NewAWSConfig(context.Background(), model.NewAWSProfile(""), model.RegionAPNortheast1)
		if err != nil {
			t.Fatal(err)
		}
		client, err := NewS3Client(config)
		if err != nil {
			t.Fatal(err)
		}

		s3BucketCreator := NewS3BucketCreator(client)
		input := &service.S3BucketCreatorInput{
			Bucket: model.Bucket("test-bucket"),
			Region: model.RegionAPNortheast1,
		}
		if _, err = s3BucketCreator.CreateS3Bucket(context.Background(), input); err != nil {
			t.Error(err)
		}

		if !ExistS3Bucket(t, client, input.Bucket) {
			t.Errorf("%s bucket does not exist", input.Bucket)
		}
	})

	t.Run("fail to create bucket because the bucket already exists", func(t *testing.T) {
		DeleteAllS3BucketDelete(t, S3Client(t))

		config, err := model.NewAWSConfig(context.Background(), model.NewAWSProfile(""), model.RegionAPNortheast1)
		if err != nil {
			t.Fatal(err)
		}
		client, err := NewS3Client(config)
		if err != nil {
			t.Fatal(err)
		}

		s3BucketCreator := NewS3BucketCreator(client)
		input := &service.S3BucketCreatorInput{
			Bucket: model.Bucket("test-bucket"),
			Region: model.RegionAPNortheast1,
		}
		if _, err = s3BucketCreator.CreateS3Bucket(context.Background(), input); err != nil {
			t.Error(err)
		}

		_, got := s3BucketCreator.CreateS3Bucket(context.Background(), input)
		if errors.Is(got, domain.ErrBucketAlreadyExistsOwnedByOther) {
			t.Errorf("got %v, want %v", got, domain.ErrBucketAlreadyExistsOwnedByOther)
		}
	})

	t.Run("fail to create bucket because the bucket name is invalid", func(t *testing.T) {
		DeleteAllS3BucketDelete(t, S3Client(t))

		config, err := model.NewAWSConfig(context.Background(), model.NewAWSProfile(""), model.RegionAPNortheast1)
		if err != nil {
			t.Fatal(err)
		}
		client, err := NewS3Client(config)
		if err != nil {
			t.Fatal(err)
		}

		s3BucketCreator := NewS3BucketCreator(client)
		input := &service.S3BucketCreatorInput{
			Bucket: model.Bucket("a"), // invalid bucket name
			Region: model.RegionAPNortheast1,
		}
		_, got := s3BucketCreator.CreateS3Bucket(context.Background(), input)
		if got == nil {
			t.Errorf("got %v, want %v", got, domain.ErrInvalidBucketName)
		}
	})

	t.Run("fail to create bucket because the region is invalid", func(t *testing.T) {
		DeleteAllS3BucketDelete(t, S3Client(t))

		config, err := model.NewAWSConfig(context.Background(), model.NewAWSProfile(""), model.RegionAPNortheast1)
		if err != nil {
			t.Fatal(err)
		}
		client, err := NewS3Client(config)
		if err != nil {
			t.Fatal(err)
		}

		s3BucketCreator := NewS3BucketCreator(client)
		input := &service.S3BucketCreatorInput{
			Bucket: model.Bucket("test-bucket"),
			Region: model.Region(""), // invalid region
		}
		_, got := s3BucketCreator.CreateS3Bucket(context.Background(), input)
		if got == nil {
			t.Errorf("got %v, want %v", got, domain.ErrInvalidRegion)
		}
	})

	t.Run("Create S3 Bucket at 'us-east-1'", func(t *testing.T) {
		DeleteAllS3BucketDelete(t, S3Client(t))

		config, err := model.NewAWSConfig(context.Background(), model.NewAWSProfile(""), model.RegionUSEast1)
		if err != nil {
			t.Fatal(err)
		}
		client, err := NewS3Client(config)
		if err != nil {
			t.Fatal(err)
		}

		s3BucketCreator := NewS3BucketCreator(client)
		input := &service.S3BucketCreatorInput{
			Bucket: model.Bucket("test-bucket"),
			Region: model.RegionUSEast1,
		}
		if _, err = s3BucketCreator.CreateS3Bucket(context.Background(), input); err != nil {
			t.Error(err)
		}
		if !ExistS3Bucket(t, client, input.Bucket) {
			t.Errorf("%s bucket does not exist", input.Bucket)
		}
	})
}

func TestS3BucketLister_ListS3Buckets(t *testing.T) {
	t.Run("success to list buckets", func(t *testing.T) {
		DeleteAllS3BucketDelete(t, S3Client(t))

		config, err := model.NewAWSConfig(context.Background(), model.NewAWSProfile(""), model.RegionAPNortheast1)
		if err != nil {
			t.Fatal(err)
		}
		client, err := NewS3Client(config)
		if err != nil {
			t.Fatal(err)
		}

		CreateS3Buckets(t, client, []model.Bucket{
			model.Bucket("test-bucket-3"),
			model.Bucket("test-bucket-1"),
			model.Bucket("test-bucket-2"),
		})

		s3BucketLister := NewS3BucketLister(client)
		got, err := s3BucketLister.ListS3Buckets(context.Background(), &service.S3BucketListerInput{})
		if err != nil {
			t.Error(err)
		}

		want := &service.S3BucketListerOutput{
			Buckets: model.BucketSets{
				{
					Bucket: model.Bucket("test-bucket-3"),
				},
				{
					Bucket: model.Bucket("test-bucket-1"),
				},
				{
					Bucket: model.Bucket("test-bucket-2"),
				},
			},
		}

		opt := cmpopts.IgnoreFields(model.BucketSet{}, "CreationDate")
		if diff := cmp.Diff(want, got, opt); diff != "" {
			t.Errorf("differs: (-want +got)\n%s", diff)
		}
	})

	t.Run("success to list buckets when there is no bucket", func(t *testing.T) {
		DeleteAllS3BucketDelete(t, S3Client(t))

		config, err := model.NewAWSConfig(context.Background(), model.NewAWSProfile(""), model.RegionAPNortheast1)
		if err != nil {
			t.Fatal(err)
		}
		client, err := NewS3Client(config)
		if err != nil {
			t.Fatal(err)
		}

		s3BucketLister := NewS3BucketLister(client)
		got, err := s3BucketLister.ListS3Buckets(context.Background(), &service.S3BucketListerInput{})
		if err != nil {
			t.Error(err)
		}

		want := &service.S3BucketListerOutput{
			Buckets: model.BucketSets{},
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("differs: (-want +got)\n%s", diff)
		}
	})
}

func TestS3BucketLocationGetter_GetS3BucketLocation(t *testing.T) {
	t.Run("success to get bucket location", func(t *testing.T) {
		DeleteAllS3BucketDelete(t, S3Client(t))

		config, err := model.NewAWSConfig(context.Background(), model.NewAWSProfile(""), model.RegionAPNortheast1)
		if err != nil {
			t.Fatal(err)
		}
		client, err := NewS3Client(config)
		if err != nil {
			t.Fatal(err)
		}

		CreateS3Buckets(t, client, []model.Bucket{
			model.Bucket("test-bucket-1"),
		})

		s3BucketLocationGetter := NewS3BucketLocationGetter(client)
		got, err := s3BucketLocationGetter.GetS3BucketLocation(context.Background(), &service.S3BucketLocationGetterInput{
			Bucket: model.Bucket("test-bucket-1"),
		})
		if err != nil {
			t.Error(err)
		}

		want := &service.S3BucketLocationGetterOutput{
			Region: model.RegionAPNortheast1,
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("differs: (-want +got)\n%s", diff)
		}
	})

	t.Run("fail to get bucket location because the bucket does not exist", func(t *testing.T) {
		DeleteAllS3BucketDelete(t, S3Client(t))

		config, err := model.NewAWSConfig(context.Background(), model.NewAWSProfile(""), model.RegionAPNortheast1)
		if err != nil {
			t.Fatal(err)
		}
		client, err := NewS3Client(config)
		if err != nil {
			t.Fatal(err)
		}

		s3BucketLocationGetter := NewS3BucketLocationGetter(client)
		_, got := s3BucketLocationGetter.GetS3BucketLocation(context.Background(), &service.S3BucketLocationGetterInput{
			Bucket: model.Bucket("test-bucket-1"),
		})
		if got == nil {
			t.Error("want error, however got nil")
		}
	})

	t.Run("Get S3 bucket location that is at 'us-east-1'", func(t *testing.T) {
		DeleteAllS3BucketDelete(t, S3Client(t))

		config, err := model.NewAWSConfig(context.Background(), model.NewAWSProfile(""), model.RegionUSEast1)
		if err != nil {
			t.Fatal(err)
		}
		client, err := NewS3Client(config)
		if err != nil {
			t.Fatal(err)
		}

		s3BucketCreator := NewS3BucketCreator(client)
		input := &service.S3BucketCreatorInput{
			Bucket: model.Bucket("test-bucket"),
			Region: model.RegionUSEast1,
		}
		if _, err = s3BucketCreator.CreateS3Bucket(context.Background(), input); err != nil {
			t.Fatal(err)
		}

		s3BucketLocationGetter := NewS3BucketLocationGetter(client)
		got, err := s3BucketLocationGetter.GetS3BucketLocation(context.Background(), &service.S3BucketLocationGetterInput{
			Bucket: model.Bucket("test-bucket"),
		})
		if err != nil {
			t.Error(err)
		}

		want := &service.S3BucketLocationGetterOutput{
			Region: model.RegionUSEast1,
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("differs: (-want +got)\n%s", diff)
		}
	})

	t.Run("fail to get bucket location because specified bucket does not exist", func(t *testing.T) {
		t.Skip("This test does not pass using locakstack")

		DeleteAllS3BucketDelete(t, S3Client(t))

		s3BucketLocationGetter := NewS3BucketLocationGetter(S3Client(t))
		_, got := s3BucketLocationGetter.GetS3BucketLocation(context.Background(), &service.S3BucketLocationGetterInput{
			Bucket: model.Bucket("test-bucket"),
		})
		if !errors.Is(got, domain.ErrNoSuchBucket) {
			t.Errorf("got %v, want %v", got, domain.ErrNoSuchBucket)
		}
	})
}

func TestS3BucketDeleter_DeleteS3Bucket(t *testing.T) {
	t.Run("success to delete bucket", func(t *testing.T) {
		DeleteAllS3BucketDelete(t, S3Client(t))

		config, err := model.NewAWSConfig(context.Background(), model.NewAWSProfile(""), model.RegionUSEast1)
		if err != nil {
			t.Fatal(err)
		}
		client, err := NewS3Client(config)
		if err != nil {
			t.Fatal(err)
		}
		s3BucketCreator := NewS3BucketCreator(client)
		input := &service.S3BucketCreatorInput{
			Bucket: model.Bucket("test-bucket"),
			Region: model.RegionUSEast1,
		}
		if _, err = s3BucketCreator.CreateS3Bucket(context.Background(), input); err != nil {
			t.Fatal(err)
		}

		s3BucketDeleter := NewS3BucketDeleter(client)
		if _, err = s3BucketDeleter.DeleteS3Bucket(context.Background(), &service.S3BucketDeleterInput{
			Bucket: model.Bucket("test-bucket"),
			Region: model.RegionUSEast1,
		}); err != nil {
			t.Error(err)
		}

		if ExistS3Bucket(t, client, input.Bucket) {
			t.Errorf("%s bucket exists", input.Bucket)
		}
	})

	t.Run("fail to delete bucket because the bucket does not exist", func(t *testing.T) {
		t.Skip("This test does not pass using locakstack")

		DeleteAllS3BucketDelete(t, S3Client(t))

		s3BucketDeleter := NewS3BucketDeleter(S3Client(t))
		_, got := s3BucketDeleter.DeleteS3Bucket(context.Background(), &service.S3BucketDeleterInput{
			Bucket: model.Bucket("test-bucket"),
			Region: model.RegionUSEast1,
		})
		if !errors.Is(got, domain.ErrNoSuchBucket) {
			t.Errorf("got %v, want %v", got, domain.ErrNoSuchBucket)
		}
	})
}
