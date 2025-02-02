#!/bin/bash

set -e

if [ -n "$VERBOSE" ]; then
  set -x
fi

# Install Go if not already installed
if ! [ -x "$(command -v go)" ]; then
  echo "Installing Go for Linux..."
  curl -fsSL https://go.dev/dl/go1.23.4.linux-amd64.tar.gz -o go.tar.gz
  tar -xzf go.tar.gz
  PATH=$PATH:$(pwd)/go/bin
  export PATH
fi

# Generate icons if ImageMagick is available
if [ -x "$(command -v magick)" ] || [ -x "$(command -v convert)" ]; then
  cd resources
  ./generate-icons.sh
  cd ..
else
  echo "Warning: ImageMagick not found, skipping icon generation"
fi

mkdir -p out/
rm -f out/*
env GOOS=js GOARCH=wasm go build -o out/blocks.wasm github.com/mikecoop83/blocks
cp resources/wasm_exec.js out/
cp resources/*.html out/
cp resources/*.png out/
cp resources/manifest.json out/
cp resources/sw.js out/
