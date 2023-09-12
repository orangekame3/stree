package pkg

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ddddddO/gtree"
)

type S3Config struct {
	AwsProfile  string
	AwsRegion   string
	EndpointURL string
	Local       bool
	NoColor     bool
}

func InitializeAWSSession(config S3Config) *s3.S3 {
	var sess *session.Session
	if config.Local {
		sessOptions := session.Options{
			Profile: config.AwsProfile,
			Config: aws.Config{
				Region:           aws.String(config.AwsRegion),
				Endpoint:         aws.String(config.EndpointURL),
				S3ForcePathStyle: aws.Bool(true),
			},
		}
		sess = session.Must(session.NewSessionWithOptions(sessOptions))
	} else {
		sessOptions := session.Options{
			Profile: config.AwsProfile,
		}
		if config.AwsRegion != "" {
			sessOptions.Config = aws.Config{Region: aws.String(config.AwsRegion)}
		}
		sess = session.Must(session.NewSessionWithOptions(sessOptions))
	}

	s3Svc := s3.New(sess)
	return s3Svc
}

func FetchS3Objects(s3Svc *s3.S3, bucket string, prefix string, root *gtree.Node, noColor bool) error {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}

	err := s3Svc.ListObjectsV2Pages(input, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range page.Contents {
			keys := strings.Split(*obj.Key, "/")
			AddNodeWithColor(root, keys, 0, noColor)
		}
		return !lastPage
	})
	return err
}
