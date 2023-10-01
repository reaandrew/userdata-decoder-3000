
# Cloud-Init Userdata Decoder

## Description

Cloud-Init Userdata Decoder is a small Golang application designed to decode `cloud-init` userdata into human-readable files. These files are saved relative to your working directory, facilitating easy investigation for Cybersecurity professionals, DevOps, and Systems Administrators.

## Features

- Decode `cloud-init` userdata into readable files
- Save decoded files relative to the working directory for easy access
- Cross-platform support: Windows, Linux, macOS
- Support for various `cloud-init` data formats (e.g., YAML, JSON)

## Requirements

- Go runtime for running from source, or
- Precompiled binaries for Windows, Linux, and macOS

## Installation

### From Source

1. Clone the repository

```
git clone https://github.com/your-repo/cloud-init-decoder.git
```

2. Navigate to the project directory

```
cd cloud-init-decoder
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
./cloud-init-decoder [options]
```

### Options

- `-i, --input` : Provide the `cloud-init` userdata file for decoding.
- `-o, --output` : (Optional) Specify the output directory within your working directory.

Example:

```
./cloud-init-decoder -i user-data.txt -o output_folder/
```

## Contribution

Feel free to open issues or PRs if you find any problems or have suggestions for improvements.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

You can directly copy-paste this into your README file. Make sure to replace placeholders like repository URLs and release URLs with the actual ones.