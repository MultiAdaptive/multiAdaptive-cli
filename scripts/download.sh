#!/bin/bash

# GitHub repository and program name
GITHUB_USER="MultiAdaptive"
REPO_NAME="multiAdaptive-cli"
PROGRAM_NAME="multiAdaptive-cli"

# Get system and architecture information
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map system and architecture to Go architecture names
case $ARCH in
  x86_64)
    ARCH="amd64"
    ;;
  aarch64)
    ARCH="arm64"
    ;;
  armv7l)
    ARCH="arm"
    ;;
  i686)
    ARCH="386"
    ;;
  ppc64)
    ARCH="ppc64"
    ;;
  ppc64le)
    ARCH="ppc64le"
    ;;
  mips)
    ARCH="mips"
    ;;
  mipsle)
    ARCH="mipsle"
    ;;
  mips64)
    ARCH="mips64"
    ;;
  mips64le)
    ARCH="mips64le"
    ;;
  s390x)
    ARCH="s390x"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Special handling for Windows
if [[ "$OS" == "mingw"* || "$OS" == "cygwin"* ]]; then
  OS="windows"
  ARCH="${ARCH}.exe"
fi

# Construct the srs file download URL
SRS_DOWNLOAD_URL="https://github.com/${GITHUB_USER}/${REPO_NAME}/releases/latest/download/srs"

# Download the srs file
echo "Downloading srs file from $SRS_DOWNLOAD_URL..."
curl -L -o "srs" "$SRS_DOWNLOAD_URL"

# Construct the download URL
DOWNLOAD_URL="https://github.com/${GITHUB_USER}/${REPO_NAME}/releases/latest/download/${PROGRAM_NAME}-${OS}-${ARCH}"

# Download the file
echo "Downloading ${PROGRAM_NAME} for ${OS}/${ARCH} from ${DOWNLOAD_URL}..."
curl -L -o ${PROGRAM_NAME} ${DOWNLOAD_URL}

# Set execute permission
chmod +x ${PROGRAM_NAME}

echo "Download complete! You can now run ${PROGRAM_NAME}."
