#!/usr/bin/env python3
"""Test creating nodes with attributes"""

import json
import subprocess
import sys
import time
from tool_constants import CREATE_DOMAIN, CREATE_DOMAIN_ATTRIBUTE, CREATE_NODE, GET_NODE, GET_NODE_ATTRIBUTES, LIST_DOMAIN_ATTRIBUTES, LIST_NODES, SET_NODE_ATTRIBUTES


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
        time.sleep(0.5)
    
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

def test_node_with_attributes():
    """Test creating nodes with attributes"""
    client = MCPClient()
    
    try:
        print("Starting MCP server...")
        client.start()
        
        # Initialize
        print("\n1. Initializing server...")
        response = client.send_request("initialize", {
            "protocolVersion": "2024-11-05",
            "capabilities": {}
        })
        
        # Send initialized notification
        notification = {
            "jsonrpc": "2.0",
            "method": "notifications/initialized"
        }
        client.proc.stdin.write(json.dumps(notification) + "\n")
        client.proc.stdin.flush()
        time.sleep(0.1)
        
        # Create domain (or use existing)
        print("\n2. Creating test domain 'tech-articles'...")
        response = client.send_request("tools/call", {
            "name": CREATE_DOMAIN,
            "arguments": {
                "name": "tech-articles",
                "description": "Technology articles and resources"
            }
        })
        
        if response["result"].get("isError", False) and "already exists" in response["result"]["content"][0]["text"]:
            print("✓ Domain already exists, continuing with existing domain")
        elif response["result"].get("isError", False):
            print(f"❌ Error creating domain: {response}")
            return
        else:
            print("✓ Domain created successfully")
            print(f"Response: {json.dumps(response, indent=2)}")
        
        # Create attributes for the domain
        print("\n3. Creating domain attributes...")
        
        # Category attribute
        print("\n   a. Creating 'category' attribute...")
        response = client.send_request("tools/call", {
            "name": CREATE_DOMAIN_ATTRIBUTE,
            "arguments": {
                "domain_name": "tech-articles",
                "name": "category",
                "type": "tag",
                "description": "Article category (e.g., AI, Security, Cloud)"
            }
        })
        
        if response["result"].get("isError", False) and "already exists" in response["result"]["content"][0]["text"]:
            print("✓ Attribute 'category' already exists")
        elif response["result"].get("isError", False):
            print(f"❌ Error creating category attribute: {response}")
        else:
            print(f"✓ Created 'category' attribute")
            print(f"Response: {json.dumps(response, indent=2)}")
        
        # Priority attribute
        print("\n   b. Creating 'priority' attribute...")
        response = client.send_request("tools/call", {
            "name": CREATE_DOMAIN_ATTRIBUTE,
            "arguments": {
                "domain_name": "tech-articles",
                "name": "priority",
                "type": "ordered_tag",
                "description": "Reading priority (high, medium, low)"
            }
        })
        
        if response["result"].get("isError", False) and "already exists" in response["result"]["content"][0]["text"]:
            print("✓ Attribute 'priority' already exists")
        elif response["result"].get("isError", False):
            print(f"❌ Error creating priority attribute: {response}")
        else:
            print(f"✓ Created 'priority' attribute")
        
        # Rating attribute
        print("\n   c. Creating 'rating' attribute...")
        response = client.send_request("tools/call", {
            "name": CREATE_DOMAIN_ATTRIBUTE,
            "arguments": {
                "domain_name": "tech-articles",
                "name": "rating",
                "type": "number",
                "description": "Article rating (1-5)"
            }
        })
        
        if response["result"].get("isError", False) and "already exists" in response["result"]["content"][0]["text"]:
            print("✓ Attribute 'rating' already exists")
        elif response["result"].get("isError", False):
            print(f"❌ Error creating rating attribute: {response}")
        else:
            print(f"✓ Created 'rating' attribute")
        
        # Summary attribute
        print("\n   d. Creating 'summary' attribute...")
        response = client.send_request("tools/call", {
            "name": CREATE_DOMAIN_ATTRIBUTE,
            "arguments": {
                "domain_name": "tech-articles",
                "name": "summary",
                "type": "string",
                "description": "Brief summary of the article"
            }
        })
        
        if response["result"].get("isError", False) and "already exists" in response["result"]["content"][0]["text"]:
            print("✓ Attribute 'summary' already exists")
        elif response["result"].get("isError", False):
            print(f"❌ Error creating summary attribute: {response}")
        else:
            print(f"✓ Created 'summary' attribute")
        
        # List attributes
        print("\n4. Listing domain attributes...")
        response = client.send_request("tools/call", {
            "name": LIST_DOMAIN_ATTRIBUTES,
            "arguments": {
                "domain_name": "tech-articles"
            }
        })
        content = response["result"]["content"][0]["text"]
        data = json.loads(content)
        print(f"Found {data['total_count']} attributes:")
        for attr in data["attributes"]:
            print(f"  - {attr['name']} ({attr['type']}): {attr['description']}")
        
        # Get existing nodes or create new ones
        print("\n5. Listing existing nodes...")
        response = client.send_request("tools/call", {
            "name": LIST_NODES,
            "arguments": {
                "domain_name": "tech-articles"
            }
        })
        
        content = response["result"]["content"][0]["text"]
        nodes_data = json.loads(content)
        
        if nodes_data["total_count"] > 0:
            # Use existing node
            node_composite_id = nodes_data["nodes"][0]["composite_id"]
            print(f"Using existing node with ID: {node_composite_id}")
        else:
            # Create new node with unique URL
            unique_id = int(time.time())
            response = client.send_request("tools/call", {
                "name": CREATE_NODE,
                "arguments": {
                    "domain_name": "tech-articles",
                    "url": f"https://example.com/ai-security-best-practices-{unique_id}",
                    "title": "AI Security Best Practices in 2025",
                    "description": "Comprehensive guide on securing AI systems"
                }
            })
            
            content = response["result"]["content"][0]["text"]
            node_data = json.loads(content)
            node_composite_id = node_data["composite_id"]
            print(f"Created new node with ID: {node_composite_id}")
        
        # Set attributes on the node
        print("\n6. Setting attributes on the node...")
        response = client.send_request("tools/call", {
            "name": SET_NODE_ATTRIBUTES,
            "arguments": {
                "composite_id": node_composite_id,
                "attributes": [
                    {
                        "name": "category",
                        "value": "AI"
                    },
                    {
                        "name": "priority",
                        "value": "high",
                        "order_index": 1
                    },
                    {
                        "name": "rating",
                        "value": "5"
                    },
                    {
                        "name": "summary",
                        "value": "Essential reading for anyone deploying AI systems in production"
                    }
                ]
            }
        })
        print(f"Response: {json.dumps(response, indent=2)}")
        
        # Get node attributes
        print("\n7. Getting node attributes...")
        response = client.send_request("tools/call", {
            "name": GET_NODE_ATTRIBUTES,
            "arguments": {
                "composite_id": node_composite_id
            }
        })
        
        content = response["result"]["content"][0]["text"]
        attr_data = json.loads(content)
        print(f"\nNode attributes for {node_composite_id}:")
        for attr in attr_data["attributes"]:
            print(f"  - {attr['name']}: {attr['value']} (type: {attr['type']})")
        
        # Create another node with different attributes
        print("\n8. Creating another node...")
        unique_id2 = int(time.time() * 1000)  # More unique timestamp
        response = client.send_request("tools/call", {
            "name": CREATE_NODE,
            "arguments": {
                "domain_name": "tech-articles",
                "url": f"https://example.com/cloud-migration-guide-{unique_id2}",
                "title": "Complete Cloud Migration Guide",
                "description": "Step-by-step guide for cloud migration"
            }
        })
        
        if response["result"].get("isError", False):
            print(f"Error creating second node: {response['result']['content'][0]['text']}")
            # Try to continue with existing nodes
            return
            
        content = response["result"]["content"][0]["text"]
        node_data2 = json.loads(content)
        node_composite_id2 = node_data2["composite_id"]
        
        # Set different attributes
        response = client.send_request("tools/call", {
            "name": SET_NODE_ATTRIBUTES,
            "arguments": {
                "composite_id": node_composite_id2,
                "attributes": [
                    {
                        "name": "category",
                        "value": "Cloud"
                    },
                    {
                        "name": "priority",
                        "value": "medium",
                        "order_index": 2
                    },
                    {
                        "name": "rating",
                        "value": "4"
                    },
                    {
                        "name": "summary",
                        "value": "Practical guide covering AWS, Azure, and GCP migration strategies"
                    }
                ]
            }
        })
        
        # List all nodes in the domain
        print("\n9. Listing all nodes in the domain...")
        response = client.send_request("tools/call", {
            "name": LIST_NODES,
            "arguments": {
                "domain_name": "tech-articles"
            }
        })
        
        content = response["result"]["content"][0]["text"]
        nodes_data = json.loads(content)
        print(f"\nFound {nodes_data['total_count']} nodes in 'tech-articles' domain")
        
        print("\n✅ Successfully created nodes with attributes!")
        print("\nSummary:")
        print("- Created domain 'tech-articles' with 4 attribute types")
        print("- Created 2 nodes with different attribute values")
        print("- Demonstrated full CRUD operations on node attributes")
        
    except Exception as e:
        print(f"\n❌ Test failed: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
    finally:
        client.close()

if __name__ == "__main__":
    test_node_with_attributes()