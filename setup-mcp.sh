#!/bin/bash

# URL-DB MCP Setup Script for Claude Desktop and Cursor

echo "🚀 Setting up URL-DB MCP Server"
echo "================================"

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
echo "📋 For Claude Desktop:"
echo "====================="
echo "Run this command:"
echo ""
echo "claude mcp add url-db \"$URL_DB_BIN\" \\"
echo "  --args=\"-mcp-mode=stdio\" \\"
echo "  --args=\"-db-path=$DATABASE_PATH\""
echo ""
echo "📋 For Cursor IDE:"
echo "=================="
echo "Add this to ~/.cursor/mcp.json:"
echo ""
echo '{
  "mcpServers": {
    "url-db": {
      "command": "'$URL_DB_BIN'",
      "args": [
        "-mcp-mode=stdio",
        "-db-path='$DATABASE_PATH'"
      ]
    }
  }
}'
echo ""
echo "📋 For Multiple Databases:"
echo "=========================="
echo "You can create multiple instances with different databases:"
echo ""
echo "# Work database"
echo "claude mcp add url-db-work \"$URL_DB_BIN\" \\"
echo "  --args=\"-mcp-mode=stdio\" \\"
echo "  --args=\"-db-path=$SCRIPT_DIR/work.db\" \\"
echo "  --args=\"-tool-name=work\""
echo ""
echo "# Personal database"
echo "claude mcp add url-db-personal \"$URL_DB_BIN\" \\"
echo "  --args=\"-mcp-mode=stdio\" \\"
echo "  --args=\"-db-path=$SCRIPT_DIR/personal.db\" \\"
echo "  --args=\"-tool-name=personal\""
echo ""
echo "📋 Manual Configuration Files:"
echo "=============================="
echo "Claude Desktop: ~/Library/Application Support/Claude/claude_desktop_config.json"
echo "Cursor: ~/.cursor/mcp.json"
echo ""
echo "After adding, restart your application to activate the MCP server."
echo ""
echo "To verify installation:"
echo "claude mcp list"
echo ""
echo "To test the server directly:"
echo "\"$URL_DB_BIN\" -mcp-mode=stdio -db-path=$DATABASE_PATH"