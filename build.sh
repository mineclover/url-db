#!/bin/bash

echo "Building URL-DB Server..."

# Set Go environment
export GO111MODULE=on

# Clean previous builds
rm -rf bin
mkdir -p bin

# Build the application
echo "Building for current platform..."
go build -o bin/url-db cmd/server/main.go
if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

echo "Build completed successfully!"
echo "Executable created: bin/url-db"

# Run tests
echo "Running tests..."
go test -v ./...
if [ $? -ne 0 ]; then
    echo "Tests failed!"
    exit 1
fi

echo "All tests passed!"
echo
echo "To run the server:"
echo "  ./bin/url-db"
echo
echo "Default configuration:"
echo "  Port: 8080"
echo "  Database: file:./url-db.sqlite"
echo "  Tool Name: url-db"