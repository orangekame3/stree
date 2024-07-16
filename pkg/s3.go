// Package pkg provides the core functionality of the program.
package pkg

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
func InitializeAWSSession(conf S3Config) (*s3.Client, error) {
	var awsConfig aws.Config
	var err error

	if conf.Local {
		awsConfig, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion("us-east-1"),
			config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           "http://localhost:4566",
					SigningRegion: "us-east-1",
				}, nil
			})),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")),
		)
		if conf.AwsRegion != "" {
			awsConfig.Region = conf.AwsRegion
		}
		if conf.EndpointURL != "" {
			awsConfig.EndpointResolver = aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           conf.EndpointURL,
					SigningRegion: conf.AwsRegion,
				}, nil
			})
		}
	} else {
		loadOptions := []func(*config.LoadOptions) error{
			config.WithSharedConfigProfile(conf.AwsProfile),
		}
		if conf.AwsRegion != "" {
			loadOptions = append(loadOptions, config.WithRegion(conf.AwsRegion))
		}
		if conf.MFA {
			loadOptions = append(loadOptions, config.WithAssumeRoleCredentialOptions(func(options *stscreds.AssumeRoleOptions) {
				options.TokenProvider = stscreds.StdinTokenProvider
			}))
		}
		awsConfig, err = config.LoadDefaultConfig(context.TODO(), loadOptions...)
	}

	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(awsConfig), nil
}

// FetchS3ObjectKeys returns a slice of keys for all objects in the specified bucket and prefix
func FetchS3ObjectKeys(s3Client *s3.Client, bucket string, prefix string, maxDepth *int) ([][]string, error) {
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
		fmt.Println("queue", queue)
		fmt.Println("currentPrefix", currentPrefix)

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

		paginator := s3.NewListObjectsV2Paginator(s3Client, input)
		// fmt.Println("paginator", paginator)
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(context.TODO())
			if err != nil {
				return nil, err
			}
			
				// fmt.Println("commonPrefixes", page.CommonPrefixes)
			
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
					queued[*commonPrefix.Prefix] = struct{}{}
				}
			}
		fmt.Println("commonPrefixes", page.CommonPrefixes)
		}
	}
	return keys, nil
}
