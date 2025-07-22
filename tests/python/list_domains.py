#!/usr/bin/env python3
"""
Script to list all domains in the URL database using MCP
"""

import json
import subprocess
import sys
import time
from tool_constants import LIST_DOMAINS


class MCPClient:
    def __init__(self, server_path):
        self.server_path = server_path
        self.process = None
        self.request_id = 1
        
    def start_server(self):
        """Start MCP server in stdio mode"""
        self.process = subprocess.Popen(
            [self.server_path, "-mcp-mode=stdio"],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            bufsize=0
        )
        
    def send_request(self, method, params=None):
        """Send JSON-RPC 2.0 request"""
        request = {
            "jsonrpc": "2.0",
            "id": self.request_id,
            "method": method,
            "params": params or {}
        }
        
        request_json = json.dumps(request)
        
        self.process.stdin.write(request_json + "\n")
        self.process.stdin.flush()
        
        # Read response
        response_line = self.process.stdout.readline().strip()
        
        if response_line:
            try:
                response = json.loads(response_line)
                self.request_id += 1
                return response
            except json.JSONDecodeError as e:
                print(f"Failed to parse response: {e}")
                return None
        
        return None
        
    def send_notification(self, method, params=None):
        """Send JSON-RPC 2.0 notification (no response expected)"""
        notification = {
            "jsonrpc": "2.0",
            "method": method,
            "params": params or {}
        }
        
        notification_json = json.dumps(notification)
        
        self.process.stdin.write(notification_json + "\n")
        self.process.stdin.flush()
        
    def close(self):
        """Close server"""
        if self.process:
            self.process.stdin.close()
            self.process.wait()

def main():
    # Server binary path
    server_path = "../../bin/url-db"
    
    client = MCPClient(server_path)
    
    try:
        # Start MCP server
        client.start_server()
        
        # Wait briefly for server to start
        time.sleep(0.5)
        
        # 1. Initialize handshake
        init_response = client.send_request("initialize", {
            "protocolVersion": "2024-11-05",
            "capabilities": {
                "experimental": {},
                "sampling": {}
            },
            "clientInfo": {
                "name": "list-domains-client",
                "version": "1.0.0"
            }
        })
        
        if not init_response or "result" not in init_response:
            print("Failed to initialize MCP connection")
            return
            
        # 2. Send initialized notification
        client.send_notification("initialized", {})
        
        # Wait briefly
        time.sleep(0.2)
        
        # 3. Call list_domains tool
        domains_response = client.send_request("tools/call", {
            "name": LIST_DOMAINS,
            "arguments": {}
        })
        
        if domains_response and "result" in domains_response:
            result = domains_response["result"]
            if not result.get("isError", False):
                # Parse and display the domains
                content = result["content"][0]["text"]
                domains_data = json.loads(content)
                
                print("Domains in the database:")
                print("-" * 50)
                
                if "domains" in domains_data:
                    for domain in domains_data["domains"]:
                        print(f"\nDomain: {domain['name']}")
                        print(f"Description: {domain['description']}")
                        print(f"Created: {domain['created_at']}")
                        print(f"Updated: {domain['updated_at']}")
                else:
                    print("No domains found in the database.")
            else:
                print("Error calling list_domains:", result["content"][0]["text"])
        else:
            print("Failed to get domains list")
        
    except Exception as e:
        print(f"Error: {e}")
        import traceback
        traceback.print_exc()
        
    finally:
        client.close()

if __name__ == "__main__":
    main()