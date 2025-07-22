#!/bin/bash
set -e

echo "🚀 URL-DB Python Integration Tests Runner"
echo "========================================"

# Check if server binary exists, if not build it
if [ ! -f "bin/url-db" ]; then
    echo "📦 Server binary not found. Building..."
    make build
fi

# Run Python tests
echo "🐍 Running Python integration tests..."
cd tests/python
python3 run_tests.py --category all --verbose

echo "✅ All Python tests completed!"