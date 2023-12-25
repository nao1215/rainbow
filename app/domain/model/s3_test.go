// Package model contains the definitions of domain models and business logic.
package model

import (
	"errors"
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
