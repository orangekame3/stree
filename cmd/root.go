//Package cmd is a command line tool for visualizing the structure of S3 buckets
/*
Copyright © 2023 Takafumi Miyanaga <miya.org.0309@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"log"
	"os"

	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/fatih/color"
	"github.com/orangekame3/stree/pkg"

	"github.com/ddddddO/gtree"
)

var (
	awsProfile  string
	awsRegion   string
	endpointURL string
	local       bool
	noColor     bool
)

var rootCmd = &cobra.Command{
	Use:   "stree [bucket/prefix]",
	Short: "stree is a command line tool for visualizing the structure of S3 buckets",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		s3Config := pkg.S3Config{
			AwsProfile:  awsProfile,
			AwsRegion:   awsRegion,
			EndpointURL: endpointURL,
			Local:       local,
		}

		s3Svc := pkg.InitializeAWSSession(s3Config)

		bucket, prefix, err := extractBucketAndPrefix(args[0])
		if err != nil {
			log.Fatalf("failed to extract bucket and prefix: %v", err)
		}

		keys, err := pkg.FetchS3ObjectKeys(s3Svc, bucket, prefix)
		if err != nil {
			log.Fatalf("failed to fetch S3 object keys: %v", err)
			return
		}

		root := gtree.NewRoot(color.BlueString(bucket))
		if noColor {
			root = pkg.BuildTreeWithoutColor(root, keys)
		} else {
			root = pkg.BuildTreeWithColor(root, keys)
		}

		if err := gtree.OutputProgrammably(os.Stdout, root); err != nil {
			log.Fatalf("failed to output tree: %v", err)
			return
		}
		fileCount, dirCount := pkg.ProcessKeys(keys)
		fmt.Printf("\n%d directories, %d files\n", dirCount, fileCount)
	},
}

// Execute executes the root command.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&awsProfile, "profile", "p", "local", "AWS profile to use")
	rootCmd.Flags().StringVarP(&awsRegion, "region", "r", "us-east-1", "AWS region to use (overrides the region specified in the profile)")
	rootCmd.Flags().StringVarP(&endpointURL, "endpoint-url", "e", "http://localhost:4566", "AWS endpoint URL to use (useful for local testing with LocalStack)")
	rootCmd.Flags().BoolVarP(&local, "local", "l", false, "Use LocalStack configuration")
	rootCmd.Flags().BoolVarP(&noColor, "no-color", "n", false, "Disable colorized output")
}

func extractBucketAndPrefix(input string) (string, string, error) {
	if input == "" {
		return "", "", errors.New("[bucket/prefix] cannot be empty")
	}

	parts := strings.SplitN(input, "/", 2)
	bucket := parts[0]

	if len(parts) == 1 {
		return bucket, "", nil
	}

	return bucket, parts[1], nil
}

func SetVersionInfo(version, date string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s)", version, date)
}
