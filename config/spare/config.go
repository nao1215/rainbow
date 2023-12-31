package spare

import (
	"errors"
	"fmt"
	"io"

	"github.com/charmbracelet/log"
	"github.com/k1LoW/runn/version"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/config"
	"github.com/nao1215/rainbow/utils/errfmt"
	"github.com/nao1215/spare/utils/xrand"
	"gopkg.in/yaml.v2"
)

// ConfigFilePath is the path of the configuration file.
const ConfigFilePath string = ".spare.yml"

// Config is a struct that corresponds to the configuration file ".spare.yml".
type Config struct {
	SpareTemplateVersion TemplateVersion `yaml:"spareTemplateVersion"`
	// DeployTarget is the path of the deploy target (it's SPA).
	DeployTarget DeployTarget `yaml:"deployTarget"`
	// Region is AWS region.
	Region model.Region `yaml:"region"`
	// CustomDomain is the domain name of the CloudFront.
	// If you do not specify this, the CloudFront default domain name is used.
	CustomDomain model.Domain `yaml:"customDomain"`
	// S3Bucket is the name of the S3 bucket.
	S3Bucket model.Bucket `yaml:"s3BucketName"` //nolint
	// AllowOrigins is the list of domains that are allowed to access the SPA.
	AllowOrigins            model.AllowOrigins `yaml:"allowOrigins"`
	DebugLocalstackEndpoint model.Endpoint     `yaml:"debugLocalstackEndpoint"`
	// TODO: WAF, HTTPS, Cache
}

// NewConfig returns a new Config.
func NewConfig() *Config {
	cfg := &Config{
		SpareTemplateVersion:    CurrentSpareTemplateVersion,
		DeployTarget:            "src",
		Region:                  model.RegionUSEast1,
		CustomDomain:            "",
		S3Bucket:                "",
		AllowOrigins:            model.AllowOrigins{},
		DebugLocalstackEndpoint: model.DebugLocalstackEndpoint,
	}
	cfg.S3Bucket = cfg.DefaultS3BucketName()
	return cfg
}

// DefaultS3BucketName returns the default S3 bucket name.
func (c *Config) DefaultS3BucketName() model.Bucket {
	const randomStrLen = 15
	randomID, err := xrand.RandomLowerAlphanumericStr(randomStrLen)
	if err != nil {
		log.Error(err)
		randomID = "cannot-generate-random-id"
	}

	return model.Bucket(
		fmt.Sprintf("%s-%s-%s", version.Name, c.Region, randomID))
}

// Write writes the Config to the io.Writer.
func (c *Config) Write(w io.Writer) (err error) {
	encoder := yaml.NewEncoder(w)
	defer func() {
		if closeErr := encoder.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()
	return encoder.Encode(c)
}

// Read reads the Config from the io.Reader.
func (c *Config) Read(r io.Reader) (err error) {
	decoder := yaml.NewDecoder(r)
	return decoder.Decode(c)
}

// Validate validates the Config.
// If debugMode is true, it validates the DebugLocalstackEndpoint.
func (c *Config) Validate(debugMode bool) error {
	validators := []model.Validator{
		c.SpareTemplateVersion,
		c.DeployTarget,
		c.Region,
		c.CustomDomain,
		c.S3Bucket,
		c.AllowOrigins,
	}
	if debugMode {
		validators = append(validators, c.DebugLocalstackEndpoint)
	}

	for _, v := range validators {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// TemplateVersion is a type that represents a spare template version.
type TemplateVersion string

// CurrentSpareTemplateVersion is the version of the template.
const CurrentSpareTemplateVersion TemplateVersion = "0.0.1"

// String returns the string representation of TemplateVersion.
func (t TemplateVersion) String() string {
	return string(t)
}

// Validate validates TemplateVersion. If TemplateVersion is invalid, it returns an error.
// TemplateVersion is invalid if it is empty.
func (t TemplateVersion) Validate() error {
	if t == "" {
		return errfmt.Wrap(config.ErrInvalidSpareTemplateVersion, "SpareTemplateVersion is empty")
	}
	return nil
}

// DeployTarget is a type that represents a deploy target path.
type DeployTarget string

// String returns the string representation of DeployTarget.
func (d DeployTarget) String() string {
	return string(d)
}

// Validate validates DeployTarget. If DeployTarget is invalid, it returns an error.
// DeployTarget is invalid if it is empty.
func (d DeployTarget) Validate() error {
	if d == "" {
		return errfmt.Wrap(config.ErrInvalidDeployTarget, "DeployTarget is empty")
	}
	// TODO: check if the path exists
	return nil
}
