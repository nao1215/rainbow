package external

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/nao1215/rainbow/app/domain/model"
)

// newS3Session returns a new session.
func newS3Session(profile model.AWSProfile, region model.Region, endpoint *model.Endpoint) *session.Session {
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable, // Ref. ~/.aws/config
		Profile:           profile.String(),
	}))

	session.Config.Region = aws.String(region.String())
	if endpoint != nil {
		// If you want to debug, uncomment the following lines.
		// session.Config.WithLogLevel(aws.LogDebugWithHTTPBody)
		session.Config.S3ForcePathStyle = aws.Bool(true)
		session.Config.Endpoint = aws.String(endpoint.String())
		session.Config.DisableSSL = aws.Bool(true)
	}
	return session
}
