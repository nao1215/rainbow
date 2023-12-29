// Package model contains the definitions of domain models and business logic.
package model

import (
	"errors"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
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
			e:       ErrEmptyRegion,
		},
		{
			name:    "failure. region is invalid",
			r:       Region("invalid"),
			wantErr: true,
			e:       ErrInvalidRegion,
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
