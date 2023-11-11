// Package pkg provides the core functionality of the program.
package pkg

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Config represents the configuration for an S3 session
type S3Config struct {
	AwsProfile  string
	AwsRegion   string
	EndpointURL string
	Local       bool
	MFA         bool
}

// InitializeAWSSession returns an AWS session based on the provided configuration
func InitializeAWSSession(config S3Config) *s3.S3 {
	var sess *session.Session
	if config.Local {
		sessOptions := session.Options{
			Config: aws.Config{
				Region:           aws.String("us-east-1"),
				Endpoint:         aws.String("http://localhost:4566"),
				S3ForcePathStyle: aws.Bool(true),
				Credentials:      credentials.NewStaticCredentials("dummy", "dummy", ""),
			},
		}
		// override region and endpoint
		if config.AwsRegion != "" {
			sessOptions.Config.Region = aws.String(config.AwsRegion)
		}
		if config.EndpointURL != "" {
			sessOptions.Config.Endpoint = aws.String(config.EndpointURL)
		}
		sess = session.Must(session.NewSessionWithOptions(sessOptions))
		return s3.New(sess)
	}

	sessOptions := session.Options{
		Profile:           config.AwsProfile,
		SharedConfigState: session.SharedConfigEnable,
	}
	// override region
	if config.AwsRegion != "" {
		sessOptions.Config.Region = aws.String(config.AwsRegion)
	}
	if config.MFA {
		sessOptions.AssumeRoleTokenProvider = stscreds.StdinTokenProvider
	}
	sess = session.Must(session.NewSessionWithOptions(sessOptions))
	return s3.New(sess)
}

// FetchS3ObjectKeys returns a slice of keys for all objects in the specified bucket and prefix
func FetchS3ObjectKeys(s3Svc *s3.S3, bucket string, prefix string, maxDepth *int) ([][]string, error) {
	var delimiter *string
	if maxDepth != nil {
		delimiter = aws.String("/")
	}
	queue := []string{prefix}
	depth := []int{0}
	queued := map[string]struct{}{}

	var keys [][]string
	for len(queue) > 0 {
		currentPrefix := queue[0]
		currentDepth := depth[0]
		queue = queue[1:]
		depth = depth[1:]

		if maxDepth != nil && currentDepth >= *maxDepth {
			key := strings.Split(currentPrefix, "/")
			keys = append(keys, key)
			continue
		}

		input := &s3.ListObjectsV2Input{
			Bucket:    aws.String(bucket),
			Prefix:    aws.String(currentPrefix),
			Delimiter: delimiter,
		}

		pageHandler := func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			for _, obj := range page.Contents {
				key := strings.Split(*obj.Key, "/")
				keys = append(keys, key)
			}
			if maxDepth != nil {
				for _, commonPrefix := range page.CommonPrefixes {
					if _, ok := queued[*commonPrefix.Prefix]; ok {
						continue
					}
					queue = append(queue, *commonPrefix.Prefix)
					depth = append(depth, currentDepth+1)
				}
			}
			return !lastPage
		}

		if err := s3Svc.ListObjectsV2Pages(input, pageHandler); err != nil {
			return nil, err
		}
	}
	return keys, nil
}
