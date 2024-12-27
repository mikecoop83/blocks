#!/bin/bash
set -e

# Install Go if not already installed
if ! [ -x "$(command -v go)" ]; then
  echo "Installing Go for Linux..."
  curl -fsSL https://go.dev/dl/go1.23.4.linux-amd64.tar.gz -o go.tar.gz
  tar -xzf go.tar.gz
  export PATH=$PATH:$(pwd)/go/bin
fi

mkdir -p out/
rm -f out/*
env GOOS=js GOARCH=wasm go build -o out/blocks.wasm github.com/mikecoop83/blocks
cp resources/wasm_exec.js out/
cp resources/*.html out/
