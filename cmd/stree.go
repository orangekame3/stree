package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/orangekame3/stree/pkg"
	"github.com/spf13/cobra"

	"github.com/ddddddO/gtree"
)

var (
	awsProfile  string
	awsRegion   string
	endpointURL string
	local       bool
	noColor     bool
)

var streeCmd = &cobra.Command{
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

func init() {
	rootCmd.AddCommand(streeCmd)
	streeCmd.Flags().StringVarP(&awsProfile, "profile", "p", "local", "AWS profile to use")
	streeCmd.Flags().StringVarP(&awsRegion, "region", "r", "us-east-1", "AWS region to use (overrides the region specified in the profile)")
	streeCmd.Flags().StringVarP(&endpointURL, "endpoint-url", "e", "http://localhost:4566", "AWS endpoint URL to use (useful for local testing with LocalStack)")
	streeCmd.Flags().BoolVar(&local, "local", false, "Use LocalStack configuration")
	streeCmd.Flags().BoolVar(&noColor, "no-color", false, "Disable colorized output")
}
