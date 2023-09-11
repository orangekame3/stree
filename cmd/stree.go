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
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var awsProfile string
var awsRegion string

// streeCmd represents the stree command
var streeCmd = &cobra.Command{
	Use:   "stree",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sessOptions := session.Options{
			Profile: awsProfile,
		}

		if awsRegion != "" {
			sessOptions.Config = aws.Config{Region: aws.String(awsRegion)}
		}

		sess, err := session.NewSessionWithOptions(sessOptions)
		if err != nil {
			log.Fatalf("Failed to create session: %v", err)
		}

		svc := s3.New(sess)

		bucket := args[0]

		input := &s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
		}
		result, err := svc.ListObjectsV2(input)
		if err != nil {
			log.Fatalf("Unable to list bucket %q objects: %v", bucket, err)
		}

		for _, item := range result.Contents {
			fmt.Printf("%s\n", aws.StringValue(item.Key))
		}
	},
}

func init() {
	rootCmd.AddCommand(streeCmd)
	streeCmd.Flags().StringVarP(&awsProfile, "profile", "p", "default", "AWS profile to use")
	streeCmd.Flags().StringVarP(&awsRegion, "region", "r", "", "AWS region to use (overrides the region specified in the profile)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// streeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// streeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
