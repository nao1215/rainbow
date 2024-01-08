package external

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/nao1215/rainbow/app/domain/model"
)

func s3client(t *testing.T) *s3.Client {
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

func deleteAllS3BucketDelete(t *testing.T, client *s3.Client) {
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

func existS3Bucket(t *testing.T, client *s3.Client, bucket model.Bucket) bool {
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
