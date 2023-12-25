// Package model contains the definitions of domain models and business logic.
package model

import (
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
			r:    Region("ap-northeast-1"),
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

func TestBucketValid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		b    Bucket
		want bool
	}{
		{
			name: "success",
			b:    Bucket("rainbow"),
			want: true,
		},
		{
			name: "failure. bucket name is empty",
			b:    Bucket(""),
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.b.Valid(); got != tt.want {
				t.Errorf("Bucket.Valid() = %v, want %v", got, tt.want)
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
			b:    Bucket("rainbow"),
			want: "rainbow",
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
