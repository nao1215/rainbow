package external

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/nao1215/rainbow/app/domain/model"
)

// S3Client returns a new S3 client.
func S3Client(t *testing.T) *s3.Client {
	t.Helper()
	config, err := model.NewAWSConfig(context.Background(), model.NewAWSProfile(""), model.RegionAPNortheast1)
	if err != nil {
		t.Fatal(err)
	}
	client, err := NewS3Client(config)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// CreateS3Buckets creates S3 buckets. Region is fixed to ap-northeast-1.
func CreateS3Buckets(t *testing.T, client *s3.Client, buckets []model.Bucket) {
	t.Helper()

	for _, bucket := range buckets {
		if _, err := client.CreateBucket(context.Background(), &s3.CreateBucketInput{
			Bucket: aws.String(bucket.String()),
			CreateBucketConfiguration: &types.CreateBucketConfiguration{
				LocationConstraint: types.BucketLocationConstraint(model.RegionAPNortheast1.String()),
			},
		}); err != nil {
			t.Fatal(err)
		}
	}
}

// DeleteAllS3BucketDelete deletes all S3 buckets.
func DeleteAllS3BucketDelete(t *testing.T, client *s3.Client) {
	t.Helper()

	buckets, err := client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	if err != nil {
		t.Fatal(err)
	}

	for _, bucket := range buckets.Buckets {
		if _, err := client.DeleteBucket(context.Background(), &s3.DeleteBucketInput{Bucket: bucket.Name}); err != nil {
			t.Fatal(err)
		}
	}
}

// ExistS3Bucket returns true if the bucket exists.
func ExistS3Bucket(t *testing.T, client *s3.Client, bucket model.Bucket) bool {
	t.Helper()

	buckets, err := client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	if err != nil {
		t.Fatal(err)
	}

	for _, b := range buckets.Buckets {
		if *b.Name == bucket.String() {
			return true
		}
	}
	return false
}
