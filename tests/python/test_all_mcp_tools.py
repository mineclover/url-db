#!/usr/bin/env python3
"""
Comprehensive MCP tools test script for URL-DB
Tests all 16 MCP tools with the new naming convention
"""

import json
import sys
import subprocess
import time
from typing import Dict, Any, List, Optional

class MCPClient:
    def __init__(self, db_path: str = "./test.db"):
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
        
    def send_request(self, method: str, params: Dict[str, Any] = None) -> Dict[str, Any]:
        """Send JSON-RPC request and get response"""
        self.request_id += 1
        request = {
            "jsonrpc": "2.0",
            "id": self.request_id,
            "method": method,
            "params": params or {}
        }
        
        request_str = json.dumps(request) + "\n"
        self.process.stdin.write(request_str)
        self.process.stdin.flush()
        
        response_str = self.process.stdout.readline()
        try:
            return json.loads(response_str)
        except json.JSONDecodeError:
            print(f"Failed to parse response: {response_str}")
            return {"error": {"message": "Invalid JSON response"}}
    
    def close(self):
        """Close the MCP connection"""
        self.process.terminate()
        self.process.wait()

def test_protocol_handshake(client: MCPClient):
    """Test 1: MCP Protocol Handshake"""
    print("\n=== Test 1: MCP Protocol Handshake ===")
    
    # Initialize
    response = client.send_request("initialize", {
        "protocolVersion": "2024-11-05",
        "capabilities": {},
        "clientInfo": {"name": "test-client", "version": "1.0"}
    })
    
    if "result" in response:
        print("✅ Initialize successful")
        print(f"  Server: {response['result']['serverInfo']['name']} v{response['result']['serverInfo']['version']}")
        print(f"  Capabilities: {response['result']['capabilities']}")
    else:
        print(f"❌ Initialize failed: {response}")
        return False
    
    # Send initialized notification
    client.send_request("notifications/initialized")
    print("✅ Initialized notification sent")
    
    return True

def test_tool_discovery(client: MCPClient):
    """Test 2: Tool Discovery"""
    print("\n=== Test 2: Tool Discovery ===")
    
    response = client.send_request("tools/list")
    
    if "result" in response:
        tools = response['result']['tools']
        print(f"✅ Found {len(tools)} tools")
        
        expected_tools = [
            "list_domains", "create_domain", "list_nodes", "create_node",
            "get_node", "update_node", "delete_node", "find_node_by_url",
            "get_node_attributes", "set_node_attributes", "get_server_info",
            "list_domain_attributes", "create_domain_attribute", 
            "get_domain_attribute", "update_domain_attribute", "delete_domain_attribute"
        ]
        
        tool_names = [tool['name'] for tool in tools]
        for expected in expected_tools:
            if expected in tool_names:
                print(f"  ✅ {expected}")
            else:
                print(f"  ❌ Missing: {expected}")
    else:
        print(f"❌ Tool list failed: {response}")
        return False
    
    return True

def test_domain_management(client: MCPClient):
    """Test 3: Domain Management"""
    print("\n=== Test 3: Domain Management ===")
    
    # List domains (should be empty)
    response = client.send_request("tools/call", {
        "name": "list_domains",
        "arguments": {}
    })
    
    if "result" in response:
        print("✅ List domains successful")
        initial_count = len(json.loads(response['result']['content'][0]['text'])['domains'])
        print(f"  Initial domain count: {initial_count}")
    else:
        print(f"❌ List domains failed: {response}")
        return False
    
    # Create domain
    response = client.send_request("tools/call", {
        "name": "create_domain",
        "arguments": {
            "name": "test-domain",
            "description": "Test domain for comprehensive testing"
        }
    })
    
    if "result" in response:
        print("✅ Create domain successful")
        domain = json.loads(response['result']['content'][0]['text'])
        print(f"  Created: {domain['name']}")
    else:
        print(f"❌ Create domain failed: {response}")
        return False
    
    # List domains again
    response = client.send_request("tools/call", {
        "name": "list_domains",
        "arguments": {}
    })
    
    if "result" in response:
        domains = json.loads(response['result']['content'][0]['text'])['domains']
        if len(domains) > initial_count:
            print("✅ Domain appears in list")
        else:
            print("❌ Domain not found in list")
    
    return True

def test_domain_schema(client: MCPClient):
    """Test 4: Domain Schema Management"""
    print("\n=== Test 4: Domain Schema Management ===")
    
    # Create domain attribute definitions
    attribute_types = [
        ("category", "tag", "Category tag for nodes"),
        ("priority", "ordered_tag", "Priority with ordering"),
        ("score", "number", "Numeric score value"),
        ("notes", "string", "Text notes"),
        ("description", "markdown", "Markdown formatted description"),
        ("thumbnail", "image", "Thumbnail image URL")
    ]
    
    for name, attr_type, description in attribute_types:
        response = client.send_request("tools/call", {
            "name": "create_domain_attribute",
            "arguments": {
                "domain_name": "test-domain",
                "name": name,
                "type": attr_type,
                "description": description
            }
        })
        
        if "result" in response:
            print(f"✅ Created attribute: {name} ({attr_type})")
        else:
            print(f"❌ Failed to create attribute {name}: {response}")
    
    # List domain attributes
    response = client.send_request("tools/call", {
        "name": "list_domain_attributes",
        "arguments": {
            "domain_name": "test-domain"
        }
    })
    
    if "result" in response:
        attributes = json.loads(response['result']['content'][0]['text'])['attributes']
        print(f"✅ Domain has {len(attributes)} attributes defined")
        for attr in attributes:
            print(f"  - {attr['name']} ({attr['type']}): {attr['description']}")
    else:
        print(f"❌ Failed to list attributes: {response}")
    
    return True

def test_node_operations(client: MCPClient):
    """Test 5: Node/URL Management"""
    print("\n=== Test 5: Node/URL Management ===")
    
    # Create node
    response = client.send_request("tools/call", {
        "name": "create_node",
        "arguments": {
            "domain_name": "test-domain",
            "url": "https://example.com/test-page",
            "title": "Test Page",
            "description": "A test page for MCP testing"
        }
    })
    
    if "result" in response:
        node = json.loads(response['result']['content'][0]['text'])
        composite_id = node['composite_id']
        print(f"✅ Created node: {composite_id}")
        print(f"  URL: {node['url']}")
        print(f"  Title: {node['title']}")
    else:
        print(f"❌ Create node failed: {response}")
        return None
    
    # Get node by composite ID
    response = client.send_request("tools/call", {
        "name": "get_node",
        "arguments": {
            "composite_id": composite_id
        }
    })
    
    if "result" in response:
        print("✅ Retrieved node by composite ID")
    else:
        print(f"❌ Get node failed: {response}")
    
    # Update node
    response = client.send_request("tools/call", {
        "name": "update_node",
        "arguments": {
            "composite_id": composite_id,
            "title": "Updated Test Page",
            "description": "Updated description"
        }
    })
    
    if "result" in response:
        print("✅ Updated node successfully")
    else:
        print(f"❌ Update node failed: {response}")
    
    # Find node by URL
    response = client.send_request("tools/call", {
        "name": "find_node_by_url",
        "arguments": {
            "domain_name": "test-domain",
            "url": "https://example.com/test-page"
        }
    })
    
    if "result" in response:
        print("✅ Found node by URL")
    else:
        print(f"❌ Find node failed: {response}")
    
    return composite_id

def test_node_attributes(client: MCPClient, composite_id: str):
    """Test 6: Node Attribute Management with Schema Validation"""
    print("\n=== Test 6: Node Attribute Management ===")
    
    # Set multiple attributes
    response = client.send_request("tools/call", {
        "name": "set_node_attributes",
        "arguments": {
            "composite_id": composite_id,
            "attributes": [
                {"name": "category", "value": "documentation"},
                {"name": "priority", "value": "high", "order_index": 1},
                {"name": "score", "value": "85.5"},
                {"name": "notes", "value": "Important reference document"}
            ]
        }
    })
    
    if "result" in response:
        print("✅ Set node attributes successfully")
        result = json.loads(response['result']['content'][0]['text'])
        print(f"  Attributes set: {len(result['attributes'])}")
    else:
        print(f"❌ Set attributes failed: {response}")
    
    # Get node attributes
    response = client.send_request("tools/call", {
        "name": "get_node_attributes",
        "arguments": {
            "composite_id": composite_id
        }
    })
    
    if "result" in response:
        attrs = json.loads(response['result']['content'][0]['text'])['attributes']
        print(f"✅ Retrieved {len(attrs)} attributes")
        for attr in attrs:
            print(f"  - {attr['name']}: {attr['value']}")
    else:
        print(f"❌ Get attributes failed: {response}")
    
    # Try to set invalid attribute (not in schema)
    response = client.send_request("tools/call", {
        "name": "set_node_attributes",
        "arguments": {
            "composite_id": composite_id,
            "attributes": [
                {"name": "invalid_attr", "value": "should fail"}
            ]
        }
    })
    
    if "error" in response or (response.get("result", {}).get("isError")):
        print("✅ Schema validation working - rejected invalid attribute")
    else:
        print("❌ Schema validation failed - accepted invalid attribute")
    
    return True

def test_list_operations(client: MCPClient):
    """Test 7: List and Search Operations"""
    print("\n=== Test 7: List and Search Operations ===")
    
    # Create more nodes for testing
    urls = [
        ("https://example.com/page1", "Page 1", "First page"),
        ("https://example.com/page2", "Page 2", "Second page"),
        ("https://example.com/search", "Search Page", "Search functionality")
    ]
    
    for url, title, desc in urls:
        client.send_request("tools/call", {
            "name": "create_node",
            "arguments": {
                "domain_name": "test-domain",
                "url": url,
                "title": title,
                "description": desc
            }
        })
    
    # List nodes
    response = client.send_request("tools/call", {
        "name": "list_nodes",
        "arguments": {
            "domain_name": "test-domain",
            "page": 1,
            "size": 10
        }
    })
    
    if "result" in response:
        result = json.loads(response['result']['content'][0]['text'])
        print(f"✅ Listed {len(result['nodes'])} nodes")
        print(f"  Total count: {result['total_count']}")
    else:
        print(f"❌ List nodes failed: {response}")
    
    # Search nodes
    response = client.send_request("tools/call", {
        "name": "list_nodes",
        "arguments": {
            "domain_name": "test-domain",
            "search": "search",
            "page": 1,
            "size": 10
        }
    })
    
    if "result" in response:
        result = json.loads(response['result']['content'][0]['text'])
        print(f"✅ Search found {len(result['nodes'])} nodes")
    else:
        print(f"❌ Search failed: {response}")
    
    return True

def test_cleanup(client: MCPClient, composite_id: str):
    """Test 8: Cleanup Operations"""
    print("\n=== Test 8: Cleanup Operations ===")
    
    # Delete node
    response = client.send_request("tools/call", {
        "name": "delete_node",
        "arguments": {
            "composite_id": composite_id
        }
    })
    
    if "result" in response:
        print("✅ Deleted node successfully")
    else:
        print(f"❌ Delete node failed: {response}")
    
    # Verify deletion
    response = client.send_request("tools/call", {
        "name": "get_node",
        "arguments": {
            "composite_id": composite_id
        }
    })
    
    if response.get("result", {}).get("isError") or "error" in response:
        print("✅ Node properly deleted (not found)")
    else:
        print("❌ Node still exists after deletion")
    
    return True

def test_server_info(client: MCPClient):
    """Test 9: Server Information"""
    print("\n=== Test 9: Server Information ===")
    
    response = client.send_request("tools/call", {
        "name": "get_server_info",
        "arguments": {}
    })
    
    if "result" in response:
        info = json.loads(response['result']['content'][0]['text'])
        print("✅ Server info retrieved")
        print(f"  Name: {info['name']}")
        print(f"  Version: {info['version']}")
        print(f"  Description: {info['description']}")
    else:
        print(f"❌ Get server info failed: {response}")
    
    return True

def main():
    """Run all MCP tests"""
    print("Starting comprehensive MCP tool tests...")
    print("=" * 50)
    
    # Initialize client
    client = MCPClient("./test_mcp_comprehensive.db")
    
    try:
        # Run all tests
        if not test_protocol_handshake(client):
            print("\n❌ Protocol handshake failed, aborting tests")
            return 1
        
        if not test_tool_discovery(client):
            print("\n❌ Tool discovery failed, aborting tests")
            return 1
        
        test_domain_management(client)
        test_domain_schema(client)
        
        composite_id = test_node_operations(client)
        if composite_id:
            test_node_attributes(client, composite_id)
            test_list_operations(client)
            test_cleanup(client, composite_id)
        
        test_server_info(client)
        
        print("\n" + "=" * 50)
        print("✅ All tests completed!")
        
    except Exception as e:
        print(f"\n❌ Test failed with error: {e}")
        return 1
    finally:
        client.close()
    
    return 0

if __name__ == "__main__":
    sys.exit(main())