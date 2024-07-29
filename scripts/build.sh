#!/bin/bash

# Program name
PROGRAM_NAME="multiAdaptive-cli"

# List of supported OS and architecture combinations
OS_ARCH_LIST=(
  "linux/amd64"
  "linux/arm64"
  "linux/386"
  "linux/arm"
  "linux/ppc64"
  "linux/ppc64le"
  "linux/mips"
  "linux/mipsle"
  "linux/mips64"
  "linux/mips64le"
  "linux/s390x"
  "darwin/amd64"
  "darwin/arm64"
  "windows/amd64"
  "windows/386"
  "windows/arm"
  "windows/arm64"
)

# Output directory
OUTPUT_DIR="./build"

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Iterate over each OS and architecture combination and build
for os_arch in "${OS_ARCH_LIST[@]}"; do
  os="${os_arch%/*}"
  arch="${os_arch#*/}"
  
  output_name="${PROGRAM_NAME}-${os}-${arch}"
  
  if [ "$os" = "windows" ]; then
    output_name="${output_name}.exe"
  fi
  
  echo "Building for $os/$arch..."
  env GOOS="$os" GOARCH="$arch" go build -o "$OUTPUT_DIR/$output_name" .
done

echo "Build complete! All binaries are in the $OUTPUT_DIR directory."
