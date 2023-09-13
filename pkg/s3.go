// Package: pkg is a package that contains the business logic for stree
package pkg

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Config struct {
	AwsProfile  string
	AwsRegion   string
	EndpointURL string
	Local       bool
}

// InitializeAWSSession returns an AWS session based on the provided configuration
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
		return s3.New(sess)
	}

	sessOptions := session.Options{
		Profile: config.AwsProfile,
	}
	if config.AwsRegion != "" {
		sessOptions.Config = aws.Config{Region: aws.String(config.AwsRegion)}
	}
	sess = session.Must(session.NewSessionWithOptions(sessOptions))

	return s3.New(sess)
}

// FetchS3ObjectKeys returns a slice of keys for all objects in the specified bucket and prefix
func FetchS3ObjectKeys(s3Svc *s3.S3, bucket string, prefix string) ([][]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}
	var keys [][]string

	pageHandler := func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range page.Contents {
			key := strings.Split(*obj.Key, "/")
			keys = append(keys, key)
		}
		return !lastPage
	}

	if err := s3Svc.ListObjectsV2Pages(input, pageHandler); err != nil {
		return nil, err
	}

	return keys, nil
}
