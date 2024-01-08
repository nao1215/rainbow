// Package external implements the external service.
package external

import (
	"context"
	"testing"

	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/domain/service"
)

func TestS3BucketCreator_CreateS3Bucket(t *testing.T) {
	deleteAllS3BucketDelete(t, s3client(t))

	t.Run("success to create bucket", func(t *testing.T) {
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

		if !existS3Bucket(t, client, input.Bucket) {
			t.Errorf("%s bucket does not exist", input.Bucket)
		}
	})
}
