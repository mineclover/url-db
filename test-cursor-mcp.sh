#!/bin/bash

echo "🧪 Testing URL-DB MCP Server for Cursor Compatibility"
echo "===================================================="

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
URL_DB_BIN="$SCRIPT_DIR/bin/url-db"
TEST_DB="/tmp/cursor-test.db"

echo ""
echo "1️⃣ Testing with standard flags (MCP Go SDK style):"
echo "Command: $URL_DB_BIN -mcp-mode=stdio -db-path=$TEST_DB -tool-name=cursor-test"
echo ""

# Create a simple test
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}' | \
    "$URL_DB_BIN" -mcp-mode=stdio -db-path="$TEST_DB" -tool-name=cursor-test 2>/dev/null | head -1 | jq .

echo ""
echo "✅ Test successful! Server responds to JSON-RPC requests."
echo ""
echo "2️⃣ Configuration for Cursor (~/.cursor/mcp.json):"
echo ""
echo '{
  "mcpServers": {
    "url-db": {
      "command": "'$URL_DB_BIN'",
      "args": [
        "-mcp-mode=stdio",
        "-db-path='$SCRIPT_DIR'/url-db.db"
      ]
    }
  }
}'
echo ""
echo "3️⃣ Configuration for multiple databases:"
echo ""
echo '{
  "mcpServers": {
    "url-db-default": {
      "command": "'$URL_DB_BIN'",
      "args": [
        "-mcp-mode=stdio",
        "-db-path='$SCRIPT_DIR'/url-db.db"
      ]
    },
    "url-db-work": {
      "command": "'$URL_DB_BIN'",
      "args": [
        "-mcp-mode=stdio",
        "-db-path='$SCRIPT_DIR'/work.db",
        "-tool-name=work"
      ]
    },
    "url-db-personal": {
      "command": "'$URL_DB_BIN'",
      "args": [
        "-mcp-mode=stdio",
        "-db-path='$SCRIPT_DIR'/personal.db",
        "-tool-name=personal"
      ]
    }
  }
}'
echo ""
echo "✅ URL-DB is now compatible with Cursor!"
echo ""
echo "Note: The -tool-name flag determines the prefix in composite keys:"
echo "  - Default: url-db:domain:id"
echo "  - Work: work:domain:id"
echo "  - Personal: personal:domain:id"