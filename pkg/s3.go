// Package pkg provides the core functionality of the program.
package pkg

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

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

// ref: https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
func formatBytes(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%c", float64(b)/float64(div), "KMGTPE"[exp])
}

// FetchS3ObjectKeys returns a slice of keys for all objects in the specified bucket and prefix
func FetchS3ObjectKeys(s3Client *s3.Client, bucket string, prefix string, maxDepth *int, size, humanReadable bool, dateTime bool, username bool, pattern string, inversePattern string) ([][]string, error) {
	var delimiter *string
	var fetchOwner *bool
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

		if username {
			fetchOwner = aws.Bool(true)
		}

		input := &s3.ListObjectsV2Input{
			Bucket:     aws.String(bucket),
			Prefix:     aws.String(currentPrefix),
			Delimiter:  delimiter,
			FetchOwner: fetchOwner,
		}

		paginator := s3.NewListObjectsV2Paginator(s3Client, input)

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(context.TODO())
			if err != nil {
				return nil, err
			}

			for _, obj := range page.Contents {
				key := strings.Split(*obj.Key, "/")

				meta := []string{}

				if username {
					meta = append(meta, *obj.Owner.DisplayName)
				}

				if humanReadable {
					meta = append(meta, formatBytes(*obj.Size))
				} else if size {
					meta = append(meta, fmt.Sprintf("%d", *obj.Size))
				}
				if dateTime {
					t := *obj.LastModified
					layout := "Jan 2 15:04"
					formatted := t.In(time.Local).Format(layout)
					fmt.Println(formatted)
					meta = append(meta, formatted)
				}
				if len(meta) > 0 {
					key[len(key)-1] = fmt.Sprintf("[%7s] %s", strings.Join(meta, " "), key[len(key)-1])
				}
				include := true
				if pattern != "" {
					include, _ = filepath.Match(pattern, filepath.Base(*obj.Key))
				}
				if inversePattern != "" {
					exclude, _ := filepath.Match(inversePattern, filepath.Base(*obj.Key))
					if exclude {
						include = false
					}
				}
				if include {
					keys = append(keys, key)
				}

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
		}
	}
	return keys, nil
}
