<br>

<p align="center">
<img src="img/stree.png" alt="S3 directory tree visualization" height="350" width="350"/>
</p>

<p align="center">
<a href="https://opensource.org/licenses/MIT">
<img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="MIT License badge">
</a>
<a href="https://pkg.go.dev/github.com/orangekame3/stree">
<img src="https://github.com/orangekame3/stree/actions/workflows/release.yml/badge.svg" alt="Release workflow status badge">
</a>
<a href="https://github.com/orangekame3/stree/actions/workflows/tagpr.yml">
<img src="https://github.com/orangekame3/stree/actions/workflows/tagpr.yml/badge.svg" alt="Tag PR workflow status badge">
</a>
</p>

<p align="center">
Directory trees of S3
</p>

<p align="center">
<img src="img/demo.gif" alt="Demonstration of stree tool in action" height="auto" width="auto"/>
</p>

# Overview

`stree` is a CLI tool designed to visualize the directory tree structure of an S3 bucket.  
By inputting an S3 bucket/prefix and utilizing various flags to customize your request, you can obtain a colorized or non-colorized directory tree right in your terminal.

Whether it's for verifying the file structure, sharing the structure with your team, or any other purpose, `stree` offers an easy and convenient way to explore your S3 buckets.

# Features

- **Colorized Output**: By default, `stree` provides a colorized tree structure, making it easy to differentiate between directories and files at a glance. This feature can be turned off with the `-n` or `--no-color` flag.
- **LocalStack Support**: `stree` supports local testing with LocalStack, a fully functional local AWS cloud stack, thanks to the `--local` and `--endpoint-url` flags.
- **Custom AWS Profile and Region**: Specify the AWS profile and region with the `--profile` and `--region` flags to override the default settings as needed.
- **Ease of Installation**: Install `stree` via Go, Homebrew, or by downloading the latest compiled binaries from the GitHub releases page.

# Install

## Go

```shell
go install github.com/orangekame3/stree@latest
```

## Homebrew

```shell
brew install orangekame3/tap/stree
```

## Download

Download the latest compiled binaries and put it anywhere in your executable path.

[Download here](https://github.com/orangekame3/stree/releases)

# Getting Started

## Prerequisites

Before using `stree`, ensure that you have Go installed on your machine, or Homebrew for macOS users. You would also need to configure your AWS credentials appropriately to access your S3 buckets.

## Running the Tool

After installing `stree`, run it with the necessary bucket/prefix and flags as shown in the usage section above. Here are a few examples to get you started:

### Display the directory tree using a specific AWS profile and region

```shell
$ stree my-bucket -p my-profile -r us-east-1
my-bucket
└── test
    ├── dir1
    │   ├── dir1_1
    │   │   └── dir1_1_1
    │   │       ├── file1.csv
    │   │       └── file2.csv
    │   └── dir1_2
    │       ├── file1.csv
    │       ├── file2.csv
    │       └── file3.csv
    ├── dir2
    │   └── dir2_1
    │       └── dir2_1_1
    │           ├── file1.csv
    │           ├── file2.csv
    │           └── file3.csv
    └── dir3
        ├── file1.csv
        └── file2.csv

9 directories, 10 files
```

### Display the sub-directory tree using a specific AWS profile and region

```shell
$ stree my-bucket/test/dir2 -p my-profile -r us-east-1
my-bucket
└── test
    └── dir2
        └── dir2_1
            └── dir2_1_1
                ├── file1.csv
                ├── file2.csv
                └── file3.csv

4 directories, 3 files
```

### Display the directory tree using Localstack

```shell
$ stree my-bucket/test/dir2 -l
my-bucket
└── test
    └── dir2
        └── dir2_1
            └── dir2_1_1
                ├── file1.csv
                ├── file2.csv
                └── file3.csv

4 directories, 3 files
```

# Usage

```shell
Usage:
  stree [bucket/prefix] [flags]

Flags:
  -e, --endpoint-url string   AWS endpoint URL to use (useful for local testing with LocalStack) (default "http://localhost:4566")
  -h, --help                  help for stree
  -l, --local                 Use LocalStack configuration
  -n, --no-color              Disable colorized output
  -p, --profile string        AWS profile to use (default "local")
  -r, --region string         AWS region to use (overrides the region specified in the profile) (default "us-east-1")
```

# License

`stree` is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

# Acknowledgments

The concept of "stree" was inspired by the pioneering work seen in [gtree](https://github.com/ddddddO/gtree). I'm grateful for the inspiration.
