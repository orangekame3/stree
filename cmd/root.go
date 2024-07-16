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
	mfa         bool
	level       int
	fullPath    bool
	fileName string

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
			MFA:         mfa,
		}

		s3Svc,err := pkg.InitializeAWSSession(s3Config)
		if err != nil {
			log.Fatalf("failed to initialize AWS session: %v", err)
			return
		}

		bucket, prefix, err := extractBucketAndPrefix(args[0])
		if err != nil {
			log.Fatalf("failed to extract bucket and prefix: %v", err)
		}
		var maxDepth *int
		if level > 0 {
			maxDepth = &level
		}
		keys, err := pkg.FetchS3ObjectKeys(s3Svc, bucket, prefix, maxDepth)
		if err != nil {
			log.Fatalf("failed to fetch S3 object keys: %v", err)
			return
		}

		root := gtree.NewRoot(bucket)
		if noColor || fileName != "" {
			root = pkg.BuildTreeWithoutColor(root, bucket, keys, fullPath)
		} else {
			root = gtree.NewRoot(color.BlueString(bucket))
			root = pkg.BuildTreeWithColor(root, bucket, keys, fullPath)
		}


		fileCount, dirCount := pkg.ProcessKeys(keys)
		
		if fileName != "" {
			f, err := os.Create(fileName)
			if err != nil {
				log.Fatalf("failed to create file: %v", err)
				return
			}
			defer f.Close()
			if err := gtree.OutputProgrammably(f, root); err != nil {
				log.Fatalf("failed to output tree: %v", err)
				return
			}
			fmt.Fprintf(f, "\n%d directories, %d files\n", dirCount, fileCount)
			
		}else{
			if err := gtree.OutputProgrammably(os.Stdout, root); err != nil {
				log.Fatalf("failed to output tree: %v", err)
				return
			}
			fmt.Printf("\n%d directories, %d files\n", dirCount, fileCount)
		}
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
	rootCmd.Flags().StringVarP(&awsProfile, "profile", "p", defaultProfile(), "AWS profile to use")
	rootCmd.Flags().StringVarP(&awsRegion, "region", "r", defaultRegion(), "AWS region to use (overrides the region specified in the profile)")
	rootCmd.Flags().StringVarP(&endpointURL, "endpoint-url", "e", "", "AWS endpoint URL to use (useful for local testing with LocalStack)")
	rootCmd.Flags().BoolVarP(&local, "local", "l", false, "Use LocalStack configuration")
	rootCmd.Flags().BoolVarP(&noColor, "no-color", "n", false, "Disable colorized output")
	rootCmd.Flags().BoolVarP(&mfa, "mfa", "m", false, "Use Multi-Factor Authentication")
	rootCmd.Flags().IntVarP(&level, "level", "L", 0, "Descend only level directories")
	rootCmd.Flags().BoolVarP(&fullPath, "full-path", "f", false, "Print the full path prefix for each file.")
	rootCmd.Flags().StringVarP(&fileName, "output", "o", "", "Send output to filename.")
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

func defaultProfile() string {
	if p, ok := os.LookupEnv("AWS_PROFILE"); ok {
		return p
	}
	return "default"
}

func defaultRegion() string {
	for _, e := range []string{"AWS_REGION", "AWS_DEFAULT_REGION"} {
		if r, ok := os.LookupEnv(e); ok {
			return r
		}
	}
	return ""
}
