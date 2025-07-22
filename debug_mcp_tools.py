#!/usr/bin/env python3
"""
Debug MCP Tools - Show raw communication
"""

import json
import subprocess
import sys
import time

def debug_mcp_communication():
    """Debug MCP server communication"""
    print("🔍 Starting MCP Debug Session")
    
    # Start server
    process = subprocess.Popen(
        ["./cmd/server/url-db", "-mcp-mode=stdio"],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
        bufsize=0
    )
    
    time.sleep(0.5)
    
    def send_and_print(method, params=None):
        """Send request and print raw response"""
        request = {
            "jsonrpc": "2.0",
            "id": 1,
            "method": method,
            "params": params or {}
        }
        
        request_str = json.dumps(request)
        print(f"\n📤 Sending: {request_str}")
        
        process.stdin.write(request_str + '\n')
        process.stdin.flush()
        
        # Read response
        response_str = process.stdout.readline()
        print(f"📥 Response: {response_str}")
        
        # Try to parse response
        if response_str.strip():
            try:
                response_obj = json.loads(response_str)
                print(f"📊 Parsed: {json.dumps(response_obj, indent=2)}")
            except json.JSONDecodeError as e:
                print(f"❌ JSON parse error: {e}")
            
        return response_str
    
    # Test initialize
    print("\n1️⃣ Testing Initialize:")
    send_and_print("initialize", {
        "protocolVersion": "2024-11-05",
        "capabilities": {"roots": {"listChanged": True}},
        "clientInfo": {"name": "debug-client", "version": "1.0"}
    })
    
    # Send initialized
    print("\n2️⃣ Sending Initialized:")
    send_and_print("notifications/initialized")
    
    # Test tools/list
    print("\n3️⃣ Testing tools/list:")
    send_and_print("tools/list")
    
    # Test a tool call
    print("\n4️⃣ Testing tool call (get_mcp_server_info):")
    send_and_print("tools/call", {
        "name": "get_mcp_server_info",
        "arguments": {}
    })
    
    process.terminate()
    print("\n✅ Debug session complete")

if __name__ == "__main__":
    debug_mcp_communication()