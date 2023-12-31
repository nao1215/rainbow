// Package service is an abstraction layer for accessing external services.
package service

import "errors"

var (
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
	// ErrNotDetectContentType is an error that occurs when the content type cannot be detected.
	ErrNotDetectContentType = errors.New("failed to detect content type")
	// ErrFileUpload is an error that occurs when the file upload fails.
	ErrFileUpload = errors.New("failed to upload file")
)
