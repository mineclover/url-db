#!/bin/bash

# Full MCP test script
echo "Testing full MCP handshake..."

# Create a named pipe for bidirectional communication
PIPE_IN=$(mktemp -u)
PIPE_OUT=$(mktemp -u)
mkfifo "$PIPE_IN"
mkfifo "$PIPE_OUT"

# Start the server in background
/Users/junwoobang/mcp/url-db/bin/url-db -mcp-mode=stdio -db-path=/tmp/test-mcp.db < "$PIPE_IN" > "$PIPE_OUT" 2>/tmp/mcp-error.log &
SERVER_PID=$!

# Function to send request and read response
send_request() {
    echo "$1" > "$PIPE_IN"
    # Read response with timeout
    timeout 2 cat "$PIPE_OUT" || echo "Timeout reading response"
}

# Keep the output pipe open
exec 3< "$PIPE_OUT"

# Send initialize
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {"tools": {}}}}' > "$PIPE_IN"
read -t 2 RESPONSE <&3
echo "Initialize response: $RESPONSE"

# Send tools/list
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}' > "$PIPE_IN"
read -t 2 RESPONSE <&3
echo "Tools list response: $RESPONSE"

# Check if server is still running
if kill -0 $SERVER_PID 2>/dev/null; then
    echo "Server is still running"
else
    echo "Server has exited"
    echo "Error log:"
    cat /tmp/mcp-error.log
fi

# Cleanup
kill $SERVER_PID 2>/dev/null
rm -f "$PIPE_IN" "$PIPE_OUT" /tmp/mcp-error.log
exec 3<&-