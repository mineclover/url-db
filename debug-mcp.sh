#!/bin/bash

# Debug wrapper for MCP server
echo "MCP Debug: Starting with args: $@" >&2
echo "MCP Debug: Working directory: $(pwd)" >&2
echo "MCP Debug: Date: $(date)" >&2

# Run the actual server
exec /Users/junwoobang/mcp/url-db/bin/url-db "$@"