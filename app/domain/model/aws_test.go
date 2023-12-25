package model

import (
	"testing"
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
