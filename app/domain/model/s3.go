// Package model contains the definitions of domain models and business logic.
package model

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/nao1215/rainbow/app/domain"
	"github.com/nao1215/rainbow/utils/errfmt"
	"github.com/nao1215/rainbow/utils/xregex"
	"github.com/wailsapp/mimetype"
)

const (
	// S3DeleteObjectChunksSize is the maximum number of objects that can be deleted in a single request.
	S3DeleteObjectChunksSize = 500
	// MaxS3DeleteObjectsParallelsCount is the maximum number of parallel executions of DeleteObjects.
	MaxS3DeleteObjectsParallelsCount = 5
	// MaxS3DeleteObjectsRetryCount is the maximum number of retries for DeleteObjects.
	MaxS3DeleteObjectsRetryCount = 6
	// S3DeleteObjectsDelayTimeSec is the delay time in seconds.
	S3DeleteObjectsDelayTimeSec = 5
	// MaxS3Keys is the maximum number of keys that can be specified in a single request.
	MaxS3Keys = 1000
)

// DeleteObjectsRetryCount is the number of retries for DeleteObjects.
type DeleteObjectsRetryCount int

// NewDeleteRetryCount creates a new DeleteRetryCount.
// If i is less than 0, it returns 0.
// If i is greater than MaxS3DeleteObjectsRetryCount, it returns MaxS3DeleteObjectsRetryCount.
func NewDeleteRetryCount(i int) DeleteObjectsRetryCount {
	if i < 0 {
		return 0
	}
	if i > MaxS3DeleteObjectsRetryCount {
		return MaxS3DeleteObjectsRetryCount
	}
	return DeleteObjectsRetryCount(i)
}

// Region is the name of the AWS region.
type Region string

const (
	// RegionUSEast1 US East (N. Virginia)
	RegionUSEast1 Region = "us-east-1"
	// RegionUSEast2 US East (Ohio)
	RegionUSEast2 Region = "us-east-2"
	// RegionUSWest1 US West (N. California)
	RegionUSWest1 Region = "us-west-1"
	// RegionUSWest2 US West (Oregon)
	RegionUSWest2 Region = "us-west-2"
	// RegionAFSouth1 Africa (Cape Town)
	RegionAFSouth1 Region = "af-south-1"
	// RegionAPEast1 Asia Pacific (Hong Kong)
	RegionAPEast1 Region = "ap-east-1"
	// RegionAPSouth1 Asia Pacific (Mumbai)
	RegionAPSouth1 Region = "ap-south-1"
	// RegionAPNortheast1 Asia Pacific (Tokyo)
	RegionAPNortheast1 Region = "ap-northeast-1"
	// RegionAPNortheast2 Asia Pacific (Seoul)
	RegionAPNortheast2 Region = "ap-northeast-2"
	// RegionAPNortheast3 Asia Pacific (Osaka-Local)
	RegionAPNortheast3 Region = "ap-northeast-3"
	// RegionAPSoutheast1 Asia Pacific (Singapore)
	RegionAPSoutheast1 Region = "ap-southeast-1"
	// RegionAPSoutheast2 Asia Pacific (Sydney)
	RegionAPSoutheast2 Region = "ap-southeast-2"
	// RegionCACentral1 Canada (Central)
	RegionCACentral1 Region = "ca-central-1"
	// RegionCNNorth1 China (Beijing)
	RegionCNNorth1 Region = "cn-north-1"
	// RegionCNNorthwest1 China (Ningxia)
	RegionCNNorthwest1 Region = "cn-northwest-1"
	// RegionEUCentral1 Europe (Frankfurt)
	RegionEUCentral1 Region = "eu-central-1"
	// RegionEUNorth1 Europe (Stockholm)
	RegionEUNorth1 Region = "eu-north-1"
	// RegionEUSouth1 Europe (Milan)
	RegionEUSouth1 Region = "eu-south-1"
	// RegionEUWest1 Europe (Ireland)
	RegionEUWest1 Region = "eu-west-1"
	// RegionEUWest2 Europe (London)
	RegionEUWest2 Region = "eu-west-2"
	// RegionEUWest3 Europe (Paris)
	RegionEUWest3 Region = "eu-west-3"
	// RegionMESouth1 Middle East (Bahrain)
	RegionMESouth1 Region = "me-south-1"
	// RegionSASouth1 South America (São Paulo)
	RegionSASouth1 Region = "sa-south-1"
	// RegionUSGovEast1 AWS GovCloud (US-East)
	RegionUSGovEast1 Region = "us-gov-east-1"
	// RegionUSGovWest1 AWS GovCloud (US)
	RegionUSGovWest1 Region = "us-gov-west-1"
)

var regions = []Region{
	RegionUSEast1, RegionUSEast2, RegionUSWest1, RegionUSWest2, RegionAFSouth1, RegionAPEast1,
	RegionAPSouth1, RegionAPNortheast1, RegionAPNortheast2, RegionAPNortheast3, RegionAPSoutheast1,
	RegionAPSoutheast2, RegionCACentral1, RegionCNNorth1, RegionCNNorthwest1, RegionEUCentral1,
	RegionEUNorth1, RegionEUSouth1, RegionEUWest1, RegionEUWest2, RegionEUWest3, RegionMESouth1,
	RegionSASouth1, RegionUSGovEast1, RegionUSGovWest1,
}

// Validate returns true if the Region exists.
func (r Region) Validate() error {
	switch r {
	case
		RegionUSEast1, RegionUSEast2, RegionUSWest1, RegionUSWest2, RegionAFSouth1,
		RegionAPEast1, RegionAPSouth1, RegionAPNortheast1, RegionAPNortheast2,
		RegionAPNortheast3, RegionAPSoutheast1, RegionAPSoutheast2, RegionCACentral1,
		RegionCNNorth1, RegionCNNorthwest1, RegionEUCentral1, RegionEUNorth1,
		RegionEUSouth1, RegionEUWest1, RegionEUWest2, RegionEUWest3, RegionMESouth1,
		RegionSASouth1, RegionUSGovEast1, RegionUSGovWest1:
		return nil
	case Region(""):
		return domain.ErrEmptyRegion
	default:
		return domain.ErrInvalidRegion
	}
}

// String returns the string representation of the Region.
func (r Region) String() string {
	return string(r)
}

// Next returns the next region.
// If the region is the last one, it returns the first region.
// If the region is invalid, it returns "ap-northeast-1".
func (r Region) Next() Region {
	for i, region := range regions {
		if r == region {
			if i == len(regions)-1 {
				return regions[0]
			}
			return regions[i+1]
		}
	}
	return RegionAPNortheast1
}

// Prev returns the previous region.
// If the region is the first one, it returns the last region.
// If the region is invalid, it returns "ap-northeast-1".
func (r Region) Prev() Region {
	for i, region := range regions {
		if r == region {
			if i == 0 {
				return regions[len(regions)-1]
			}
			return regions[i-1]
		}
	}
	return RegionAPNortheast1
}

const (
	// MinBucketNameLength is the minimum length of the bucket name.
	MinBucketNameLength = 3
	// MaxBucketNameLength is the maximum length of the bucket name.
	MaxBucketNameLength = 63
	// S3Protocol is the protocol of the S3 bucket.
	S3Protocol = "s3://"
)

// Bucket is the name of the S3 bucket.
type Bucket string

// NewBucketWithoutProtocol creates a new Bucket.
func NewBucketWithoutProtocol(s string) Bucket {
	return Bucket(strings.TrimPrefix(s, S3Protocol))
}

// WithProtocol returns the Bucket with the protocol.
func (b Bucket) WithProtocol() Bucket {
	if strings.HasPrefix(b.String(), S3Protocol) {
		return b
	}
	return Bucket(S3Protocol + b.String())
}

// Join returns the Bucket with the S3Key.
// e.g. "bucket" + "key" -> "bucket/key"
func (b Bucket) Join(key S3Key) Bucket {
	if b.Empty() || key.Empty() {
		return b
	}
	if strings.HasSuffix(key.String(), "/") {
		key = S3Key(strings.TrimSuffix(key.String(), "/"))
	}
	return Bucket(fmt.Sprintf("%s/%s", b.String(), key.String()))
}

// String returns the string representation of the Bucket.
func (b Bucket) String() string {
	return string(b)
}

// Empty is whether bucket name is empty
func (b Bucket) Empty() bool {
	return b == ""
}

// Domain returns the domain name of the Bucket.
func (b Bucket) Domain() string {
	return fmt.Sprintf("%s.s3.amazonaws.com", b.String())
}

// TrimKey returns the Bucket without the key.
// e.g. "bucket/key" -> "bucket"
func (b Bucket) TrimKey() Bucket {
	return Bucket(strings.Split(b.String(), "/")[0])
}

// Split returns the Bucket and the S3Key.
// If the Bucket does not contain "/", the S3Key is empty.
func (b Bucket) Split() (Bucket, S3Key) {
	s := strings.Split(b.String(), "/")
	if len(s) == 1 {
		return b, ""
	}

	key := strings.Join(s[1:], "/")
	if key == "" {
		return Bucket(s[0]), S3Key("")
	}
	return Bucket(s[0]), S3Key(filepath.Clean(key))
}

// Validate returns true if the Bucket is valid.
// Bucket naming rules: https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html
func (b Bucket) Validate() error {
	if b.Empty() {
		return errfmt.Wrap(domain.ErrInvalidBucketName, "s3 bucket name is empty")
	}

	validators := []func() error{
		b.validateLength,
		b.validatePattern,
		b.validatePrefix,
		b.validateSuffix,
		b.validateCharSequence,
	}
	for _, v := range validators {
		if err := v(); err != nil {
			return err
		}
	}
	return nil
}

// validateLength validates the length of the bucket name.
func (b Bucket) validateLength() error {
	if len(b) < MinBucketNameLength || len(b) > MaxBucketNameLength {
		return fmt.Errorf("s3 bucket name must be between 3 and 63 characters long")
	}
	return nil
}

var s3RegexPattern xregex.Regex //nolint:gochecknoglobals

// validatePattern validates the pattern of the bucket name.
func (b Bucket) validatePattern() error {
	s3RegexPattern.InitOnce(`^[a-z0-9][a-z0-9.-]*[a-z0-9]$`)
	if err := s3RegexPattern.MatchString(string(b)); err != nil {
		return errfmt.Wrap(domain.ErrInvalidBucketName, "s3 bucket name must use only lowercase letters, numbers, periods, and hyphens")
	}
	return nil
}

// validatePrefix validates the prefix of the bucket name.
func (b Bucket) validatePrefix() error {
	for _, prefix := range []string{"xn--", "sthree-", "sthree-configurator"} {
		if strings.HasPrefix(string(b), prefix) {
			return errfmt.Wrap(domain.ErrInvalidBucketName, "s3 bucket name must not start with \"xn--\", \"sthree-\", or \"sthree-configurator\"")
		}
	}
	return nil
}

// validateSuffix validates the suffix of the bucket name.
func (b Bucket) validateSuffix() error {
	for _, suffix := range []string{"-s3alias", "--ol-s3"} {
		if strings.HasSuffix(string(b), suffix) {
			return errfmt.Wrap(domain.ErrInvalidBucketName, "s3 bucket name must not end with \"-s3alias\" or \"--ol-s3\"")
		}
	}
	return nil
}

// validateCharSequence validates the character sequence of the bucket name.
func (b Bucket) validateCharSequence() error {
	if strings.Contains(string(b), "..") || strings.Contains(string(b), "--") {
		return errfmt.Wrap(domain.ErrInvalidBucketName, "s3 bucket name must not contain consecutive periods or hyphens")
	}
	return nil
}

// BucketSets is the set of the BucketSet.
type BucketSets []BucketSet

// Len returns the length of the BucketSets.
func (b BucketSets) Len() int {
	return len(b)
}

// Empty returns true if the BucketSets is empty.
func (b BucketSets) Empty() bool {
	return b.Len() == 0
}

// Contains returns true if the BucketSets contains the bucket.
func (b BucketSets) Contains(bucket Bucket) bool {
	for _, bs := range b {
		if bs.Bucket == bucket {
			return true
		}
	}
	return false
}

// BucketSet is the set of the Bucket and the Region.
type BucketSet struct {
	// Bucket is the name of the S3 bucket.
	Bucket Bucket
	// Region is the name of the AWS region.
	Region Region
	// CreationDate is date the bucket was created.
	// This date can change when making changes to your bucket, such as editing its bucket policy.
	CreationDate time.Time
}

// S3ObjectIdentifiers is the set of the S3ObjectSet.
type S3ObjectIdentifiers []S3ObjectIdentifier

// Len returns the length of the S3ObjectIdentifiers.
func (s S3ObjectIdentifiers) Len() int {
	return len(s)
}

// Less defines the ordering of S3ObjectIdentifier instances.
func (s S3ObjectIdentifiers) Less(i, j int) bool {
	return s[i].S3Key < s[j].S3Key
}

// Swap swaps the elements with indexes i and j.
func (s S3ObjectIdentifiers) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// ToS3ObjectIdentifiers converts the S3ObjectSets to the ObjectIdentifiers.
func (s S3ObjectIdentifiers) ToS3ObjectIdentifiers() []types.ObjectIdentifier {
	ids := make([]types.ObjectIdentifier, 0, s.Len())
	for _, o := range s {
		ids = append(ids, *o.ToAWSS3ObjectIdentifier())
	}
	return ids
}

// S3ObjectIdentifier is the object  identifier in the S3 bucket.
type S3ObjectIdentifier struct {
	// S3Key is the name of the object.
	S3Key S3Key
	// VersionID is the version ID for the specific version of the object to delete.
	VersionID VersionID
}

// ToAWSS3ObjectIdentifier converts the S3ObjectIdentifier to the ObjectIdentifier.
func (o S3ObjectIdentifier) ToAWSS3ObjectIdentifier() *types.ObjectIdentifier {
	return &types.ObjectIdentifier{
		Key:       aws.String(o.S3Key.String()),
		VersionId: aws.String(o.VersionID.String()),
	}
}

// S3Key is the name of the object.
// Replacement must be made for object keys containing special characters (such as carriage returns) when using XML requests.
// For more information, see XML related object key constraints (https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-keys.html#object-key-xml-related-constraints).
type S3Key string

// String returns the string representation of the S3Key.
func (k S3Key) String() string {
	return string(k)
}

// Empty is whether S3Key is empty
func (k S3Key) Empty() bool {
	return k == ""
}

// IsAll is whether S3Key is "*"
func (k S3Key) IsAll() bool {
	return k == "*"
}

func (k S3Key) Join(key S3Key) S3Key {
	if key.Empty() {
		return k
	}
	if strings.HasPrefix(key.String(), "/") {
		key = S3Key(strings.TrimPrefix(key.String(), "/"))
	}
	if strings.HasSuffix(key.String(), "/") {
		key = S3Key(strings.TrimSuffix(key.String(), "/"))
	}
	if k.Empty() {
		return key
	}
	return S3Key(fmt.Sprintf("%s/%s", k.String(), key))
}

// VersionID is the version ID for the specific version of the object to delete.
// This functionality is not supported for directory buckets.
type VersionID string

// String returns the string representation of the VersionID.
func (v VersionID) String() string {
	return string(v)
}

// S3Object is the object in the S3 bucket.
type S3Object struct {
	*bytes.Buffer
}

// NewS3Object creates a new S3Object.
func NewS3Object(b []byte) *S3Object {
	return &S3Object{Buffer: bytes.NewBuffer(b)}
}

// ToFile writes the S3Object to the file.
func (s *S3Object) ToFile(path string, perm fs.FileMode) error {
	return os.WriteFile(filepath.Clean(path), s.Bytes(), perm)
}

// ContentType returns the content type of the S3Object.
// If the content type cannot be detected, it returns "plain/text".
func (s *S3Object) ContentType() string {
	mtype, err := mimetype.DetectReader(s.Buffer)
	if err != nil {
		return "plain/text"
	}
	return mtype.String()
}

// ContentLength returns the content length of the S3Object.
func (s *S3Object) ContentLength() int64 {
	return int64(s.Len())
}
