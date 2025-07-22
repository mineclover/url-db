#!/usr/bin/env python3
"""Test MCP domain attribute operations"""

import json
import subprocess
import sys
from tool_constants import CREATE_DOMAIN, CREATE_DOMAIN_ATTRIBUTE, DELETE_DOMAIN_ATTRIBUTE, GET_DOMAIN_ATTRIBUTE, LIST_DOMAIN_ATTRIBUTES, UPDATE_DOMAIN_ATTRIBUTE


def send_mcp_request(method, params=None):
    """Send an MCP request via stdio"""
    request = {
        "jsonrpc": "2.0",
        "id": 1,
        "method": method
    }
    if params:
        request["params"] = params
    
    # Start the MCP server
    proc = subprocess.Popen(
        ["../../bin/url-db", "-mcp-mode=stdio"],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True
    )
    
    # Send request
    proc.stdin.write(json.dumps(request) + "\n")
    proc.stdin.flush()
    
    # Read response
    response_line = proc.stdout.readline()
    proc.terminate()
    
    try:
        return json.loads(response_line)
    except json.JSONDecodeError as e:
        print(f"Failed to parse response: {response_line}")
        raise e

def test_domain_attributes():
    """Test domain attribute CRUD operations"""
    print("Testing MCP Domain Attribute Operations...")
    
    # -1. Initialize the server first
    print("\n-1. Initializing server...")
    response = send_mcp_request("initialize", {
        "protocolVersion": "0.1.0",
        "capabilities": {}
    })
    print(f"Response: {json.dumps(response, indent=2)}")
    
    # 0. List available tools first
    print("\n0. Listing available tools...")
    response = send_mcp_request("tools/list", {})
    print(f"Response: {json.dumps(response, indent=2)}")
    
    # 1. Create a test domain first
    print("\n1. Creating test domain...")
    response = send_mcp_request(CREATE_DOMAIN, {
        "name": "test-domain",
        "description": "Test domain for attribute testing"
    })
    print(f"Response: {json.dumps(response, indent=2)}")
    
    # 2. List domain attributes (should be empty)
    print("\n2. Listing domain attributes...")
    response = send_mcp_request(LIST_DOMAIN_ATTRIBUTES, {
        "domain_name": "test-domain"
    })
    print(f"Response: {json.dumps(response, indent=2)}")
    
    # 3. Create a domain attribute
    print("\n3. Creating domain attribute...")
    response = send_mcp_request(CREATE_DOMAIN_ATTRIBUTE, {
        "domain_name": "test-domain",
        "name": "category",
        "type": "tag",
        "description": "Category tag for URLs"
    })
    print(f"Response: {json.dumps(response, indent=2)}")
    
    if "result" in response:
        composite_id = response["result"]["composite_id"]
        print(f"Created attribute with ID: {composite_id}")
        
        # 4. Get the attribute by composite ID
        print("\n4. Getting attribute by composite ID...")
        response = send_mcp_request(GET_DOMAIN_ATTRIBUTE, {
            "composite_id": composite_id
        })
        print(f"Response: {json.dumps(response, indent=2)}")
        
        # 5. Update the attribute
        print("\n5. Updating attribute description...")
        response = send_mcp_request(UPDATE_DOMAIN_ATTRIBUTE, {
            "composite_id": composite_id,
            "description": "Updated category tag description"
        })
        print(f"Response: {json.dumps(response, indent=2)}")
        
        # 6. List attributes again
        print("\n6. Listing domain attributes after creation...")
        response = send_mcp_request(LIST_DOMAIN_ATTRIBUTES, {
            "domain_name": "test-domain"
        })
        print(f"Response: {json.dumps(response, indent=2)}")
        
        # 7. Delete the attribute
        print("\n7. Deleting attribute...")
        response = send_mcp_request(DELETE_DOMAIN_ATTRIBUTE, {
            "composite_id": composite_id
        })
        print(f"Response: {json.dumps(response, indent=2)}")

if __name__ == "__main__":
    try:
        test_domain_attributes()
        print("\n✓ All tests completed!")
    except Exception as e:
        print(f"\n✗ Test failed: {e}")
        sys.exit(1)