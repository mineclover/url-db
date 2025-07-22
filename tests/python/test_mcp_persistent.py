#!/usr/bin/env python3
"""Test MCP domain attribute operations with persistent connection"""

import json
import subprocess
import sys
import time
from tool_constants import CREATE_DOMAIN, CREATE_DOMAIN_ATTRIBUTE, DELETE_DOMAIN_ATTRIBUTE, GET_DOMAIN_ATTRIBUTE, LIST_DOMAIN_ATTRIBUTES, UPDATE_DOMAIN_ATTRIBUTE


class MCPClient:
    def __init__(self):
        self.proc = None
        self.request_id = 0
    
    def start(self):
        """Start the MCP server process"""
        self.proc = subprocess.Popen(
            ["../../bin/url-db", "-mcp-mode=stdio"],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            bufsize=1
        )
        time.sleep(0.5)  # Give server time to start
    
    def send_request(self, method, params=None):
        """Send a request and get response"""
        self.request_id += 1
        request = {
            "jsonrpc": "2.0",
            "id": self.request_id,
            "method": method
        }
        if params:
            request["params"] = params
        
        # Send request
        request_str = json.dumps(request)
        print(f"Sending: {request_str}")
        self.proc.stdin.write(request_str + "\n")
        self.proc.stdin.flush()
        
        # Read response
        response_line = self.proc.stdout.readline()
        print(f"Received: {response_line.strip()}")
        
        try:
            return json.loads(response_line)
        except json.JSONDecodeError as e:
            print(f"Failed to parse response: {response_line}")
            raise e
    
    def close(self):
        """Close the connection"""
        if self.proc:
            self.proc.terminate()
            self.proc.wait()

def test_domain_attributes():
    """Test domain attribute CRUD operations"""
    client = MCPClient()
    
    try:
        print("Starting MCP server...")
        client.start()
        
        # 1. Initialize
        print("\n1. Initializing server...")
        response = client.send_request("initialize", {
            "protocolVersion": "2024-11-05",
            "capabilities": {}
        })
        print(f"Response: {json.dumps(response, indent=2)}")
        
        # 1.5. Send initialized notification
        print("\n1.5. Sending initialized notification...")
        # Notification without ID
        notification = {
            "jsonrpc": "2.0",
            "method": "notifications/initialized"
        }
        print(f"Sending: {json.dumps(notification)}")
        client.proc.stdin.write(json.dumps(notification) + "\n")
        client.proc.stdin.flush()
        time.sleep(0.1)  # Give server time to process
        
        # 2. List tools
        print("\n2. Listing available tools...")
        response = client.send_request("tools/list")
        print(f"Response: {json.dumps(response, indent=2)}")
        
        if "result" in response and "tools" in response["result"]:
            tools = response["result"]["tools"]
            print(f"\nFound {len(tools)} tools:")
            for tool in tools:
                print(f"  - {tool['name']}")
        
        # 3. Create domain
        print("\n3. Creating test domain...")
        response = client.send_request("tools/call", {
            "name": CREATE_DOMAIN,
            "arguments": {
                "name": "test-domain",
                "description": "Test domain for attribute testing"
            }
        })
        print(f"Response: {json.dumps(response, indent=2)}")
        
        # 4. List domain attributes
        print("\n4. Listing domain attributes...")
        response = client.send_request("tools/call", {
            "name": LIST_DOMAIN_ATTRIBUTES,
            "arguments": {
                "domain_name": "test-domain"
            }
        })
        print(f"Response: {json.dumps(response, indent=2)}")
        
        # 5. Create domain attribute
        print("\n5. Creating domain attribute...")
        response = client.send_request("tools/call", {
            "name": CREATE_DOMAIN_ATTRIBUTE,
            "arguments": {
                "domain_name": "test-domain",
                "name": "category",
                "type": "tag",
                "description": "Category tag for URLs"
            }
        })
        print(f"Response: {json.dumps(response, indent=2)}")
        
        if "result" in response and not response["result"].get("isError", False):
            result_content = response["result"]["content"][0]["text"]
            result_data = json.loads(result_content)
            composite_id = result_data["composite_id"]
            print(f"\nCreated attribute with ID: {composite_id}")
        elif response["result"].get("isError", False) and "already exists" in response["result"]["content"][0]["text"]:
            print("\n✓ Attribute already exists, using existing one")
            # Get existing attribute from list
            list_response = client.send_request("tools/call", {
                "name": LIST_DOMAIN_ATTRIBUTES,
                "arguments": {
                    "domain_name": "test-domain"
                }
            })
            if list_response and not list_response["result"].get("isError", False):
                list_content = json.loads(list_response["result"]["content"][0]["text"])
                if list_content["attributes"]:
                    composite_id = list_content["attributes"][0]["composite_id"]
                    print(f"Using existing attribute ID: {composite_id}")
                else:
                    print("❌ No attributes found")
                    return
            else:
                print("❌ Failed to get attribute list")
                return
        else:
            print(f"❌ Failed to create or find attribute: {response}")
            return
        
        # 6. Get attribute
        print("\n6. Getting attribute by composite ID...")
        response = client.send_request("tools/call", {
            "name": GET_DOMAIN_ATTRIBUTE,
            "arguments": {
                "composite_id": composite_id
            }
        })
        print(f"Response: {json.dumps(response, indent=2)}")
        
        # 7. Update attribute
        print("\n7. Updating attribute description...")
        response = client.send_request("tools/call", {
            "name": UPDATE_DOMAIN_ATTRIBUTE,
            "arguments": {
                "composite_id": composite_id,
                "description": "Updated category tag description"
            }
        })
        print(f"Response: {json.dumps(response, indent=2)}")
        
        # 8. List attributes again
        print("\n8. Listing domain attributes after update...")
        response = client.send_request("tools/call", {
            "name": LIST_DOMAIN_ATTRIBUTES,
            "arguments": {
                "domain_name": "test-domain"
            }
        })
        print(f"Response: {json.dumps(response, indent=2)}")
        
        # 9. Delete attribute
        print("\n9. Deleting attribute...")
        response = client.send_request("tools/call", {
            "name": DELETE_DOMAIN_ATTRIBUTE,
            "arguments": {
                "composite_id": composite_id
            }
        })
        print(f"Response: {json.dumps(response, indent=2)}")
        
        print("\n✓ All tests completed!")
        
    except Exception as e:
        print(f"\n✗ Test failed: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
    finally:
        client.close()

if __name__ == "__main__":
    test_domain_attributes()