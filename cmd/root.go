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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stree",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if noColor {
			color.NoColor = true // disables colorized output
		}

		s3Config := pkg.S3Config{
			AwsProfile:  awsProfile,
			AwsRegion:   awsRegion,
			EndpointURL: endpointURL,
			Local:       local,
		}

		s3Svc := pkg.InitializeAWSSession(s3Config)

		bucketAndPrefix := strings.SplitN(args[0], "/", 2)
		bucket := bucketAndPrefix[0]
		prefix := ""
		if len(bucketAndPrefix) > 1 {
			prefix = bucketAndPrefix[1]
		}

		root := gtree.NewRoot(color.BlueString(bucket))

		keys, err := pkg.FetchS3ObjectKeys(s3Svc, bucket, prefix)
		if err != nil {
			fmt.Println("failed to fetch S3 object keys:", err)
			return
		}

		root = pkg.BuildTree(root, keys, noColor)

		if err := gtree.OutputProgrammably(os.Stdout, root); err != nil {
			fmt.Println(err)
			return
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
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
	rootCmd.Flags().BoolVar(&local, "local", false, "Use LocalStack configuration")
	rootCmd.Flags().BoolVar(&noColor, "no-color", false, "Disable colorized output")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
