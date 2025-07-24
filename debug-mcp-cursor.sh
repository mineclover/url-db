#!/bin/bash

echo "=== MCP Server Debug Script ==="

# 1. MCP 서버 직접 테스트
echo "1. Testing MCP server directly..."
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "cursor", "version": "1.3.0"}}}' | ./bin/url-db -mcp-mode=stdio -db-path=./url-db.sqlite

echo -e "\n2. Testing tools/list..."
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}' | ./bin/url-db -mcp-mode=stdio -db-path=./url-db.sqlite

echo -e "\n3. Testing resources/list..."
echo '{"jsonrpc": "2.0", "id": 3, "method": "resources/list", "params": {}}' | ./bin/url-db -mcp-mode=stdio -db-path=./url-db.sqlite

echo -e "\n=== Cursor MCP Configuration ==="
echo "Global MCP config:"
cat ~/.cursor/mcp.json

echo -e "\nLocal MCP config:"
cat .cursor/mcp.json

echo -e "\n=== File Permissions ==="
ls -la bin/url-db
ls -la url-db.sqlite

echo -e "\n=== Cursor Logs ==="
echo "Latest Cursor logs:"
find ~/Library/Application\ Support/Cursor/logs/ -name "*.log" -type f -exec ls -la {} \; | tail -5

echo -e "\n=== Process Check ==="
ps aux | grep -E "(cursor|url-db)" | grep -v grep

echo -e "\n=== Network Check ==="
lsof -i :8080 2>/dev/null || echo "No process on port 8080"
lsof -i :8082 2>/dev/null || echo "No process on port 8082"

echo -e "\n=== Debug Complete ===" 