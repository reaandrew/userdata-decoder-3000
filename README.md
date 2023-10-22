
# Cloud Startup Data Decoder

## Description

Cloud Startup Data Decoder decodes cloud instance startup data. It works with AWS and handles both cloud-init and plain text formats. Useful for developers, sysadmins, and cybersecurity roles. Azure and GCP support are on the roadmap.

## Features

### Current Features

- AWS User Data Support: Decode and analyze user data in Amazon EC2 instances.
- Cloud-Init Support: Decode cloud-init formatted user data, widely used for cloud instance initialization.
- Plain User Data: Decode plain-text user data that doesn't follow the cloud-init format.

### Upcoming Features (TODO)

- Azure Support: Extend decoding capabilities to Azure's Custom Script Extension and user data.
- GCP Support: Extend decoding capabilities to include Google Cloud Platform's startup scripts and metadata.

### Usage Scenarios

- Cloud Migration: Easily understand startup configurations when migrating across cloud platforms.
- Security Auditing: Utilize this tool to audit startup scripts for any security vulnerabilities.
- Debugging: Decode and analyze startup data to debug instance initialization issues.

### How it works

- Decode userdata into readable files in the output directory
- Save decoded files relative to the working directory for easy access
- Cross-platform support: Windows, Linux, macOS

## Requirements

- Go runtime for running from source, or
- Precompiled binaries for Windows, Linux, and macOS

## Installation

### From Source

1. Clone the repository

```
git clone git@github.com:reaandrew/cloud-startup-data-decoder.git
```

2. Navigate to the project directory

```
cd cloud-startup-data-decoder
```

3. Build the project

```
go build
```

### Precompiled Binaries

Download the precompiled binaries for your respective operating system from the Releases section.

- [Windows](url-to-windows-release)
- [Linux](url-to-linux-release)
- [macOS](url-to-macos-release)

## Usage

1. Navigate to your working directory where you want the decoded files to be saved.

2. Run the application

```
Usages:
  cloud-startup-data-decoder [OPTIONS]             : Specify options for data provider and output directory.
  cloud-startup-data-decoder [CONTENT_TO_DECODE]   : Specify the content to decode as the first argument.

Options:
  -o, --output-dir: Specify the output directory within your working directory. (default "output")
  -p, --provider:   Specify the data provider (e.g., aws).

```

### Options

- `-o, --output-dir` : (Optional) Specify the output directory within your working directory.
- `-p, --provider` : (Optional) Specify a cloud provider to get startup data from instances

Example:

```
./cloud-startup-data-decoder $(cat test_data/gzip_base64_userdata.txt) 

OR

./cloud-startup-data-decoder --provider aws
```

## Contribution

Feel free to open issues or PRs if you find any problems or have suggestions for improvements.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

You can directly copy-paste this into your README file. Make sure to replace placeholders like repository URLs and release URLs with the actual ones.