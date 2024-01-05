package model

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestNewAWSProfile(t *testing.T) { //nolint
	type args struct {
		p string
	}
	tests := []struct {
		name string
		args args
		want AWSProfile
	}{
		{
			name: "success",
			args: args{
				p: "test",
			},
			want: AWSProfile("test"),
		},
		{
			name: "success. p is empty",
			args: args{
				p: "",
			},
			want: AWSProfile("from env"),
		},
		{
			name: "success. p is empty and $AWS_PROFILE is empty",
			args: args{
				p: "",
			},
			want: AWSProfile("default"),
		},
	}
	for _, tt := range tests { //nolint
		if tt.name == "success. p is empty" {
			t.Setenv("AWS_PROFILE", "from env")
		} else if tt.name == "success. p is empty and $AWS_PROFILE is empty" {
			t.Setenv("AWS_PROFILE", "")
		}

		t.Run(tt.name, func(t *testing.T) {
			if got := NewAWSProfile(tt.args.p); got != tt.want {
				t.Errorf("NewAWSProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAWSProfileString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		p    AWSProfile
		want string
	}{
		{
			name: "success",
			p:    AWSProfile("test"),
			want: "test",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.p.String(); got != tt.want {
				t.Errorf("AWSProfile.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAWSConfig_Region(t *testing.T) {
	t.Parallel()

	type fields struct {
		Config *aws.Config
	}
	tests := []struct {
		name   string
		fields fields
		want   Region
	}{
		{
			name: "If aws config region is ap-northeast-1, return RegionAPNortheast1",
			fields: fields{
				Config: &aws.Config{
					Region: string(RegionAPNortheast1),
				},
			},
			want: RegionAPNortheast1,
		},
		{
			name: "If aws config region isempty, return RegionUSEast1",
			fields: fields{
				Config: &aws.Config{
					Region: "",
				},
			},
			want: RegionUSEast1,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &AWSConfig{
				Config: tt.fields.Config,
			}
			if got := c.Region(); got != tt.want {
				t.Errorf("AWSConfig.Region() = %v, want %v", got, tt.want)
			}
		})
	}
}
