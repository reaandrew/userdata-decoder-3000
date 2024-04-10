
# Userdata Decoder 3000
See what things are being stored in user data to ensure there is no sensitive information in there!

![logo.png](html_client%2Fsrc%2Fimages%2Flogo.png)


## Getting Started

### Linux

```shell
curl -L -o udd https://github.com/reaandrew/userdata-decoder-3000/releases/download/v1.4.1/udd-linux-amd64
./udd --provider aws --output output-path
```

### Mac

```shell
curl -L -o udd https://github.com/reaandrew/userdata-decoder-3000/releases/download/v1.4.1/udd-darwin-amd64
./udd --provider aws --output output-path
```

### Windows

```shell
Invoke-WebRequest -Uri "https://github.com/reaandrew/userdata-decoder-3000/releases/download/v1.4.1/udd-windows-amd64.exe" -OutFile "udd.exe"
./udd.exe --provider aws --output output-path
```

## Description

User Data Decoder 3000 (name inspiration from Joe Sparkes) decodes cloud instance startup data. It works with AWS and handles both cloud-init and plain text formats. Useful for developers, sysadmins, and cybersecurity roles. Azure and GCP support are on the roadmap.

This tool was developed primarily as a demonstration and also as a way to check no sensitive data is included in the userdata once it is created and attached to either an EC2 instance or AWS Launch Template.

> Although you can only access instance metadata and user data from within the instance itself, the data is not protected by authentication or cryptographic methods. Anyone who has direct access to the instance, and potentially any software running on the instance, can view its metadata. Therefore, you should not store sensitive data, such as passwords or long-lived encryption keys, as user data.
>
> https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html

There is a client side version available too just to serve as a demonstration and quick view of some user data.

## Some info about user data

EC2 User Data is utilised to pass startup scripts or other configuration data into instances at launch. The flexibility of User Data allows for various formats, depending on the instance's needs and the tasks you want to automate. Below are some of the formats that EC2 User Data can take:

1. **Plain Text**

Plain text can be used for simple commands or instructions that don't require encoding. It's directly interpreted by the shell if the first line is a shebang (#!), indicating what interpreter to use (e.g., #!/bin/bash for a shell script).

2. **Shell Scripts**

Shell scripts are perhaps the most common use case, allowing you to execute a series of commands automatically upon instance startup. 

3. **Cloud-Init Directives**

Cloud-init is a widely used method for early initialization of cloud instances. Supported by many cloud providers, cloud-init scripts can be used to perform tasks such as installing packages, writing files, and configuring users or security settings. 

4. **Multi-part MIME Messages**

When you need to pass multiple pieces of information or scripts of different types, you can use multi-part MIME messages. This format allows you to combine shell scripts, cloud-init directives, and other data types into a single User Data payload. Each part of the message can be of a different MIME type, enabling complex initialization sequences.

5. **Gzipped Content**

User Data supports gzipped content, which is useful for compressing large initialization scripts or data. Gzipped content must be decompressed by the receiving script or application. This is particularly useful for optimizing the use of the 16 KB limit on User Data.

6. **Base64 Encoded Data**

Any User Data must be base64 encoded. 

Read more at [https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instancedata-add-user-data.html](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instancedata-add-user-data.html)

## CLI

### Installation

Download the precompiled binaries for your respective operating system from the Releases section.


- [Linux](https://github.com/reaandrew/userdata-decoder-3000/releases/download/v1.4.1/udd-linux-amd64)
- [macOS](https://github.com/reaandrew/userdata-decoder-3000/releases/download/v1.4.1/udd-darwin-amd64)
- [Windows](https://github.com/reaandrew/userdata-decoder-3000/releases/download/v1.4.1/udd-windows-amd64.exe)

## Usage

1. Navigate to your working directory where you want the decoded files to be saved.

2. Run the application

```
Usages:
  udd [OPTIONS]             : Specify options for data provider and output directory.
  udd [CONTENT_TO_DECODE]   : Specify the content to decode as the first argument.

Options:
  -o, --output-dir: Specify the output directory within your working directory. (default "output")
  -p, --provider:   Specify the data provider (e.g., aws).

```

### Options

- `-o, --output-dir` : (Optional) Specify the output directory within your working directory.
- `-p, --provider` : (Optional) Specify a cloud provider to get startup data from instances

Example:

```
./udd $(cat test_data/gzip_base64_userdata.txt) 

OR

./udd --provider aws
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