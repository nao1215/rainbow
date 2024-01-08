// Package domain implements the domain layer. This package is used only for application domain.
package domain

import "errors"

var (
	// ErrInvalidRegion is an error that occurs when the region is invalid.
	ErrInvalidRegion = errors.New("invalid region")
	// ErrEmptyRegion is an error that occurs when the region is empty.
	ErrEmptyRegion = errors.New("region is empty")
	// ErrInvalidBucketName is an error that occurs when the bucket name is invalid.
	ErrInvalidBucketName = errors.New("bucket name is invalid")
	// ErrNoSuchBucket is an error that occurs when the bucket does not exist.
	ErrNoSuchBucket = errors.New("the specified bucket does not exist")
	// ErrInvalidDomain is an error that occurs when the domain is invalid.
	ErrInvalidDomain = errors.New("invalid domain")
	// ErrNotDetectContentType is an error that occurs when the content type cannot be detected.
	ErrNotDetectContentType = errors.New("failed to detect content type")
	// ErrInvalidEndpoint is an error that occurs when the endpoint is invalid.
	ErrInvalidEndpoint = errors.New("invalid endpoint")
	// ErrBucketAlreadyExistsOwnedByOther is an error that occurs when the bucket already exists and is owned by another account.
	ErrBucketAlreadyExistsOwnedByOther = errors.New("bucket already exists and is owned by another account")
	// ErrBucketAlreadyOwnedByYou is an error that occurs when the bucket already exists and is owned by you.
	ErrBucketAlreadyOwnedByYou = errors.New("bucket already exists and is owned by you")
	// ErrBucketPublicAccessBlock is an error that occurs when the bucket public access block setting fails.
	ErrBucketPublicAccessBlock = errors.New("failed to set public access block")
	// ErrBucketPolicySet is an error that occurs when the bucket policy setting fails.
	ErrBucketPolicySet = errors.New("failed to set bucket policy")
	// ErrCDNAlreadyExists is an error that occurs when the CDN already exists.
	ErrCDNAlreadyExists = errors.New("CDN already exists")
	// ErrOriginAccessIdentifyAlreadyExists is an error that occurs when the origin access identify already exists.
	ErrOriginAccessIdentifyAlreadyExists = errors.New("origin access identify already exists")
	// ErrFileUpload is an error that occurs when the file upload fails.
	ErrFileUpload = errors.New("failed to upload file")
)
