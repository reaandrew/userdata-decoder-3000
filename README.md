
# Userdata Decoder 3000

![logo.png](html_client%2Fsrc%2Fimages%2Flogo.png)

## Description

User Data Decoder 3000 (name inspiration from Joe Sparkes) decodes cloud instance startup data. It works with AWS and handles both cloud-init and plain text formats. Useful for developers, sysadmins, and cybersecurity roles. Azure and GCP support are on the roadmap.

This tool was developed as a way to check no sensitive data is included in the userdata once it is created and attached to either an EC2 instance or AWS Launch Template.

There is a client side version available too just to serve as a demonstration and quick view of some user data.

## CLI

### Installation

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

Given the following User Data (example and partial content) which is base64 encoded gzipped multi-part mime message with embedded encoded files inside, you can see before and after running the tool where the content is decoded and rebuilt into a directory structure consistent with the data:

```shell
50PcHBExXr4k2LwUZUpAyiY2hKtPFzKLS179+8galc2Ad8WTPqPGS4vAtBvbHcSsWMLqpPcJ6/Ln
sg+VmrXNj90ETOYYhusFekUp7EHs/hbF11ftgljwYtBpVR21lwMU/6FNd8kxU6Mgjs5RfmPt+/EU
xp3aOU+6WNYKPt/0VjbV4GNJ5/iS4OKQ/8M+KmZNpSXfle9Zo6zKQf4m4wTCC6TDRpNHFbVboW/N
UWCPsu+mdPBO58XWZL6F0wuaUAwbSLxDQF3wADN3SJDVocOqHPaYnd+NoDd.... AND THE REST
```

**OUTPUT**

```shell
-- i-0c9086lee0g83cg3d
    -- etc
       |-- chrony.conf
       |-- containerd
       |   -- config.toml
       |-- eks
       |   -- containerd
       |      -- pull-content.sh
       |-- environment
       |-- kubernetes
       |   |-- kubelet
       |       |-- kubeconfig
       |        -- kubelete-config.json
       AND THE REST ...
```

Even if your userdata is simple content or not encoded, it will still be pulled down and output into a folder matching the ID of the instance

```shell
-- i-0c9086lee0g83cg3d
   -- userdata
````

In this example, the content of the userdata is placed in a file called user data.

## Next Steps

After you have decoded all the user data for the instances, you can then run analysis over the files including secret detection.


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

## Installation from source

### From Source

1. Clone the repository

```
git clone git@github.com:reaandrew/userdata-decoder-3000.git
```

2. Navigate to the project directory

```
cd userdata-decoder-3000
```

3. Build the project

```
go build
```

## Website

[https://userdata-decoder-3000.andrewrea.co.uk/](https://userdata-decoder-3000.andrewrea.co.uk/)

![website_image.png](html_client%2Fsrc%2Fimages%2Fwebsite_image.png)


## Contribution

Feel free to open issues or PRs if you find any problems or have suggestions for improvements.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

You can directly copy-paste this into your README file. Make sure to replace placeholders like repository URLs and release URLs with the actual ones.