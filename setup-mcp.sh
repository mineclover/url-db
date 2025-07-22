#!/bin/bash

# URL-DB MCP Setup Script for Claude Desktop

echo "🚀 Setting up URL-DB MCP Server for Claude Desktop"
echo "=================================================="

# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
URL_DB_BIN="$SCRIPT_DIR/bin/url-db"
DATABASE_PATH="$SCRIPT_DIR/url-db.db"

# Check if the binary exists
if [ ! -f "$URL_DB_BIN" ]; then
    echo "❌ Binary not found at: $URL_DB_BIN"
    echo "   Running build script..."
    cd "$SCRIPT_DIR"
    ./build.sh
    if [ $? -ne 0 ]; then
        echo "❌ Build failed!"
        exit 1
    fi
fi

# Check if the binary is executable
if [ ! -x "$URL_DB_BIN" ]; then
    echo "🔧 Making binary executable..."
    chmod +x "$URL_DB_BIN"
fi

echo ""
echo "✅ URL-DB binary found at: $URL_DB_BIN"
echo "📁 Database will be at: $DATABASE_PATH"
echo ""
echo "To add URL-DB to Claude Desktop, run ONE of these commands:"
echo ""
echo "Option 1 - Using environment variables (recommended):"
echo "--------------------------------------------------------"
echo "claude mcp add url-db \"$URL_DB_BIN\" \\"
echo "  --args=\"-mcp-mode=stdio\" \\"
echo "  --env=\"DATABASE_URL=file:$DATABASE_PATH\""
echo ""
echo "Option 2 - Using command line arguments:"
echo "--------------------------------------------------------"
echo "claude mcp add url-db -- \"$URL_DB_BIN\" -mcp-mode=stdio DATABASE_URL=\"file:$DATABASE_PATH\""
echo ""
echo "Option 3 - Manual configuration:"
echo "--------------------------------------------------------"
echo "Add this to your Claude Desktop config file:"
echo "~/Library/Application Support/Claude/claude_desktop_config.json"
echo ""
echo '{
  "mcpServers": {
    "url-db": {
      "command": "'$URL_DB_BIN'",
      "args": ["-mcp-mode=stdio"],
      "env": {
        "DATABASE_URL": "file:'$DATABASE_PATH'"
      }
    }
  }
}'
echo ""
echo "After adding, restart Claude Desktop to activate the MCP server."
echo ""
echo "To verify installation:"
echo "claude mcp list"
echo ""
echo "To test the server directly:"
echo "\"$URL_DB_BIN\" -mcp-mode=stdio"