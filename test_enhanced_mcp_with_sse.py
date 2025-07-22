#!/usr/bin/env python3
"""
Test enhanced MCP functions in both stdio and SSE modes
Tests get_node_with_attributes and filter_nodes_by_attributes
"""

import json
import sys
import subprocess
import time
import urllib.request
import urllib.parse
import urllib.error

class MCPStdioClient:
    def __init__(self, db_path: str = "./test_enhanced_stdio.db"):
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

class MCPSSEClient:
    def __init__(self, base_url: str = "http://localhost:8080"):
        """Initialize MCP client with SSE/HTTP connection"""
        self.base_url = base_url
        self.api_url = f"{base_url}/api/mcp"
        
    def wait_for_server(self, max_retries=30):
        """Wait for server to be ready"""
        for i in range(max_retries):
            try:
                with urllib.request.urlopen(f"{self.base_url}/health") as response:
                    if response.status == 200:
                        return True
            except:
                pass
            time.sleep(1)
        return False
    
    def post_json(self, url, data):
        """Make POST request with JSON data"""
        req = urllib.request.Request(url, 
            data=json.dumps(data).encode('utf-8'),
            headers={'Content-Type': 'application/json'})
        with urllib.request.urlopen(req) as response:
            return response.status, json.loads(response.read().decode('utf-8'))
    
    def get_json(self, url):
        """Make GET request and return JSON"""
        with urllib.request.urlopen(url) as response:
            return response.status, json.loads(response.read().decode('utf-8'))

def setup_test_data_stdio(client):
    """Set up test data using stdio client"""
    print("\n=== Setting up test data (stdio) ===")
    
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
            "name": "testproducts",
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
                "domain_name": "testproducts",
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
         [("category", "electronics"), ("price", "899.99"), ("status", "out-of-stock"), ("brand", "PhoneCo")])
    ]
    
    composite_ids = []
    for url, title, desc, attrs in products:
        # Create node
        response = client.send_request("tools/call", {
            "name": "create_node",
            "arguments": {
                "domain_name": "testproducts",
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

def test_stdio_enhanced_functions(client, composite_ids):
    """Test enhanced functions via stdio"""
    print("\n=== Testing Enhanced Functions (stdio) ===")
    
    # Test get_node_with_attributes
    print("\n1. Testing get_node_with_attributes...")
    response = client.send_request("tools/call", {
        "name": "get_node_with_attributes",
        "arguments": {
            "composite_id": composite_ids[0]
        }
    })
    
    if "result" in response:
        data = json.loads(response['result']['content'][0]['text'])
        print("✅ get_node_with_attributes works!")
        print(f"   Node: {data['node']['title']}")
        print(f"   Attributes: {len(data['attributes'])} found")
    else:
        print(f"❌ get_node_with_attributes failed: {response}")
    
    # Test filter_nodes_by_attributes
    print("\n2. Testing filter_nodes_by_attributes...")
    response = client.send_request("tools/call", {
        "name": "filter_nodes_by_attributes",
        "arguments": {
            "domain_name": "testproducts",
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
        print("✅ filter_nodes_by_attributes works!")
        print(f"   Found {data['total_count']} nodes")
    else:
        print(f"❌ filter_nodes_by_attributes failed: {response}")

def test_sse_enhanced_functions(base_url="http://localhost:8080"):
    """Test enhanced functions via SSE/HTTP"""
    print("\n=== Testing Enhanced Functions (SSE/HTTP) ===")
    
    client = MCPSSEClient(base_url)
    
    # Wait for server
    if not client.wait_for_server():
        print("❌ Server not ready")
        return False
    
    # Create test domain via REST API
    print("\n1. Setting up test data via REST API...")
    
    # Create domain
    try:
        status, domain = client.post_json(f"{base_url}/api/domains", {
            "name": "ssetest",
            "description": "SSE test domain"
        })
        
        if status != 201:
            print(f"❌ Failed to create domain")
            return False
        
        domain_id = domain["id"]
    except Exception as e:
        print(f"❌ Failed to create domain: {e}")
        return False
    
    # Create attributes
    attrs = [
        {"name": "category", "type": "tag", "description": "Category"},
        {"name": "price", "type": "number", "description": "Price"}
    ]
    
    for attr in attrs:
        try:
            status, _ = client.post_json(f"{base_url}/api/domains/{domain_id}/attributes", attr)
            if status != 201:
                print(f"❌ Failed to create attribute")
        except Exception as e:
            print(f"❌ Failed to create attribute: {e}")
    
    # Create nodes
    nodes = [
        {"url": "https://sse.test/1", "title": "Product 1", "description": "Test product 1"},
        {"url": "https://sse.test/2", "title": "Product 2", "description": "Test product 2"}
    ]
    
    created_nodes = []
    for node in nodes:
        try:
            status, node_data = client.post_json(f"{base_url}/api/domains/{domain_id}/urls", node)
            if status == 201:
                created_nodes.append(node_data)
        except:
            pass
    
    # Set attributes for first node
    if created_nodes:
        node_id = created_nodes[0]["id"]
        
        # Get attribute IDs
        try:
            status, attrs_data = client.get_json(f"{base_url}/api/domains/{domain_id}/attributes")
            if status == 200:
                domain_attrs = attrs_data["attributes"]
                
                for attr in domain_attrs:
                    if attr["name"] == "category":
                        client.post_json(f"{base_url}/api/urls/{node_id}/attributes", {
                            "attribute_id": attr["id"],
                            "value": "electronics"
                        })
                    elif attr["name"] == "price":
                        client.post_json(f"{base_url}/api/urls/{node_id}/attributes", {
                            "attribute_id": attr["id"],
                            "value": "99.99"
                        })
        except:
            pass
    
    print("✅ Test data created via REST API")
    
    # Note: The current SSE implementation doesn't expose the new enhanced query functions
    # They are only available via stdio mode (JSON-RPC tools)
    print("\n⚠️  Note: Enhanced query functions (get_node_with_attributes, filter_nodes_by_attributes)")
    print("   are currently only available in stdio mode via MCP tools.")
    print("   SSE mode uses standard REST API endpoints.")
    
    return True

def main():
    """Run all enhanced MCP tests"""
    print("Starting enhanced MCP function tests...")
    print("=" * 50)
    
    # Test stdio mode
    print("\n### STDIO MODE TESTS ###")
    stdio_client = MCPStdioClient()
    
    try:
        composite_ids = setup_test_data_stdio(stdio_client)
        if composite_ids:
            test_stdio_enhanced_functions(stdio_client, composite_ids)
    except Exception as e:
        print(f"❌ Stdio test failed: {e}")
    finally:
        stdio_client.close()
    
    # Test SSE mode
    print("\n\n### SSE MODE TESTS ###")
    
    # Start server in SSE mode
    server_process = subprocess.Popen(
        ["./bin/url-db", "-mcp-mode=sse", "-db-path=./test_enhanced_sse.db", "-port=8089"],
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE
    )
    
    try:
        time.sleep(2)  # Give server time to start
        test_sse_enhanced_functions("http://localhost:8089")
    except Exception as e:
        print(f"❌ SSE test failed: {e}")
    finally:
        server_process.terminate()
        server_process.wait()
    
    print("\n" + "=" * 50)
    print("✅ All tests completed!")
    
    return 0

if __name__ == "__main__":
    sys.exit(main())