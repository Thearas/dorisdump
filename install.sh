#!/bin/bash

set -e

os=$(uname -s)
if [ -n "$OS" ]; then
    os="$OS"
fi

arch=$(uname -m)
if [ -n "$ARCH" ]; then
    arch="$ARCH"
fi

os=$(echo "$os" | tr  "[:upper:]" "[:lower:]")
arch=$(echo "$arch" | tr  "[:upper:]" "[:lower:]")

if [[ "$arch" == "aarch64" ]]; then
    arch="arm64"
fi

binary_name="dorisdump-${os}-${arch}"

# get latest
latest_download_url=$(curl -s https://api.github.com/repos/Thearas/dorisdump/releases/latest | grep "browser_download_url.*${binary_name}" | cut -d : -f 2,3 | tr -d \")
if [[ -z "$latest_download_url" ]]; then
    echo "No release found for ${os}-${arch}"
    exit 1
fi

name=$(basename "$latest_download_url")

wget "$latest_download_url" -O "/tmp/$name"
tar -xzf "/tmp/$name" -C /usr/local/bin/ && mv "/usr/local/bin/$binary_name" /usr/local/bin/dorisdump
rm -f "/tmp/$name"

echo "dorisdump installed successfully"

/usr/local/bin/dorisdump completion print-help
