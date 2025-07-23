#!/bin/bash

# MCP wrapper script for debugging
LOG_FILE="/tmp/mcp-debug-$(date +%Y%m%d-%H%M%S).log"

echo "MCP Wrapper started at $(date)" >> "$LOG_FILE"
echo "Arguments: $@" >> "$LOG_FILE"
echo "Working directory: $(pwd)" >> "$LOG_FILE"

# Change to project directory to ensure go.mod is found
cd /Users/junwoobang/mcp/url-db

echo "Changed to directory: $(pwd)" >> "$LOG_FILE"

# Run the actual server and log everything
exec /Users/junwoobang/mcp/url-db/bin/url-db "$@" 2>> "$LOG_FILE"