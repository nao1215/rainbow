// Package model contains the definitions of domain models and business logic.
package model

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

// Valid returns true if the Region exists.
func (r Region) Valid() bool {
	switch r {
	case
		RegionUSEast1, RegionUSEast2, RegionUSWest1, RegionUSWest2, RegionAFSouth1,
		RegionAPEast1, RegionAPSouth1, RegionAPNortheast1, RegionAPNortheast2,
		RegionAPNortheast3, RegionAPSoutheast1, RegionAPSoutheast2, RegionCACentral1,
		RegionCNNorth1, RegionCNNorthwest1, RegionEUCentral1, RegionEUNorth1,
		RegionEUSouth1, RegionEUWest1, RegionEUWest2, RegionEUWest3, RegionMESouth1,
		RegionSASouth1, RegionUSGovEast1, RegionUSGovWest1:
		return true
	default:
		return false
	}
}

// String returns the string representation of the Region.
func (r Region) String() string {
	return string(r)
}

// Bucket is the name of the S3 bucket.
type Bucket string

// Valid returns true if the Bucket is valid (it's not empty).
func (b Bucket) Valid() bool {
	// TODO: check strictly
	return b != ""
}

// String returns the string representation of the Bucket.
func (b Bucket) String() string {
	return string(b)
}
