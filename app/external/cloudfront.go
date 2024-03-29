package external

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/domain/service"
	"github.com/nao1215/rainbow/utils/errfmt"
)

// NewCloudFrontClient returns a new CloudFront client.
func NewCloudFrontClient(cfg *model.AWSConfig) *cloudfront.Client {
	return cloudfront.NewFromConfig(*cfg.Config)
}

// CloudFrontCreatorSet is a provider set for CloudFrontCreator.
var CloudFrontCreatorSet = wire.NewSet(
	wire.Bind(new(service.CloudFrontCreator), new(*CloudFrontCreator)),
	NewCloudFrontCreator,
)

// CloudFrontCreator is an implementation for CloudFrontCreator.
type CloudFrontCreator struct {
	*cloudfront.Client
}

var _ service.CloudFrontCreator = &CloudFrontCreator{}

// NewCloudFrontCreator creates a new CloudFrontCreator.
func NewCloudFrontCreator(c *cloudfront.Client) *CloudFrontCreator {
	return &CloudFrontCreator{
		Client: c,
	}
}

// CreateCloudFront creates a CDN.
func (c *CloudFrontCreator) CreateCloudFront(ctx context.Context, input *service.CloudFrontCreatorInput) (*service.CloudFrontCreatorOutput, error) {
	createDistributionInput := &cloudfront.CreateDistributionInput{
		DistributionConfig: &types.DistributionConfig{
			Comment:         aws.String("CloudFront Distribution Generated by Rainbow Project"),
			CallerReference: aws.String(uuid.New().String()),
			DefaultCacheBehavior: &types.DefaultCacheBehavior{
				TargetOriginId:       aws.String("S3 Origin ID Generated by Rainbow Project"),
				ViewerProtocolPolicy: types.ViewerProtocolPolicyRedirectToHttps,
				MinTTL:               aws.Int64(300), //nolint:gomnd
				MaxTTL:               aws.Int64(300), //nolint:gomnd
				DefaultTTL:           aws.Int64(300), //nolint:gomnd
				AllowedMethods: &types.AllowedMethods{
					Items: []types.Method{
						types.Method("GET"),
						types.Method("HEAD"),
						types.Method("OPTIONS"),
					},
					Quantity: aws.Int32(3),
					CachedMethods: &types.CachedMethods{
						Items: []types.Method{
							types.Method("GET"),
							types.Method("HEAD"),
						},
						Quantity: aws.Int32(2), //nolint:gomnd
					},
				},
				// Deprecated fields
				ForwardedValues: &types.ForwardedValues{
					QueryString: aws.Bool(true),
					Cookies: &types.CookiePreference{
						Forward: types.ItemSelection("none"),
					},
				},
			},
			DefaultRootObject: aws.String("index.html"),
			HttpVersion:       types.HttpVersion("http2and3"),
			PriceClass:        types.PriceClass("PriceClass_100"),
			Origins: &types.Origins{
				Items: []types.Origin{
					{
						Id:         aws.String("S3 Origin ID Generated by Spare"),
						DomainName: aws.String(input.Bucket.Domain()),
						S3OriginConfig: &types.S3OriginConfig{
							OriginAccessIdentity: aws.String(
								fmt.Sprintf("origin-access-identity/cloudfront/%s", *input.OAIID),
							),
						},
					},
				},
				Quantity: aws.Int32(1),
			},
			Enabled: aws.Bool(true),
		},
	}

	output, err := c.CreateDistribution(ctx, createDistributionInput)
	if err != nil {
		return nil, errfmt.Wrap(err, "failed to create a CloudFront distribution")
	}
	return &service.CloudFrontCreatorOutput{
		Domain: model.Domain(*output.Distribution.DomainName),
	}, nil
}

// OAICreatorSet is a provider set for OAICreator.
var OAICreatorSet = wire.NewSet(
	NewCloudFrontOAICreator,
	wire.Bind(new(service.OAICreator), new(*CloudFrontOAICreator)),
)

// CloudFrontOAICreator is an implementation for OAICreator.
type CloudFrontOAICreator struct {
	*cloudfront.Client
}

var _ service.OAICreator = &CloudFrontOAICreator{}

// NewCloudFrontOAICreator creates a new CloudFrontOAICreator.
func NewCloudFrontOAICreator(c *cloudfront.Client) *CloudFrontOAICreator {
	return &CloudFrontOAICreator{
		Client: c,
	}
}

// CreateOAI creates a new OAI.
func (c *CloudFrontOAICreator) CreateOAI(ctx context.Context, _ *service.OAICreatorInput) (*service.OAICreatorOutput, error) {
	createOAIInput := &cloudfront.CreateCloudFrontOriginAccessIdentityInput{
		CloudFrontOriginAccessIdentityConfig: &types.CloudFrontOriginAccessIdentityConfig{
			CallerReference: aws.String(uuid.NewString()),
			Comment:         aws.String("Origin Access Identity (OAI) Generated by Spare"),
		},
	}

	output, err := c.CreateCloudFrontOriginAccessIdentity(ctx, createOAIInput)
	if err != nil {
		return nil, err
	}

	return &service.OAICreatorOutput{
		ID: output.CloudFrontOriginAccessIdentity.Id,
	}, nil
}
