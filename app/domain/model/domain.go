package model

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/nao1215/rainbow/utils/errfmt"
)

// Domain is a type that represents a domain name.
type Domain string

// String returns the string representation of Domain.
func (d Domain) String() string {
	return string(d)
}

// Validate validates Domain. If Domain is invalid, it returns an error.
// If domain is empty, it returns nil and the default CloudFront domain will be used.
func (d Domain) Validate() error {
	for _, part := range strings.Split(d.String(), ".") {
		if !isAlphaNumeric(part) {
			return errfmt.Wrap(ErrInvalidDomain, fmt.Sprintf("domain %s is invalid", d))
		}
	}
	return nil
}

// isAlphaNumeric　returns true if s is alphanumeric.
func isAlphaNumeric(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}

// Empty is whether domain is empty
func (d Domain) Empty() bool {
	return d == ""
}

// AllowOrigins is list of origins (domain names) that CloudFront can use as
// the value for the Access-Control-Allow-Origin HTTP response header.
type AllowOrigins []Domain

// Validate validates AllowOrigins. If AllowOrigins is invalid, it returns an error.
func (a AllowOrigins) Validate() (err error) {
	for _, origin := range a {
		if e := origin.Validate(); e != nil {
			err = errors.Join(err, e)
		}
	}
	return err
}

// String returns the string representation of AllowOrigins.
func (a AllowOrigins) String() string {
	origins := make([]string, 0, len(a))
	for _, origin := range a {
		if origin.Empty() {
			continue
		}
		origins = append(origins, origin.String())
	}
	return strings.Join(origins, ",")
}

// Endpoint is a type that represents an endpoint.
type Endpoint string

// String returns the string representation of Endpoint.
func (e Endpoint) String() string {
	return string(e)
}

// Validate validates Endpoint. If Endpoint is invalid, it returns an error.
func (e Endpoint) Validate() error {
	if e == "" {
		return errfmt.Wrap(ErrInvalidEndpoint, "endpoint is empty")
	}

	parsedURL, err := url.Parse(e.String())
	if err != nil {
		return errfmt.Wrap(ErrInvalidDomain, err.Error())
	}
	host := parsedURL.Host
	if host == "" || parsedURL.Scheme == "" {
		return errfmt.Wrap(ErrInvalidDomain, host)
	}
	return nil
}

// DebugLocalstackEndpoint is the endpoint for localstack. It's used for testing.
const DebugLocalstackEndpoint = "http://localhost:4566"
