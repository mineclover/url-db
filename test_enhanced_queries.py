#!/usr/bin/env python3
"""
Test script for enhanced MCP query tools
Tests get_node_with_attributes and filter_nodes_by_attributes
"""

import json
import sys
import subprocess
import time

class MCPClient:
    def __init__(self, db_path: str = "./test_enhanced.db"):
        """Initialize MCP client with stdio connection"""
        self.process = subprocess.Popen(
            ["./bin/url-db", "-mcp-mode=stdio", f"-db-path={db_path}"],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            bufsize=0
        )
        self.request_id = 0
        
    def send_request(self, method: str, params=None):
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

def setup_test_data(client):
    """Set up test domain and nodes with attributes"""
    print("\n=== Setting up test data ===")
    
    # Initialize
    client.send_request("initialize", {
        "protocolVersion": "2024-11-05",
        "capabilities": {},
        "clientInfo": {"name": "test-client", "version": "1.0"}
    })
    client.send_request("notifications/initialized")
    
    # Create domain
    print("Creating test domain...")
    client.send_request("tools/call", {
        "name": "create_domain",
        "arguments": {
            "name": "products",
            "description": "Product catalog for testing"
        }
    })
    
    # Create domain attributes
    print("Creating domain schema...")
    attributes = [
        ("category", "tag", "Product category"),
        ("price", "number", "Product price"),
        ("status", "tag", "Product status"),
        ("brand", "string", "Product brand")
    ]
    
    for name, attr_type, desc in attributes:
        client.send_request("tools/call", {
            "name": "create_domain_attribute",
            "arguments": {
                "domain_name": "products",
                "name": name,
                "type": attr_type,
                "description": desc
            }
        })
    
    # Create nodes with different attributes
    print("Creating test nodes...")
    products = [
        ("https://example.com/laptop1", "Laptop Pro 15", "High-end laptop", 
         [("category", "electronics"), ("price", "1299.99"), ("status", "available"), ("brand", "TechCorp")]),
        ("https://example.com/laptop2", "Budget Laptop", "Affordable laptop",
         [("category", "electronics"), ("price", "499.99"), ("status", "available"), ("brand", "ValueTech")]),
        ("https://example.com/phone1", "SmartPhone X", "Latest smartphone",
         [("category", "electronics"), ("price", "899.99"), ("status", "out-of-stock"), ("brand", "PhoneCo")]),
        ("https://example.com/book1", "Python Programming", "Learn Python",
         [("category", "books"), ("price", "39.99"), ("status", "available"), ("brand", "TechBooks")]),
        ("https://example.com/book2", "Data Science Guide", "DS fundamentals",
         [("category", "books"), ("price", "49.99"), ("status", "available"), ("brand", "TechBooks")])
    ]
    
    composite_ids = []
    for url, title, desc, attrs in products:
        # Create node
        response = client.send_request("tools/call", {
            "name": "create_node",
            "arguments": {
                "domain_name": "products",
                "url": url,
                "title": title,
                "description": desc
            }
        })
        
        if "result" in response:
            node = json.loads(response['result']['content'][0]['text'])
            composite_id = node['composite_id']
            composite_ids.append(composite_id)
            
            # Set attributes
            attr_list = [{"name": name, "value": value} for name, value in attrs]
            client.send_request("tools/call", {
                "name": "set_node_attributes",
                "arguments": {
                    "composite_id": composite_id,
                    "attributes": attr_list
                }
            })
    
    return composite_ids

def test_get_node_with_attributes(client, composite_id):
    """Test get_node_with_attributes tool"""
    print("\n=== Test 1: Get Node with Attributes ===")
    
    response = client.send_request("tools/call", {
        "name": "get_node_with_attributes",
        "arguments": {
            "composite_id": composite_id
        }
    })
    
    if "result" in response:
        data = json.loads(response['result']['content'][0]['text'])
        print(f"✅ Successfully retrieved node with attributes")
        print(f"  Node: {data['node']['title']} ({data['node']['url']})")
        print(f"  Attributes:")
        for attr in data['attributes']:
            print(f"    - {attr['name']}: {attr['value']}")
    else:
        print(f"❌ Failed to get node with attributes: {response}")

def test_filter_by_single_attribute(client):
    """Test filtering by a single attribute"""
    print("\n=== Test 2: Filter by Single Attribute ===")
    
    response = client.send_request("tools/call", {
        "name": "filter_nodes_by_attributes",
        "arguments": {
            "domain_name": "products",
            "filters": [
                {
                    "name": "category",
                    "value": "electronics",
                    "operator": "equals"
                }
            ],
            "page": 1,
            "size": 10
        }
    })
    
    if "result" in response:
        data = json.loads(response['result']['content'][0]['text'])
        print(f"✅ Found {data['total_count']} electronics products")
        for node in data['nodes']:
            print(f"  - {node['title']}")
    else:
        print(f"❌ Filter failed: {response}")

def test_filter_by_multiple_attributes(client):
    """Test filtering by multiple attributes"""
    print("\n=== Test 3: Filter by Multiple Attributes ===")
    
    response = client.send_request("tools/call", {
        "name": "filter_nodes_by_attributes",
        "arguments": {
            "domain_name": "products",
            "filters": [
                {
                    "name": "category",
                    "value": "electronics",
                    "operator": "equals"
                },
                {
                    "name": "status",
                    "value": "available",
                    "operator": "equals"
                }
            ],
            "page": 1,
            "size": 10
        }
    })
    
    if "result" in response:
        data = json.loads(response['result']['content'][0]['text'])
        print(f"✅ Found {data['total_count']} available electronics")
        for node in data['nodes']:
            print(f"  - {node['title']}")
    else:
        print(f"❌ Filter failed: {response}")

def test_filter_with_operators(client):
    """Test different filter operators"""
    print("\n=== Test 4: Filter with Different Operators ===")
    
    # Test contains operator
    print("\nTesting 'contains' operator for brand...")
    response = client.send_request("tools/call", {
        "name": "filter_nodes_by_attributes",
        "arguments": {
            "domain_name": "products",
            "filters": [
                {
                    "name": "brand",
                    "value": "Tech",
                    "operator": "contains"
                }
            ],
            "page": 1,
            "size": 10
        }
    })
    
    if "result" in response:
        data = json.loads(response['result']['content'][0]['text'])
        print(f"✅ Found {data['total_count']} products with 'Tech' in brand")
        for node in data['nodes']:
            print(f"  - {node['title']}")
    else:
        print(f"❌ Filter failed: {response}")

def test_pagination(client):
    """Test pagination in filtered results"""
    print("\n=== Test 5: Pagination ===")
    
    response = client.send_request("tools/call", {
        "name": "filter_nodes_by_attributes",
        "arguments": {
            "domain_name": "products",
            "filters": [
                {
                    "name": "status",
                    "value": "available",
                    "operator": "equals"
                }
            ],
            "page": 1,
            "size": 2
        }
    })
    
    if "result" in response:
        data = json.loads(response['result']['content'][0]['text'])
        print(f"✅ Page 1 of {data['total_pages']} (size=2)")
        print(f"  Total items: {data['total_count']}")
        print(f"  Items on this page: {len(data['nodes'])}")
    else:
        print(f"❌ Pagination test failed: {response}")

def main():
    """Run all enhanced query tests"""
    print("Starting enhanced MCP query tests...")
    print("=" * 50)
    
    client = MCPClient()
    
    try:
        # Set up test data
        composite_ids = setup_test_data(client)
        
        if composite_ids:
            # Run tests
            test_get_node_with_attributes(client, composite_ids[0])
            test_filter_by_single_attribute(client)
            test_filter_by_multiple_attributes(client)
            test_filter_with_operators(client)
            test_pagination(client)
        
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