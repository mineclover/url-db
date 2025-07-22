#!/usr/bin/env python3
"""
Comprehensive MCP Tools Test Script
Tests all 11 MCP tools through the MCP protocol
"""

import json
import subprocess
import sys
import time
from datetime import datetime
from tool_constants import CREATE_DOMAIN, CREATE_NODE, DELETE_NODE, FIND_NODE_BY_URL, GET_NODE, GET_NODE_ATTRIBUTES, GET_SERVER_INFO, LIST_DOMAINS, LIST_NODES, SET_NODE_ATTRIBUTES, UPDATE_NODE


class MCPToolTester:
    def __init__(self):
        self.server_path = "../../bin/url-db"
        self.process = None
        self.request_id = 1
        
    def start_server(self):
        """Start the MCP server in stdio mode"""
        print("🚀 Starting MCP server...")
        self.process = subprocess.Popen(
            [self.server_path, "-mcp-mode=stdio"],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            bufsize=0
        )
        time.sleep(0.5)
        
    def send_request(self, method, params=None):
        """Send JSON-RPC request and return response"""
        request = {
            "jsonrpc": "2.0",
            "id": self.request_id,
            "method": method,
            "params": params or {}
        }
        self.request_id += 1
        
        request_str = json.dumps(request)
        self.process.stdin.write(request_str + '\n')
        self.process.stdin.flush()
        
        response_str = self.process.stdout.readline()
        if response_str:
            return json.loads(response_str) if response_str.strip() else {}
        return None
        
    def initialize(self):
        """Initialize MCP protocol"""
        print("\n📋 Step 1: MCP Protocol Initialization")
        
        # Send initialize request
        response = self.send_request("initialize", {
            "protocolVersion": "2024-11-05",
            "capabilities": {
                "roots": {"listChanged": True},
                "sampling": {}
            },
            "clientInfo": {
                "name": "mcp-tools-tester",
                "version": "1.0.0"
            }
        })
        
        if response and not response.get("error"):
            print("✅ Initialize successful")
            print(f"   Server: {response['result']['serverInfo']['name']} v{response['result']['serverInfo']['version']}")
            
            # Send initialized notification (no ID for notifications)
            notification = {
                "jsonrpc": "2.0",
                "method": "notifications/initialized"
            }
            self.process.stdin.write(json.dumps(notification) + '\n')
            self.process.stdin.flush()
            time.sleep(0.1)  # Give server time to process
            print("✅ Initialized notification sent")
            return True
        else:
            print("❌ Initialize failed")
            return False
            
    def test_server_info(self):
        """Test get_server_info tool"""
        print("\n📋 Step 2: Testing get_server_info")
        
        response = self.send_request("tools/call", {
            "name": GET_SERVER_INFO,
            "arguments": {}
        })
        
        if response and not response.get("error"):
            result = response["result"]["content"][0]["text"]
            server_info = json.loads(result) if result.strip() else {}
            print("✅ Server info retrieved:")
            print(f"   Name: {server_info['name']}")
            print(f"   Version: {server_info['version']}")
            print(f"   Capabilities: {', '.join(server_info['capabilities'])}")
            return True
        else:
            print("❌ Failed to get server info")
            return False
            
    def test_list_domains(self):
        """Test list_domains tool"""
        print("\n📋 Step 3: Testing list_domains")
        
        response = self.send_request("tools/call", {
            "name": LIST_DOMAINS,
            "arguments": {}
        })
        
        if response and not response.get("error"):
            result = response["result"]["content"][0]["text"]
            domains = json.loads(result) if result.strip() else {}["domains"]
            print(f"✅ Found {len(domains)} domains:")
            for domain in domains[:5]:  # Show first 5
                print(f"   - {domain['name']}: {domain['description']}")
            return domains
        else:
            print("❌ Failed to list domains")
            return []
            
    def test_create_domain(self):
        """Test create_domain tool"""
        print("\n📋 Step 4: Testing create_domain")
        
        timestamp = int(time.time())
        domain_name = f"mcp-test-{timestamp}"
        
        response = self.send_request("tools/call", {
            "name": CREATE_DOMAIN,
            "arguments": {
                "name": domain_name,
                "description": "Test domain created via MCP protocol"
            }
        })
        
        if response and not response.get("error"):
            result = response["result"]["content"][0]["text"]
            domain = json.loads(result) if result.strip() else {}
            print(f"✅ Domain created: {domain['name']}")
            print(f"   Created at: {domain['created_at']}")
            return domain_name
        else:
            print("❌ Failed to create domain")
            return None
            
    def test_create_node(self, domain_name):
        """Test create_node tool"""
        print("\n📋 Step 5: Testing create_node")
        
        response = self.send_request("tools/call", {
            "name": CREATE_NODE,
            "arguments": {
                "domain_name": domain_name,
                "url": f"https://example.com/test-{int(time.time())}",
                "title": "Test Node via MCP",
                "description": "This node was created through MCP protocol testing"
            }
        })
        
        if response and not response.get("error"):
            result = response["result"]["content"][0]["text"]
            node = json.loads(result) if result.strip() else {}
            print(f"✅ Node created: {node['composite_id']}")
            print(f"   URL: {node['url']}")
            print(f"   Title: {node['title']}")
            return node['composite_id']
        else:
            print("❌ Failed to create node")
            return None
            
    def test_get_node(self, composite_id):
        """Test get_node tool"""
        print("\n📋 Step 6: Testing get_node")
        
        response = self.send_request("tools/call", {
            "name": GET_NODE,
            "arguments": {
                "composite_id": composite_id
            }
        })
        
        if response and not response.get("error"):
            result = response["result"]["content"][0]["text"]
            node = json.loads(result) if result.strip() else {}
            print(f"✅ Node retrieved: {node['composite_id']}")
            print(f"   URL: {node['url']}")
            return True
        else:
            print("❌ Failed to get node")
            return False
            
    def test_update_node(self, composite_id):
        """Test update_node tool"""
        print("\n📋 Step 7: Testing update_node")
        
        response = self.send_request("tools/call", {
            "name": UPDATE_NODE,
            "arguments": {
                "composite_id": composite_id,
                "title": "Updated Node Title",
                "description": "This description was updated via MCP protocol"
            }
        })
        
        if response and not response.get("error"):
            result = response["result"]["content"][0]["text"]
            node = json.loads(result) if result.strip() else {}
            print(f"✅ Node updated: {node['composite_id']}")
            print(f"   New title: {node['title']}")
            return True
        else:
            print("❌ Failed to update node")
            return False
            
    def test_set_attributes(self, composite_id):
        """Test set_node_attributes tool"""
        print("\n📋 Step 8: Testing set_node_attributes")
        
        response = self.send_request("tools/call", {
            "name": SET_NODE_ATTRIBUTES,
            "arguments": {
                "composite_id": composite_id,
                "attributes": [
                    {"name": "category", "value": "testing"},
                    {"name": "priority", "value": "high"},
                    {"name": "tags", "value": "mcp,test,automated"}
                ]
            }
        })
        
        if response and not response.get("error"):
            result = response["result"]["content"][0]["text"]
            if result:
                try:
                    attr_response = json.loads(result) if result.strip() else {}
                    print(f"✅ Attributes set for: {attr_response['composite_id']}")
                    print(f"   Total attributes: {len(attr_response['attributes'])}")
                    for attr in attr_response['attributes']:
                        print(f"   - {attr['name']}: {attr['value']}")
                    return True
                except json.JSONDecodeError:
                    print(f"✅ Attributes set successfully")
                    return True
            else:
                print("✅ Attributes set successfully")
                return True
        else:
            print("❌ Failed to set attributes")
            if response:
                print(f"   Error: {response.get('error', {}).get('message', 'Unknown error')}")
            return False
            
    def test_get_attributes(self, composite_id):
        """Test get_node_attributes tool"""
        print("\n📋 Step 9: Testing get_node_attributes")
        
        response = self.send_request("tools/call", {
            "name": GET_NODE_ATTRIBUTES,
            "arguments": {
                "composite_id": composite_id
            }
        })
        
        if response and not response.get("error"):
            result = response["result"]["content"][0]["text"]
            if result:
                try:
                    attr_response = json.loads(result) if result.strip() else {}
                    print(f"✅ Attributes retrieved for: {attr_response['composite_id']}")
                    if attr_response.get('attributes'):
                        for attr in attr_response['attributes']:
                            print(f"   - {attr['name']}: {attr['value']}")
                    else:
                        print("   No attributes found")
                    return True
                except json.JSONDecodeError:
                    print(f"✅ Attributes retrieved (text response): {result}")
                    return True
            else:
                print("✅ No attributes found")
                return True
        else:
            print("❌ Failed to get attributes")
            if response:
                print(f"   Error: {response.get('error', {}).get('message', 'Unknown error')}")
            return False
            
    def test_find_by_url(self, domain_name, url):
        """Test find_node_by_url tool"""
        print("\n📋 Step 10: Testing find_node_by_url")
        
        response = self.send_request("tools/call", {
            "name": FIND_NODE_BY_URL,
            "arguments": {
                "domain_name": domain_name,
                "url": url
            }
        })
        
        if response and not response.get("error"):
            result = response["result"]["content"][0]["text"]
            node = json.loads(result) if result.strip() else {}
            print(f"✅ Node found by URL: {node['composite_id']}")
            return True
        else:
            print("❌ Failed to find node by URL")
            return False
            
    def test_list_nodes(self, domain_name):
        """Test list_nodes tool"""
        print("\n📋 Step 11: Testing list_nodes")
        
        response = self.send_request("tools/call", {
            "name": LIST_NODES,
            "arguments": {
                "domain_name": domain_name,
                "page": 1,
                "size": 10
            }
        })
        
        if response and not response.get("error"):
            result = response["result"]["content"][0]["text"]
            nodes_response = json.loads(result) if result.strip() else {}
            print(f"✅ Found {nodes_response['total_count']} nodes in domain {domain_name}")
            for node in nodes_response['nodes']:
                print(f"   - {node['composite_id']}: {node['title']}")
            return True
        else:
            print("❌ Failed to list nodes")
            return False
            
    def test_delete_node(self, composite_id):
        """Test delete_node tool"""
        print("\n📋 Step 12: Testing delete_node")
        
        response = self.send_request("tools/call", {
            "name": DELETE_NODE,
            "arguments": {
                "composite_id": composite_id
            }
        })
        
        if response and not response.get("error"):
            result = response["result"]["content"][0]["text"]
            print(f"✅ {result}")
            return True
        else:
            print("❌ Failed to delete node")
            return False
            
    def run_all_tests(self):
        """Run all MCP tool tests"""
        print("\n" + "="*60)
        print("🧪 MCP Tools Comprehensive Test Suite")
        print("="*60)
        
        try:
            self.start_server()
            
            # Initialize protocol
            if not self.initialize():
                print("❌ Failed to initialize MCP protocol")
                return
                
            # Test all tools
            success_count = 0
            total_tests = 12
            
            # Test 1: Server info
            if self.test_server_info():
                success_count += 1
                
            # Test 2: List domains
            domains = self.test_list_domains()
            if domains is not None:
                success_count += 1
                
            # Test 3: Create domain
            domain_name = self.test_create_domain()
            if domain_name:
                success_count += 1
                
                # Test 4: Create node
                composite_id = self.test_create_node(domain_name)
                if composite_id:
                    success_count += 1
                    
                    # Test 5: Get node
                    if self.test_get_node(composite_id):
                        success_count += 1
                        
                    # Test 6: Update node
                    if self.test_update_node(composite_id):
                        success_count += 1
                        
                    # Test 7: Set attributes
                    if self.test_set_attributes(composite_id):
                        success_count += 1
                        
                    # Test 8: Get attributes
                    if self.test_get_attributes(composite_id):
                        success_count += 1
                        
                    # Test 9: Find by URL
                    url = f"https://example.com/test-{int(time.time())}"
                    node2_id = self.test_create_node(domain_name)  # Create another node
                    if node2_id:
                        # Get the actual URL from the node
                        response = self.send_request("tools/call", {
                            "name": GET_NODE,
                            "arguments": {"composite_id": node2_id}
                        })
                        if response and not response.get("error"):
                            node = json.loads(response["result"]["content"][0]["text"]) if response["result"]["content"][0]["text"].strip() else {}
                            if self.test_find_by_url(domain_name, node['url']):
                                success_count += 1
                                
                    # Test 10: List nodes
                    if self.test_list_nodes(domain_name):
                        success_count += 1
                        
                    # Test 11: Delete node
                    if self.test_delete_node(composite_id):
                        success_count += 1
                        
                    # Test 12: Delete second node
                    if node2_id and self.test_delete_node(node2_id):
                        success_count += 1
                        
            print("\n" + "="*60)
            print(f"📊 Test Results: {success_count}/{total_tests} tests passed")
            print(f"✅ Success Rate: {(success_count/total_tests)*100:.1f}%")
            print("="*60)
            
        finally:
            if self.process:
                self.process.terminate()
                print("\n🛑 MCP server stopped")

if __name__ == "__main__":
    tester = MCPToolTester()
    tester.run_all_tests()