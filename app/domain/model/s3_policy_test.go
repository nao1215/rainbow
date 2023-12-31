package model

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBucketPolicyString(t *testing.T) {
	t.Parallel()
	t.Run("output s3 policy for cloudfront", func(t *testing.T) {
		t.Parallel()

		data, err := os.ReadFile(filepath.Join("testdata", "s3policy.json"))
		if err != nil {
			t.Fatal()
		}

		bp := NewAllowCloudFrontS3BucketPolicy("bucket")
		got, err := bp.String()
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(string(data), got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}
