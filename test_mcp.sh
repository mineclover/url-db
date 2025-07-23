#!/bin/bash

# Test script for MCP functionality
# This script tests the basic MCP server functionality

echo "Testing MCP Server Implementation..."
echo "====================================="

# Build the server
echo "1. Building server..."
go build -o bin/url-db ./cmd/server
if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi
echo "✅ Build successful"

# Test initialization
echo "2. Testing initialization..."
result=$(echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null)
if echo "$result" | grep -q '"protocolVersion":"2024-11-05"'; then
    echo "✅ Initialization successful"
else
    echo "❌ Initialization failed"
    echo "Result: $result"
    exit 1
fi

# Test tools list
echo "3. Testing tools list..."
result=$(echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null)
if echo "$result" | grep -q '"name":"get_server_info"'; then
    echo "✅ Tools list successful"
else
    echo "❌ Tools list failed"
    echo "Result: $result"
    exit 1
fi

# Test tool call
echo "4. Testing tool call..."
result=$(echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"get_server_info","arguments":{}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null)
if echo "$result" | grep -q "url-db-mcp-server"; then
    echo "✅ Tool call successful"
else
    echo "❌ Tool call failed"
    echo "Result: $result"
    exit 1
fi

# Test list domains
echo "5. Testing list_domains tool..."
result=$(echo '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"list_domains","arguments":{}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null)
if echo "$result" | grep -q '"type":"text"'; then
    echo "✅ List domains successful"
else
    echo "❌ List domains failed"
    echo "Result: $result"
    exit 1
fi

# Test create domain (should fail if exists)  
echo "6. Testing create_domain tool..."
result=$(echo '{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"create_domain","arguments":{"name":"test-mcp-domain","description":"Test domain for MCP"}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null)
if echo "$result" | grep -q -E '(Successfully created|already exists)'; then
    echo "✅ Create domain working"
else
    echo "❌ Create domain failed"
    echo "Result: $result"
    exit 1
fi

# Test list nodes
echo "7. Testing list_nodes tool..."
result=$(echo '{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"list_nodes","arguments":{"domain_name":"tech"}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null)
if echo "$result" | grep -q -E '(Node ID|No nodes found)'; then
    echo "✅ List nodes working"
else
    echo "❌ List nodes failed"
    echo "Result: $result"
    exit 1
fi

# Test invalid mode
echo "8. Testing invalid mode validation..."
result=$(./bin/url-db -mcp-mode=invalid 2>&1)
if echo "$result" | grep -q "Invalid MCP mode"; then
    echo "✅ Invalid mode validation working"
else
    echo "❌ Invalid mode validation failed"
    echo "Result: $result"
    exit 1
fi

echo ""
echo "🎉 All MCP tests passed!"
echo "✅ MCP server implementation is working correctly"
echo ""
echo "Usage examples:"
echo "- HTTP mode: ./bin/url-db"
echo "- MCP stdio mode: ./bin/url-db -mcp-mode=stdio"
echo "- MCP SSE mode: ./bin/url-db -mcp-mode=sse (not implemented yet)"
echo "- MCP HTTP mode: ./bin/url-db -mcp-mode=http (not implemented yet)"