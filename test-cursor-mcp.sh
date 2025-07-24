#!/bin/bash

echo "🧪 Cursor MCP Connection Test"
echo "============================="

# Check if binary exists and is executable
if [ ! -x "./bin/url-db" ]; then
    echo "❌ Binary not executable: ./bin/url-db"
    echo "Run: chmod +x ./bin/url-db"
    exit 1
fi

echo "✅ Binary is executable"

# Check if config files exist
if [ -f "$HOME/.cursor/mcp.json" ]; then
    echo "✅ Global MCP config exists: ~/.cursor/mcp.json"
else
    echo "❌ Global MCP config missing: ~/.cursor/mcp.json"
fi

if [ -f "./.cursor/mcp.json" ]; then
    echo "✅ Project MCP config exists: ./.cursor/mcp.json"
else
    echo "❌ Project MCP config missing: ./.cursor/mcp.json"
fi

# Test MCP server directly
echo ""
echo "🔧 Testing MCP server directly..."
echo '{"jsonrpc": "2.0", "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}}, "id": 1}' | ./bin/url-db -mcp-mode=stdio -db-path=/Users/junwoobang/mcp/url-db/url-db.sqlite > /tmp/mcp-test.json 2>&1

if [ $? -eq 0 ]; then
    echo "✅ MCP server responds correctly"
    echo "Response preview:"
    cat /tmp/mcp-test.json | head -c 200
    echo "..."
else
    echo "❌ MCP server failed"
    cat /tmp/mcp-test.json
fi

echo ""
echo "📋 Next steps:"
echo "1. Open Cursor IDE"
echo "2. Go to Settings (Cmd+,)"
echo "3. Search for 'MCP' and enable MCP Servers"
echo "4. Restart Cursor completely"
echo "5. Check Settings > MCP for server status"
echo ""
echo "💡 Test in Cursor by asking:"
echo '   "What MCP tools are available?"'
echo '   "Can you list domains in my URL database?"'

rm -f /tmp/mcp-test.json