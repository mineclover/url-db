#!/usr/bin/env python3
"""
Debug test to identify the JSON parsing issue
"""

import json
import sys
import subprocess
import time
from tool_constants import LIST_DOMAINS, CREATE_DOMAIN

class MCPClient:
    def __init__(self, db_path: str = "./debug_test.db"):
        """Initialize MCP client with stdio connection"""
        self.process = subprocess.Popen(
            ["../../bin/url-db", "-mcp-mode=stdio", f"-db-path={db_path}"],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            bufsize=0
        )
        self.request_id = 0
        
    def send_request(self, method: str, params: dict = None) -> dict:
        """Send JSON-RPC request and get response"""
        self.request_id += 1
        request = {
            "jsonrpc": "2.0",
            "id": self.request_id,
            "method": method,
            "params": params or {}
        }
        
        request_str = json.dumps(request) + "\n"
        print(f">>> Sending: {request_str.strip()}")
        
        self.process.stdin.write(request_str)
        self.process.stdin.flush()
        
        response_str = self.process.stdout.readline()
        print(f"<<< Received: {response_str.strip()}")
        
        try:
            response = json.loads(response_str) if response_str.strip() else {}
            print(f"    Parsed response: {response}")
            return response
        except json.JSONDecodeError as e:
            print(f"    JSON decode error: {e}")
            print(f"    Raw response: {repr(response_str)}")
            return {"error": {"message": "Invalid JSON response"}}
    
    def close(self):
        """Close the MCP connection"""
        self.process.terminate()
        self.process.wait()

def main():
    """Debug the JSON parsing issues"""
    print("🐛 Debug test for JSON parsing issues")
    print("=" * 50)
    
    client = MCPClient("./debug_test.db")
    
    try:
        # Initialize
        print("\n1. Initialize:")
        response = client.send_request("initialize", {
            "protocolVersion": "2024-11-05",
            "capabilities": {},
            "clientInfo": {"name": "debug-client", "version": "1.0"}
        })
        
        if "result" not in response:
            print("❌ Initialize failed!")
            return 1
        
        # Initialized notification
        print("\n2. Initialized notification:")
        client.send_request("notifications/initialized")
        
        # List domains
        print("\n3. List domains:")
        response = client.send_request("tools/call", {
            "name": LIST_DOMAINS,
            "arguments": {}
        })
        
        if "result" in response:
            try:
                content = response['result']['content'][0]['text']
                print(f"    Content: {repr(content)}")
                if content.strip():
                    domains_data = json.loads(content)
                    print(f"    Parsed domains: {domains_data}")
                else:
                    print("    Empty content!")
            except Exception as e:
                print(f"    Error processing domains: {e}")
        
        # Create domain
        print("\n4. Create domain:")
        response = client.send_request("tools/call", {
            "name": CREATE_DOMAIN,
            "arguments": {
                "name": "debug-domain",
                "description": "Debug test domain"
            }
        })
        
        if "result" in response:
            try:
                content = response['result']['content'][0]['text']
                print(f"    Content: {repr(content)}")
                if content.strip():
                    domain_data = json.loads(content)
                    print(f"    Parsed domain: {domain_data}")
                else:
                    print("    Empty content!")
            except Exception as e:
                print(f"    Error processing domain: {e}")
        
        # List domains again
        print("\n5. List domains again:")
        response = client.send_request("tools/call", {
            "name": LIST_DOMAINS,
            "arguments": {}
        })
        
        if "result" in response:
            try:
                content = response['result']['content'][0]['text']
                print(f"    Content: {repr(content)}")
                if content.strip():
                    domains_data = json.loads(content)
                    print(f"    Parsed domains: {domains_data}")
                    print(f"    Domain count: {len(domains_data.get('domains', []))}")
                else:
                    print("    Empty content!")
            except Exception as e:
                print(f"    Error processing domains: {e}")
        
        print("\n✅ Debug test completed successfully!")
        
    except Exception as e:
        print(f"\n❌ Debug test failed: {e}")
        import traceback
        traceback.print_exc()
        return 1
    finally:
        client.close()
    
    return 0

if __name__ == "__main__":
    sys.exit(main())