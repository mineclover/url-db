#!/usr/bin/env python3
"""Final test for MCP domain attribute operations"""

import json
import subprocess
import sys
import time
import random

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
        self.proc.stdin.write(request_str + "\n")
        self.proc.stdin.flush()
        
        # Read response
        response_line = self.proc.stdout.readline()
        
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
    
    # Generate unique names to avoid conflicts
    test_suffix = random.randint(1000, 9999)
    domain_name = f"test-domain-{test_suffix}"
    attr_name = f"category-{test_suffix}"
    
    try:
        print("Starting MCP server...")
        client.start()
        
        # 1. Initialize
        print("\n1. Initializing server...")
        response = client.send_request("initialize", {
            "protocolVersion": "2024-11-05",
            "capabilities": {}
        })
        assert response.get("result"), f"Initialize failed: {response}"
        
        # Send initialized notification
        notification = {
            "jsonrpc": "2.0",
            "method": "notifications/initialized"
        }
        client.proc.stdin.write(json.dumps(notification) + "\n")
        client.proc.stdin.flush()
        time.sleep(0.1)
        
        # 2. Create domain
        print(f"\n2. Creating domain '{domain_name}'...")
        response = client.send_request("tools/call", {
            "name": "create_domain",
            "arguments": {
                "name": domain_name,
                "description": "Test domain for attribute testing"
            }
        })
        assert not response["result"].get("isError"), f"Create domain failed: {response}"
        print("✓ Domain created successfully")
        
        # 3. List domain attributes (should be empty)
        print("\n3. Listing domain attributes (should be empty)...")
        response = client.send_request("tools/call", {
            "name": "list_domain_attributes",
            "arguments": {
                "domain_name": domain_name
            }
        })
        assert not response["result"].get("isError"), f"List attributes failed: {response}"
        
        # Parse the response
        content = response["result"]["content"][0]["text"]
        data = json.loads(content)
        assert data["total_count"] == 0, "Expected 0 attributes"
        print("✓ No attributes found (as expected)")
        
        # 4. Create domain attribute
        print(f"\n4. Creating domain attribute '{attr_name}'...")
        response = client.send_request("tools/call", {
            "name": "create_domain_attribute",
            "arguments": {
                "domain_name": domain_name,
                "name": attr_name,
                "type": "tag",
                "description": "Test category tag"
            }
        })
        assert not response["result"].get("isError"), f"Create attribute failed: {response}"
        
        # Parse the response to get composite_id
        content = response["result"]["content"][0]["text"]
        attr_data = json.loads(content)
        composite_id = attr_data["composite_id"]
        print(f"✓ Attribute created with ID: {composite_id}")
        
        # Verify the composite ID format
        assert composite_id.startswith("url-db:"), f"Invalid composite ID format: {composite_id}"
        assert ":attr-" in composite_id, f"Missing 'attr-' prefix in composite ID: {composite_id}"
        print("✓ Composite ID has correct format for attributes")
        
        # 5. Get attribute by composite ID
        print("\n5. Getting attribute by composite ID...")
        response = client.send_request("tools/call", {
            "name": "get_domain_attribute",
            "arguments": {
                "composite_id": composite_id
            }
        })
        assert not response["result"].get("isError"), f"Get attribute failed: {response}"
        print("✓ Attribute retrieved successfully")
        
        # 6. Update attribute
        print("\n6. Updating attribute description...")
        response = client.send_request("tools/call", {
            "name": "update_domain_attribute",
            "arguments": {
                "composite_id": composite_id,
                "description": "Updated test category tag"
            }
        })
        assert not response["result"].get("isError"), f"Update attribute failed: {response}"
        print("✓ Attribute updated successfully")
        
        # 7. List attributes again
        print("\n7. Listing domain attributes after creation...")
        response = client.send_request("tools/call", {
            "name": "list_domain_attributes",
            "arguments": {
                "domain_name": domain_name
            }
        })
        assert not response["result"].get("isError"), f"List attributes failed: {response}"
        
        content = response["result"]["content"][0]["text"]
        data = json.loads(content)
        assert data["total_count"] == 1, f"Expected 1 attribute, got {data['total_count']}"
        assert data["attributes"][0]["description"] == "Updated test category tag"
        print("✓ Found 1 attribute with updated description")
        
        # 8. Delete attribute
        print("\n8. Deleting attribute...")
        response = client.send_request("tools/call", {
            "name": "delete_domain_attribute",
            "arguments": {
                "composite_id": composite_id
            }
        })
        assert not response["result"].get("isError"), f"Delete attribute failed: {response}"
        print("✓ Attribute deleted successfully")
        
        # 9. Verify deletion
        print("\n9. Verifying deletion...")
        response = client.send_request("tools/call", {
            "name": "list_domain_attributes",
            "arguments": {
                "domain_name": domain_name
            }
        })
        assert not response["result"].get("isError"), f"List attributes failed: {response}"
        
        content = response["result"]["content"][0]["text"]
        data = json.loads(content)
        assert data["total_count"] == 0, f"Expected 0 attributes after deletion, got {data['total_count']}"
        print("✓ Attribute successfully deleted")
        
        print("\n✅ All tests passed!")
        print("\nKey findings:")
        print("1. NodeAttributeWithInfo struct was missing db tags - FIXED")
        print("2. Attribute composite IDs now use 'attr-' prefix for differentiation")
        print("3. MCP domain attribute CRUD operations work correctly")
        
    except Exception as e:
        print(f"\n❌ Test failed: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
    finally:
        client.close()

if __name__ == "__main__":
    test_domain_attributes()