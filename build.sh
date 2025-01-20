#!/bin/bash
set -e

mkdir -p out/
rm -f out/*
env GOOS=js GOARCH=wasm go build -o out/blocks.wasm github.com/mikecoop83/blocks

