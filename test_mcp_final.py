#!/usr/bin/env python3
"""
Final MCP Tools Test - Safe version with better error handling
"""

import json
import subprocess
import sys
import time
from datetime import datetime

class MCPFinalTester:
    def __init__(self):
        self.server_path = "./cmd/server/url-db"
        self.process = None
        self.request_id = 1
        self.test_results = []
        
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
            try:
                return json.loads(response_str)
            except json.JSONDecodeError:
                return {"error": {"message": "Invalid JSON response"}}
        return None
        
    def safe_parse_result(self, response):
        """Safely parse tool response"""
        if not response or response.get("error"):
            return None, response.get("error", {}).get("message", "Unknown error") if response else "No response"
            
        try:
            result_text = response["result"]["content"][0]["text"]
            if result_text:
                # Try to parse as JSON
                try:
                    return json.loads(result_text), None
                except json.JSONDecodeError:
                    # Return as text if not JSON
                    return result_text, None
            else:
                return "", None
        except (KeyError, IndexError, TypeError) as e:
            return None, f"Response parsing error: {str(e)}"
            
    def run_test(self, test_name, test_func):
        """Run a single test and record results"""
        print(f"\n📋 {test_name}")
        try:
            success = test_func()
            self.test_results.append((test_name, success))
            return success
        except Exception as e:
            print(f"❌ Test failed with exception: {str(e)}")
            self.test_results.append((test_name, False))
            return False
            
    def test_protocol_init(self):
        """Test MCP protocol initialization"""
        # Send initialize
        response = self.send_request("initialize", {
            "protocolVersion": "2024-11-05",
            "capabilities": {
                "roots": {"listChanged": True},
                "sampling": {}
            },
            "clientInfo": {
                "name": "mcp-final-tester",
                "version": "1.0.0"
            }
        })
        
        if response and not response.get("error"):
            print("✅ Initialize successful")
            print(f"   Server: {response['result']['serverInfo']['name']} v{response['result']['serverInfo']['version']}")
            
            # Send initialized notification
            notification = {
                "jsonrpc": "2.0",
                "method": "notifications/initialized"
            }
            self.process.stdin.write(json.dumps(notification) + '\n')
            self.process.stdin.flush()
            time.sleep(0.1)
            print("✅ Initialized notification sent")
            return True
        return False
        
    def test_tools_list(self):
        """Test tools/list"""
        response = self.send_request("tools/list")
        if response and not response.get("error"):
            tools = response["result"]["tools"]
            print(f"✅ Found {len(tools)} tools")
            # Show first 3 tools
            for tool in tools[:3]:
                print(f"   - {tool['name']}: {tool['description']}")
            return True
        return False
        
    def test_server_info(self):
        """Test get_mcp_server_info"""
        response = self.send_request("tools/call", {
            "name": "get_mcp_server_info",
            "arguments": {}
        })
        
        result, error = self.safe_parse_result(response)
        if error:
            print(f"❌ Error: {error}")
            return False
            
        print("✅ Server info retrieved")
        if isinstance(result, dict):
            print(f"   Name: {result.get('name', 'N/A')}")
            print(f"   Version: {result.get('version', 'N/A')}")
        return True
        
    def test_list_domains(self):
        """Test list_mcp_domains"""
        response = self.send_request("tools/call", {
            "name": "list_mcp_domains",
            "arguments": {}
        })
        
        result, error = self.safe_parse_result(response)
        if error:
            print(f"❌ Error: {error}")
            return False
            
        if isinstance(result, dict) and "domains" in result:
            domains = result["domains"]
            print(f"✅ Found {len(domains)} domains")
            for domain in domains[:3]:
                print(f"   - {domain['name']}")
            return True
        return False
        
    def test_create_domain(self):
        """Test domain creation"""
        timestamp = int(time.time())
        domain_name = f"final-test-{timestamp}"
        
        response = self.send_request("tools/call", {
            "name": "create_mcp_domain",
            "arguments": {
                "name": domain_name,
                "description": "Final test domain"
            }
        })
        
        result, error = self.safe_parse_result(response)
        if error:
            print(f"❌ Error: {error}")
            return False
            
        print(f"✅ Domain created: {domain_name}")
        return domain_name
        
    def test_create_and_manage_node(self, domain_name):
        """Test node creation and management"""
        # Create node
        response = self.send_request("tools/call", {
            "name": "create_mcp_node",
            "arguments": {
                "domain_name": domain_name,
                "url": f"https://example.com/final-{int(time.time())}",
                "title": "Final Test Node",
                "description": "Node for final testing"
            }
        })
        
        result, error = self.safe_parse_result(response)
        if error:
            print(f"❌ Create node error: {error}")
            return False
            
        if isinstance(result, dict):
            composite_id = result.get("composite_id")
            print(f"✅ Node created: {composite_id}")
            
            # Get node
            response = self.send_request("tools/call", {
                "name": "get_mcp_node",
                "arguments": {"composite_id": composite_id}
            })
            
            result, error = self.safe_parse_result(response)
            if not error:
                print(f"✅ Node retrieved successfully")
                
            # Update node
            response = self.send_request("tools/call", {
                "name": "update_mcp_node",
                "arguments": {
                    "composite_id": composite_id,
                    "title": "Updated Final Node"
                }
            })
            
            result, error = self.safe_parse_result(response)
            if not error:
                print(f"✅ Node updated successfully")
                
            # Delete node
            response = self.send_request("tools/call", {
                "name": "delete_mcp_node",
                "arguments": {"composite_id": composite_id}
            })
            
            result, error = self.safe_parse_result(response)
            if not error:
                print(f"✅ Node deleted successfully")
                
            return True
            
        return False
        
    def test_resources(self):
        """Test resource system"""
        # List resources
        response = self.send_request("resources/list")
        
        if response and not response.get("error"):
            resources = response["result"]["resources"]
            print(f"✅ Found {len(resources)} resources")
            
            # Read server info resource
            response = self.send_request("resources/read", {
                "uri": "mcp://server/info"
            })
            
            if response and not response.get("error"):
                print("✅ Server info resource read successfully")
                return True
                
        return False
        
    def run_all_tests(self):
        """Run all tests"""
        print("\n" + "="*60)
        print("🧪 MCP Final Integration Test")
        print("="*60)
        
        try:
            self.start_server()
            
            # Run tests
            self.run_test("1. Protocol Initialization", self.test_protocol_init)
            self.run_test("2. Tools Discovery", self.test_tools_list)
            self.run_test("3. Server Information", self.test_server_info)
            self.run_test("4. List Domains", self.test_list_domains)
            
            domain_name = self.run_test("5. Create Domain", self.test_create_domain)
            if domain_name and isinstance(domain_name, str):
                self.run_test("6. Node Operations", lambda: self.test_create_and_manage_node(domain_name))
                
            self.run_test("7. Resource System", self.test_resources)
            
            # Calculate results
            passed = sum(1 for _, success in self.test_results if success)
            total = len(self.test_results)
            
            print("\n" + "="*60)
            print(f"📊 Final Results: {passed}/{total} tests passed")
            print(f"✅ Success Rate: {(passed/total)*100:.1f}%")
            
            if passed == total:
                print("🎉 All tests passed! MCP server is working correctly.")
            elif passed >= total * 0.8:
                print("👍 Most tests passed. MCP server is mostly functional.")
            else:
                print("⚠️  Several tests failed. Check implementation.")
                
            print("="*60)
            
        finally:
            if self.process:
                self.process.terminate()
                print("\n🛑 MCP server stopped")

if __name__ == "__main__":
    tester = MCPFinalTester()
    tester.run_all_tests()