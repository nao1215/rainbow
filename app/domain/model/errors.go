package model

import "errors"

var (
	// ErrInvalidRegion is an error that occurs when the region is invalid.
	ErrInvalidRegion = errors.New("invalid region")
	// ErrEmptyRegion is an error that occurs when the region is empty.
	ErrEmptyRegion = errors.New("region is empty")
	// ErrInvalidBucketName is an error that occurs when the bucket name is invalid.
	ErrInvalidBucketName = errors.New("bucket name is invalid")
	// ErrInvalidDomain is an error that occurs when the domain is invalid.
	ErrInvalidDomain = errors.New("invalid domain")
	// ErrNotDetectContentType is an error that occurs when the content type cannot be detected.
	ErrNotDetectContentType = errors.New("failed to detect content type")
	// ErrInvalidEndpoint is an error that occurs when the endpoint is invalid.
	ErrInvalidEndpoint = errors.New("invalid endpoint")
)
