// Package model contains the definitions of domain models and business logic.
package model

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/nao1215/rainbow/app/domain"
)

func TestRegionString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		r    Region
		want string
	}{
		{
			name: "success",
			r:    RegionAPNortheast1,
			want: "ap-northeast-1",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.r.String(); got != tt.want {
				t.Errorf("Region.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegionValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		r       Region
		wantErr bool
		e       error
	}{
		{
			name:    "success",
			r:       RegionAPNortheast1,
			wantErr: false,
			e:       nil,
		},
		{
			name:    "failure. region is empty",
			r:       Region(""),
			wantErr: true,
			e:       domain.ErrEmptyRegion,
		},
		{
			name:    "failure. region is invalid",
			r:       Region("invalid"),
			wantErr: true,
			e:       domain.ErrInvalidRegion,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.r.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Region.Validate() error = %v, wantErr %v", err, tt.wantErr)
				if tt.wantErr {
					if errors.Is(err, tt.e) {
						t.Errorf("error mismatch got = %v, wantErr %v", err, tt.wantErr)
					}
				}
			}
		})
	}
}

func TestBucketString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		b    Bucket
		want string
	}{
		{
			name: "success",
			b:    Bucket("spare"),
			want: "spare",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.b.String(); got != tt.want {
				t.Errorf("Bucket.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketValidateLength(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		b       Bucket
		wantErr bool
	}{
		{
			name:    "success: minimum length",
			b:       Bucket("abc"),
			wantErr: false,
		},
		{
			name:    "success: maximum length",
			b:       Bucket(strings.Repeat("a", 63)),
			wantErr: false,
		},
		{
			name:    "failure. bucket name is too short",
			b:       Bucket("ab"),
			wantErr: true,
		},
		{
			name:    "failure. bucket name is too long",
			b:       Bucket(strings.Repeat("a", 64)),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.b.validateLength(); (err != nil) != tt.wantErr {
				t.Errorf("Bucket.validateLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBucketValidatePattern(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		b       Bucket
		wantErr bool
	}{
		{
			name:    "success",
			b:       Bucket("abc"),
			wantErr: false,
		},
		{
			name:    "failure. bucket name contains invalid character",
			b:       Bucket("abc!"),
			wantErr: true,
		},
		{
			name:    "failure. bucket name contains uppercase character",
			b:       Bucket("Abc"),
			wantErr: true,
		},
		{
			name:    "failure. bucket name contains underscore",
			b:       Bucket("abc_def"),
			wantErr: true,
		},
		{
			name:    "failure. bucket name starts with hyphen",
			b:       Bucket("-abc"),
			wantErr: true,
		},
		{
			name:    "failure. bucket name ends with hyphen",
			b:       Bucket("abc-"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.b.validatePattern(); (err != nil) != tt.wantErr {
				t.Errorf("Bucket.validatePattern() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBucketValidatePrefix(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		b       Bucket
		wantErr bool
	}{
		{
			name:    "success",
			b:       Bucket("abc"),
			wantErr: false,
		},
		{
			name:    "failure. bucket name starts with 'xn--'",
			b:       Bucket("xn--abc"),
			wantErr: true,
		},
		{
			name:    "failure. bucket name starts with 'sthree-'",
			b:       Bucket("sthree-abc"),
			wantErr: true,
		},
		{
			name:    "failure. bucket name starts with 'sthree-configurator'",
			b:       Bucket("sthree-configurator-abc"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.b.validatePrefix(); (err != nil) != tt.wantErr {
				t.Errorf("Bucket.validatePrefix() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBucketValidateSuffix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		b       Bucket
		wantErr bool
	}{
		{
			name:    "success",
			b:       Bucket("abc"),
			wantErr: false,
		},
		{
			name:    "failure. bucket name ends with '-s3alias'",
			b:       Bucket("abc-s3alias"),
			wantErr: true,
		},
		{
			name:    "failure. bucket name ends with '--ol-s3'",
			b:       Bucket("abc--ol-s3"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.b.validateSuffix(); (err != nil) != tt.wantErr {
				t.Errorf("Bucket.validateSuffix() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBucketValidateCharSequence(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		b       Bucket
		wantErr bool
	}{
		{
			name:    "success",
			b:       Bucket("abc"),
			wantErr: false,
		},
		{
			name:    "failure. bucket name contains consecutive periods",
			b:       Bucket("abc..def"),
			wantErr: true,
		},
		{
			name:    "failure. bucket name contains consecutive hyphens",
			b:       Bucket("abc--def"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.b.validateCharSequence(); (err != nil) != tt.wantErr {
				t.Errorf("Bucket.validateCharSequence() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBucketValidate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		b       Bucket
		wantErr bool
	}{
		{
			name:    "success",
			b:       Bucket("abc"),
			wantErr: false,
		},
		{
			name:    "failure. bucket name is empty",
			b:       Bucket(""),
			wantErr: true,
		},
		{
			name:    "failure. bucket name is too short",
			b:       Bucket("ab"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.b.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Bucket.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBucketDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		b    Bucket
		want string
	}{
		{
			name: "success",
			b:    Bucket("abc"),
			want: "abc.s3.amazonaws.com",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.b.Domain(); got != tt.want {
				t.Errorf("Bucket.Domain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegion_Next(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		r    Region
		want Region
	}{
		{
			name: "success",
			r:    RegionAPNortheast1,
			want: RegionAPNortheast2,
		},
		{
			name: "success. last region",
			r:    RegionUSGovWest1,
			want: RegionUSEast1,
		},
		{
			name: "failure. invalid region",
			r:    Region("invalid"),
			want: RegionAPNortheast1,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.r.Next(); got != tt.want {
				t.Errorf("Region.Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegion_Prev(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		r    Region
		want Region
	}{
		{
			name: "success",
			r:    RegionAPNortheast2,
			want: RegionAPNortheast1,
		},
		{
			name: "success. first region",
			r:    RegionUSEast1,
			want: RegionUSGovWest1,
		},
		{
			name: "failure. invalid region",
			r:    Region("invalid"),
			want: RegionAPNortheast1,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.r.Prev(); got != tt.want {
				t.Errorf("Region.Prev() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketSets_Len(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		b    BucketSets
		want int
	}{
		{
			name: "If BucketSets has two BucketSet, Len() returns 2",
			b: BucketSets{
				BucketSet{
					Bucket: Bucket("abc"),
				},
				BucketSet{
					Bucket: Bucket("def"),
				},
			},
			want: 2,
		},
		{
			name: "If BucketSets is empty, Len() returns 0",
			b:    BucketSets{},
			want: 0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.b.Len(); got != tt.want {
				t.Errorf("BucketSets.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketSets_Empty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		b    BucketSets
		want bool
	}{
		{
			name: "If BucketSets has two BucketSet, Empty() returns false",
			b: BucketSets{
				BucketSet{
					Bucket: Bucket("abc"),
				},
				BucketSet{
					Bucket: Bucket("def"),
				},
			},
			want: false,
		},
		{
			name: "If BucketSets is empty, Empty() returns true",
			b:    BucketSets{},
			want: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.b.Empty(); got != tt.want {
				t.Errorf("BucketSets.Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketSets_Contains(t *testing.T) {
	t.Parallel()

	type args struct {
		bucket Bucket
	}
	tests := []struct {
		name string
		b    BucketSets
		args args
		want bool
	}{
		{
			name: "If BucketSets has two BucketSet including the bucket 'abc' and input is 'abc', Contains() returns true",
			b: BucketSets{
				BucketSet{
					Bucket: Bucket("abc"),
				},
				BucketSet{
					Bucket: Bucket("def"),
				},
			},
			args: args{
				bucket: Bucket("abc"),
			},
			want: true,
		},
		{
			name: "If BucketSets has two BucketSet including the bucket 'abc' and input is 'def', Contains() returns false",
			b: BucketSets{
				BucketSet{
					Bucket: Bucket("abc"),
				},
				BucketSet{
					Bucket: Bucket("def"),
				},
			},
			args: args{
				bucket: Bucket("def"),
			},
			want: true,
		},
		{
			name: "If BucketSets is empty, Contains() returns false",
			b:    BucketSets{},
			args: args{
				bucket: Bucket("abc"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.b.Contains(tt.args.bucket); got != tt.want {
				t.Errorf("BucketSets.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucket_TrimKey(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skip this test on Windows")
	}
	t.Parallel()

	tests := []struct {
		name string
		b    Bucket
		want Bucket
	}{
		{
			name: "If Bucket is 'abc', TrimKey() returns 'abc'",
			b:    Bucket("abc"),
			want: Bucket("abc"),
		},
		{
			name: "If Bucket is 'abc/', TrimKey() returns 'abc'",
			b:    Bucket("abc/"),
			want: Bucket("abc"),
		},
		{
			name: "If Bucket is 'abc/def', TrimKey() returns 'abc/def'",
			b:    Bucket(filepath.Join("abc", "def")),
			want: Bucket("abc"),
		},
		{
			name: "If Bucket is 'abc/def/', TrimKey() returns 'abc'",
			b:    Bucket(filepath.Join("abc", "def/")),
			want: Bucket("abc"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.b.TrimKey(); got != tt.want {
				t.Errorf("Bucket.TrimKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucket_Split(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skip this test on Windows")
	}

	t.Parallel()

	tests := []struct {
		name  string
		b     Bucket
		want  Bucket
		want1 S3Key
	}{
		{
			name:  "If Bucket is 'abc', Split() returns 'abc' and ''",
			b:     Bucket("abc"),
			want:  Bucket("abc"),
			want1: S3Key(""),
		},
		{
			name:  "If Bucket is 'abc/', Split() returns 'abc' and ''",
			b:     Bucket("abc/"),
			want:  Bucket("abc"),
			want1: S3Key(""),
		},
		{
			name:  "If Bucket is 'abc/def', Split() returns 'abc' and 'def'",
			b:     Bucket(filepath.Join("abc", "def")),
			want:  Bucket("abc"),
			want1: S3Key("def"),
		},
		{
			name:  "If Bucket is 'abc/def/', Split() returns 'abc' and 'def/'",
			b:     Bucket(filepath.Join("abc", "def/")),
			want:  Bucket("abc"),
			want1: S3Key("def"),
		},
		{
			name:  "If Bucket is 'abc/def/ghi', Split() returns 'abc' and 'def/ghi'",
			b:     Bucket(filepath.Join("abc", "def", "ghi")),
			want:  Bucket("abc"),
			want1: S3Key(filepath.Join("def", "ghi")),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, got1 := tt.b.Split()
			if got != tt.want {
				t.Errorf("Bucket.Split() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Bucket.Split() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestS3Key_Empty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		k    S3Key
		want bool
	}{
		{
			name: "If S3Key is 'abc', Empty() returns false",
			k:    S3Key("abc"),
			want: false,
		},
		{
			name: "If S3Key is '', Empty() returns true",
			k:    S3Key(""),
			want: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.k.Empty(); got != tt.want {
				t.Errorf("S3Key.Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestS3Key_IsAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		k    S3Key
		want bool
	}{
		{
			name: "If S3Key is 'abc', IsAll() returns false",
			k:    S3Key("abc"),
			want: false,
		},
		{
			name: "If S3Key is '', IsAll() returns false",
			k:    S3Key(""),
			want: false,
		},
		{
			name: "If S3Key is '*', IsAll() returns true",
			k:    S3Key("*"),
			want: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.k.IsAll(); got != tt.want {
				t.Errorf("S3Key.IsAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDeleteRetryCount(t *testing.T) {
	t.Parallel()

	type args struct {
		i int
	}
	tests := []struct {
		name string
		args args
		want DeleteObjectsRetryCount
	}{
		{
			name: "input is 1, NewDeleteRetryCount() returns 1",
			args: args{
				i: 1,
			},
			want: DeleteObjectsRetryCount(1),
		},
		{
			name: "input is 0, NewDeleteRetryCount() returns 0",
			args: args{
				i: 0,
			},
			want: DeleteObjectsRetryCount(0),
		},
		{
			name: "input is -1, NewDeleteRetryCount() returns 0",
			args: args{
				i: -1,
			},
			want: DeleteObjectsRetryCount(0),
		},
		{
			name: "input is over MaxS3DeleteObjectsRetryCount, NewDeleteRetryCount() returns MaxS3DeleteObjectsRetryCount",
			args: args{
				i: MaxS3DeleteObjectsRetryCount + 1,
			},
			want: DeleteObjectsRetryCount(MaxS3DeleteObjectsRetryCount),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NewDeleteRetryCount(tt.args.i); got != tt.want {
				t.Errorf("NewDeleteRetryCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucket_WithProtocol(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		b    Bucket
		want Bucket
	}{
		{
			name: "If Bucket is 'abc', WithProtocol() returns 's3://abc'",
			b:    Bucket("abc"),
			want: Bucket("s3://abc"),
		},
		{
			name: "If Bucket is 's3://abc', WithProtocol() returns 's3://abc'",
			b:    Bucket("s3://abc"),
			want: Bucket("s3://abc"),
		},
		{
			name: "If Bucket is 's3://abc/def', WithProtocol() returns 's3://abc/def'",
			b:    Bucket("s3://abc/def"),
			want: Bucket("s3://abc/def"),
		},
		{
			name: "If Bucket is '', WithProtocol() returns 's3://'",
			b:    Bucket(""),
			want: Bucket("s3://"),
		},
		{
			name: "If Bucket is 's3://', WithProtocol() returns 's3://'",
			b:    Bucket("s3://"),
			want: Bucket("s3://"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.b.WithProtocol(); got != tt.want {
				t.Errorf("Bucket.WithProtocol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBucketWithoutProtocol(t *testing.T) {
	t.Parallel()
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want Bucket
	}{
		{
			name: "If input is 's3://abc', NewBucketWithoutProtocol() returns 'abc'",
			args: args{
				s: "s3://abc",
			},
			want: Bucket("abc"),
		},
		{
			name: "If input is 's3://abc/def', NewBucketWithoutProtocol() returns 'abc/def'",
			args: args{
				s: "s3://abc/def",
			},
			want: Bucket("abc/def"),
		},
		{
			name: "If input is 'abc', NewBucketWithoutProtocol() returns 'abc'",
			args: args{
				s: "abc",
			},
			want: Bucket("abc"),
		},
		{
			name: "If input is 'abc/def', NewBucketWithoutProtocol() returns 'abc/def'",
			args: args{
				s: "abc/def",
			},
			want: Bucket("abc/def"),
		},
		{
			name: "If input is '', NewBucketWithoutProtocol() returns ''",
			args: args{
				s: "",
			},
			want: Bucket(""),
		},
		{
			name: "If input is 's3://', NewBucketWithoutProtocol() returns ''",
			args: args{
				s: "s3://",
			},
			want: Bucket(""),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NewBucketWithoutProtocol(tt.args.s); got != tt.want {
				t.Errorf("NewBucketWithoutProtocol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucket_Join(t *testing.T) {
	t.Parallel()

	type args struct {
		key S3Key
	}
	tests := []struct {
		name string
		b    Bucket
		args args
		want Bucket
	}{
		{
			name: "If Bucket is 'abc' and key is 'def', Join() returns 'abc/def'",
			b:    Bucket("abc"),
			args: args{
				S3Key("def"),
			},
			want: Bucket("abc/def"),
		},
		{
			name: "If Bucket is 'abc' and key is 'def/ghi', Join() returns 'abc/def/ghi'",
			b:    Bucket("abc"),
			args: args{
				S3Key("def/ghi"),
			},
			want: Bucket("abc/def/ghi"),
		},
		{
			name: "If Bucket is 'abc' and key is '', Join() returns 'abc'",
			b:    Bucket("abc"),
			args: args{
				S3Key(""),
			},
			want: Bucket("abc"),
		},
		{
			name: "If Bucket is 'abc' and key is 'def/', Join() returns 'abc/def'",
			b:    Bucket("abc"),
			args: args{
				S3Key("def/"),
			},
			want: Bucket("abc/def"),
		},
		{
			name: "If Bucket is '' and key is 'def', Join() returns ''",
			b:    Bucket(""),
			args: args{
				S3Key("def"),
			},
			want: Bucket(""),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.b.Join(tt.args.key); got != tt.want {
				t.Errorf("Bucket.Join() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestS3ObjectIdentifiers_Len(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		s    S3ObjectIdentifiers
		want int
	}{
		{
			name: "If S3ObjectIdentifiers has two S3ObjectIdentifierSet, Len() returns 2",
			s:    S3ObjectIdentifiers{S3ObjectIdentifier{}, S3ObjectIdentifier{}},
			want: 2,
		},
		{
			name: "If S3ObjectIdentifiers is empty, Len() returns 0",
			s:    S3ObjectIdentifiers{},
			want: 0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.s.Len(); got != tt.want {
				t.Errorf("S3ObjectIdentifiers.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortS3ObjectIdentifiers(t *testing.T) {
	t.Parallel()
	t.Run("If S3ObjectIdentifiers has three S3ObjectIdentifierSet, sort.Sort returns sorted S3ObjectIdentifiers", func(t *testing.T) {
		t.Parallel()
		s := S3ObjectIdentifiers{
			S3ObjectIdentifier{
				S3Key: S3Key("ghi"),
			},
			S3ObjectIdentifier{
				S3Key: S3Key("abc"),
			},
			S3ObjectIdentifier{
				S3Key: S3Key("def"),
			},
		}
		want := S3ObjectIdentifiers{
			S3ObjectIdentifier{
				S3Key: S3Key("abc"),
			},
			S3ObjectIdentifier{
				S3Key: S3Key("def"),
			},
			S3ObjectIdentifier{
				S3Key: S3Key("ghi"),
			},
		}
		sort.Sort(s)
		if diff := cmp.Diff(s, want); diff != "" {
			t.Errorf("sort.Sort() mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestS3ObjectIdentifiers_ToS3ObjectIdentifiers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		s    S3ObjectIdentifiers
		want []types.ObjectIdentifier
	}{
		{
			name: "If S3ObjectIdentifiers has two S3ObjectIdentifierSet, ToS3ObjectIdentifiers() returns []types.ObjectIdentifier",
			s: S3ObjectIdentifiers{
				{
					S3Key:     S3Key("abc"),
					VersionID: VersionID("def"),
				},
				{
					S3Key:     S3Key("ghi"),
					VersionID: VersionID("jkl"),
				},
				{
					S3Key:     S3Key("mno"),
					VersionID: VersionID("pqr"),
				},
			},
			want: []types.ObjectIdentifier{
				{
					Key:       aws.String("abc"),
					VersionId: aws.String("def"),
				},
				{
					Key:       aws.String("ghi"),
					VersionId: aws.String("jkl"),
				},
				{
					Key:       aws.String("mno"),
					VersionId: aws.String("pqr"),
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.s.ToS3ObjectIdentifiers()

			opt := cmpopts.IgnoreUnexported(types.ObjectIdentifier{})
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("S3ObjectIdentifiers.ToS3ObjectIdentifiers() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestS3Key_Join(t *testing.T) {
	t.Parallel()

	type args struct {
		key S3Key
	}
	tests := []struct {
		name string
		k    S3Key
		args args
		want S3Key
	}{
		{
			name: "If S3Key is 'abc' and key is 'def', Join() returns 'abc/def'",
			k:    S3Key("abc"),
			args: args{
				key: S3Key("def"),
			},
			want: S3Key("abc/def"),
		},
		{
			name: "If S3Key is 'abc' and key is 'def/ghi', Join() returns 'abc/def/ghi'",
			k:    S3Key("abc"),
			args: args{
				key: S3Key("def/ghi"),
			},
			want: S3Key("abc/def/ghi"),
		},
		{
			name: "If S3Key is 'abc' and key is '', Join() returns 'abc'",
			k:    S3Key("abc"),
			args: args{
				key: S3Key(""),
			},
			want: S3Key("abc"),
		},
		{
			name: "If S3Key is 'abc' and key is 'def/', Join() returns 'abc/def'",
			k:    S3Key("abc"),
			args: args{
				key: S3Key("def/"),
			},
			want: S3Key("abc/def"),
		},
		{
			name: "If S3Key is '' and key is '/def', Join() returns 'def'",
			k:    S3Key(""),
			args: args{
				key: S3Key("/def"),
			},
			want: S3Key("def"),
		},
		{
			name: "If S3Key is '' and key is 'def', Join() returns 'def'",
			args: args{
				key: S3Key("def"),
			},
			want: S3Key("def"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.k.Join(tt.args.key); got != tt.want {
				t.Errorf("S3Key.Join() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestS3Object_ToFile(t *testing.T) {
	t.Parallel()

	t.Run("If S3Object is 'abc', ToFile() writes 'abc' to the file", func(t *testing.T) {
		t.Parallel()

		want := []byte("abc")
		obj := NewS3Object(want)
		tmpDir := os.TempDir()
		tmpFilePath := filepath.Join(tmpDir, "s3object.txt")

		if err := obj.ToFile(tmpFilePath, 0600); err != nil {
			t.Fatalf("S3Object.ToFile() error = %v", err)
		}

		got, err := os.ReadFile(filepath.Clean(tmpFilePath))
		if err != nil {
			t.Fatalf("os.ReadFile() error = %v", err)
		}

		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("S3Object.ToFile() mismatch (-want +got):\n%s", diff)
		}

		if err := os.RemoveAll(tmpFilePath); err != nil {
			t.Fatalf("os.RemoveAll() error = %v", err)
		}
	})
}

func TestS3Object_ContentType(t *testing.T) {
	t.Parallel()

	t.Run("If S3Object is png file, ContentType() returns 'image/png'", func(t *testing.T) {
		t.Parallel()

		b, err := os.ReadFile(filepath.Join("testdata", "lena.png"))
		if err != nil {
			t.Fatalf("os.ReadFile() error = %v", err)
		}

		got := NewS3Object(b).ContentType()
		want := "image/png"
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("S3Object.ContentType() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("If S3Object is json file, ContentType() returns 'application/json'", func(t *testing.T) {
		t.Parallel()

		b, err := os.ReadFile(filepath.Join("testdata", "s3policy.json"))
		if err != nil {
			t.Fatalf("os.ReadFile() error = %v", err)
		}

		got := NewS3Object(b).ContentType()
		want := "application/json"
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("S3Object.ContentType() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("If S3Object is markdown file, ContentType() returns 'text/plain; charset=utf-8'", func(t *testing.T) {
		t.Parallel()

		b, err := os.ReadFile(filepath.Join("testdata", "sample.md"))
		if err != nil {
			t.Fatalf("os.ReadFile() error = %v", err)
		}

		got := NewS3Object(b).ContentType()
		want := "text/plain; charset=utf-8"

		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("S3Object.ContentType() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("If S3Object is empty, ContentType() returns 'text/plain'", func(t *testing.T) {
		t.Parallel()

		got := NewS3Object([]byte{}).ContentType()
		want := "text/plain"

		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("S3Object.ContentType() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("If S3Object is nil, ContentType() returns 'text/plain'", func(t *testing.T) {
		t.Parallel()

		got := NewS3Object(nil).ContentType()
		want := "text/plain"

		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("S3Object.ContentType() mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestS3Object_ContentLength(t *testing.T) {
	t.Parallel()

	type fields struct {
		Buffer *bytes.Buffer
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			name: "If S3Object is 'abc', ContentLength() returns 3",
			fields: fields{
				Buffer: bytes.NewBuffer([]byte("abc")),
			},
			want: 3,
		},
		{
			name: "If S3Object is empty, ContentLength() returns 0",
			fields: fields{
				Buffer: bytes.NewBuffer([]byte{}),
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &S3Object{
				Buffer: tt.fields.Buffer,
			}
			if got := s.ContentLength(); got != tt.want {
				t.Errorf("S3Object.ContentLength() = %v, want %v", got, tt.want)
			}
		})
	}
}
