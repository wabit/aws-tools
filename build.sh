#!/bin/bash

# Build for macOS (64-bit)
echo "Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -o aws-tools-macos
echo "Done!"

# Build for Linux (64-bit)
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o aws-tools-linux
echo "Done!"

# Build for Windows (64-bit)
echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o aws-tools.exe
echo "Done!"

echo "All builds completed!"
