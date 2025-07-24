#!/bin/bash

# MCP Server Debug Script
# This script helps debug MCP server connection issues

echo "🔍 MCP Server Debug Session"
echo "=========================="

# Check if binary exists and is executable
if [ ! -f "./bin/url-db" ]; then
    echo "❌ Binary not found: ./bin/url-db"
    echo "Run 'make build' first"
    exit 1
fi

if [ ! -x "./bin/url-db" ]; then
    echo "❌ Binary not executable: ./bin/url-db"
    chmod +x ./bin/url-db
    echo "✅ Fixed permissions"
fi

# Test basic functionality
echo "📋 Testing basic functionality..."
./bin/url-db -version
echo ""

# Test MCP initialization
echo "📡 Testing MCP initialization..."
echo '{"jsonrpc": "2.0", "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}}, "id": 1}' | ./bin/url-db -mcp-mode=stdio -db-path=/Users/junwoobang/mcp/url-db/url-db.sqlite > /tmp/mcp-init.json 2>&1

if [ $? -eq 0 ]; then
    echo "✅ MCP initialization successful"
    echo "Response:"
    cat /tmp/mcp-init.json
    echo ""
else
    echo "❌ MCP initialization failed"
    cat /tmp/mcp-init.json
    exit 1
fi

# Test tools list
echo "🛠️ Testing tools list..."
(echo '{"jsonrpc": "2.0", "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}}, "id": 1}'; echo '{"jsonrpc": "2.0", "method": "notifications/initialized"}'; echo '{"jsonrpc": "2.0", "method": "tools/list", "params": {}, "id": 2}') | ./bin/url-db -mcp-mode=stdio -db-path=/Users/junwoobang/mcp/url-db/url-db.sqlite > /tmp/mcp-tools.json 2>&1

# Count tools
TOOL_COUNT=$(grep -o '"name":"[^"]*"' /tmp/mcp-tools.json | wc -l | tr -d ' ')
echo "✅ Found $TOOL_COUNT tools"

# Show first few tools
echo "📝 First 3 tools:"
grep -o '"name":"[^"]*"' /tmp/mcp-tools.json | head -3 | sed 's/"name":"\([^"]*\)"/- \1/'

# Test a simple tool call
echo ""
echo "⚡ Testing get_server_info tool..."
(echo '{"jsonrpc": "2.0", "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}}, "id": 1}'; echo '{"jsonrpc": "2.0", "method": "notifications/initialized"}'; echo '{"jsonrpc": "2.0", "method": "tools/call", "params": {"name": "get_server_info", "arguments": {}}, "id": 3}') | ./bin/url-db -mcp-mode=stdio -db-path=/Users/junwoobang/mcp/url-db/url-db.sqlite > /tmp/mcp-call.json 2>&1

if grep -q '"text":' /tmp/mcp-call.json; then
    echo "✅ Tool call successful"
    grep -o '"text":"[^"]*"' /tmp/mcp-call.json | sed 's/"text":"\([^"]*\)"/Result: \1/'
else
    echo "❌ Tool call failed"
    cat /tmp/mcp-call.json
fi

echo ""
echo "🎯 Configuration for Cursor (.cursor/mcp.json):"
echo "{"
echo '  "mcpServers": {'
echo '    "url-db": {'
echo "      \"command\": \"$(pwd)/bin/url-db\","
echo '      "args": ['
echo '        "-mcp-mode=stdio",'
echo "        \"-db-path=$(pwd)/url-db.sqlite\""
echo '      ],'
echo '      "env": {}'
echo '    }'
echo '  }'
echo "}"

echo ""
echo "🔧 Next steps:"
echo "1. Copy the configuration above to ~/.cursor/mcp.json"
echo "2. Restart Cursor completely"
echo "3. Check if MCP server appears in Cursor's tool list"
echo ""
echo "💡 If still having issues:"
echo "- Check Cursor's developer console for errors"
echo "- Verify the absolute path is correct"
echo "- Try running this debug script again"

# Cleanup
rm -f /tmp/mcp-*.json