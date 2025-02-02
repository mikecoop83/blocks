#!/bin/bash

if [ -n "$VERBOSE" ]; then
    set -x
fi

# Check if original logo.png exists
if [ ! -f "logo.png" ]; then
    echo "Error: logo.png not found in current directory"
    exit 1
fi

# Determine which ImageMagick command to use
if [ -x "$(command -v magick)" ]; then
    CMD="magick"
elif [ -x "$(command -v convert)" ]; then
    CMD="convert"
else
    echo "Error: ImageMagick not found (neither 'magick' nor 'convert' commands available)"
    exit 1
fi

rm -f logo-*.png

# Generate all required sizes
$CMD logo.png -resize 152x152 logo-152.png
$CMD logo.png -resize 167x167 logo-167.png
$CMD logo.png -resize 180x180 logo-180.png
$CMD logo.png -resize 192x192 logo-192.png
$CMD logo.png -resize 512x512 logo-512.png
