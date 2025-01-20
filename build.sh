#!/bin/bash
set -e

# Install Go if not already installed
if ! [ -x "$(command -v go)" ]; then
  echo "Installing Go for Linux..."
  curl -fsSL https://go.dev/dl/go1.23.4.linux-amd64.tar.gz -o go.tar.gz
  tar -xzf go.tar.gz
  PATH=$PATH:$(pwd)/go/bin
fi

# Define output directories
OUTPUT_DIR="dist"

# Create directories if they don't exist
mkdir -p OUTPUT_DIR

echo "Building WebAssembly (WASM)..."
GOOS=js GOARCH=wasm go build -o $OUTPUT_DIR/blocks.wasm

echo "Copying static files..."
cp -r static/* $OUTPUT_DIR

echo "Build completed successfully!"