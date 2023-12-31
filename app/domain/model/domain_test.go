package model

import (
	"errors"
	"testing"
)

const (
	exampleCom             = "example.com"
	exampleNet             = "example.net"
	exampleComWithProtocol = "https://example.com"
)

func TestDomainString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    Domain
		want string
	}{
		{
			name: exampleCom,
			d:    exampleCom,
			want: exampleCom,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.d.String(); got != tt.want {
				t.Errorf("Domain.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		d       Domain
		wantErr error
	}{
		{
			name:    "success",
			d:       exampleCom,
			wantErr: nil,
		},
		{
			name:    "failure. protocol is included",
			d:       exampleComWithProtocol,
			wantErr: ErrInvalidDomain,
		},
		{
			name:    "success. domain is empty",
			d:       "",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.d.Validate(); !errors.Is(err, tt.wantErr) {
				t.Errorf("Domain.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDomainEmpty(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		d    Domain
		want bool
	}{
		{
			name: "success",
			d:    exampleCom,
			want: false,
		},
		{
			name: "failure",
			d:    "",
			want: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.d.Empty(); got != tt.want {
				t.Errorf("Domain.Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsAlphaNumeric(t *testing.T) {
	t.Parallel()

	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "success",
			args: args{s: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"},
			want: true,
		},
		{
			name: "failure",
			args: args{s: "abc123/"},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := isAlphaNumeric(tt.args.s); got != tt.want {
				t.Errorf("isAlphaNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllowOriginsValidate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		a       AllowOrigins
		wantErr bool
	}{
		{
			name:    "success",
			a:       AllowOrigins{exampleCom, exampleNet},
			wantErr: false,
		},
		{
			name:    "success. include empty string",
			a:       AllowOrigins{exampleCom, ""},
			wantErr: false,
		},
		{
			name:    "failure. origin is invalid",
			a:       AllowOrigins{exampleCom, exampleComWithProtocol},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.a.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("AllowOrigins.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEndpointString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		e    Endpoint
		want string
	}{
		{
			name: "success",
			e:    Endpoint(exampleComWithProtocol),
			want: exampleComWithProtocol,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Endpoint.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndpointValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		e       Endpoint
		wantErr bool
	}{
		{
			name:    "success",
			e:       Endpoint(exampleComWithProtocol),
			wantErr: false,
		},
		{
			name:    "failure. protocol is not included",
			e:       exampleCom,
			wantErr: true,
		},
		{
			name:    "failure. endpoint is empty",
			e:       "",
			wantErr: true,
		},
		{
			name:    "failure. include Ctrl character",
			e:       "\x00",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.e.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Endpoint.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAllowOriginsString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		a    AllowOrigins
		want string
	}{
		{
			name: "success",
			a:    AllowOrigins{exampleCom, exampleNet},
			want: "example.com,example.net",
		},
		{
			name: "success. include empty string",
			a:    AllowOrigins{exampleCom, ""},
			want: "example.com",
		},
		{
			name: "success. empty",
			a:    AllowOrigins{},
			want: "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.a.String(); got != tt.want {
				t.Errorf("AllowOrigins.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
