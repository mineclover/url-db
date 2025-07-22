#!/usr/bin/env python3
"""
Test script for Dependency Relationship Creation and Querying scenario
"""
import json
import subprocess
import sys
import os
import time

class MCPClient:
    def __init__(self, server_path):
        self.server_path = server_path
        self.request_id = 1
        self.process = None
        self.initialized = False
        
    def start_server(self):
        """Start the MCP server process"""
        if self.process:
            return
            
        self.process = subprocess.Popen(
            [self.server_path, "-mcp-mode=stdio"],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            bufsize=0
        )
        
    def initialize(self):
        """Initialize the MCP server"""
        self.start_server()
        
        # Send initialization request
        init_request = {
            "jsonrpc": "2.0",
            "id": self.request_id,
            "method": "initialize",
            "params": {
                "protocolVersion": "0.1.0",
                "capabilities": {
                    "roots": {
                        "listChanged": True
                    }
                },
                "clientInfo": {
                    "name": "test-client",
                    "version": "1.0.0"
                }
            }
        }
        
        self.request_id += 1
        request_str = json.dumps(init_request) + "\n"
        
        try:
            self.process.stdin.write(request_str)
            self.process.stdin.flush()
            
            # Wait for response
            response_line = self.process.stdout.readline()
            response = json.loads(response_line.strip())
            print(f"Initialization response: {response}")
            
            # Send initialized notification
            initialized_notification = {
                "jsonrpc": "2.0",
                "method": "notifications/initialized"
            }
            
            notification_str = json.dumps(initialized_notification) + "\n"
            self.process.stdin.write(notification_str)
            self.process.stdin.flush()
            
            self.initialized = True
            return response
            
        except Exception as e:
            print(f"Initialization failed: {e}")
            return {"error": {"code": -1, "message": str(e)}}
        
    def send_request(self, method, params=None):
        """Send an MCP request to the URL-DB server"""
        if not self.initialized:
            init_result = self.initialize()
            if "error" in init_result:
                return init_result
        
        request = {
            "jsonrpc": "2.0",
            "id": self.request_id,
            "method": method
        }
        if params:
            request["params"] = params
            
        self.request_id += 1
        
        try:
            request_str = json.dumps(request) + "\n"
            self.process.stdin.write(request_str)
            self.process.stdin.flush()
            
            # Wait for response
            response_line = self.process.stdout.readline()
            response = json.loads(response_line.strip())
            return response
            
        except Exception as e:
            print(f"Request failed: {e}")
            return {"error": {"code": -1, "message": str(e)}}
    
    def call_tool(self, tool_name, arguments=None):
        """Call an MCP tool"""
        response = self.send_request("tools/call", {
            "name": tool_name,
            "arguments": arguments or {}
        })
        
        # Parse MCP tool response format
        if "result" in response and "content" in response["result"]:
            content = response["result"]["content"]
            if content and len(content) > 0 and "text" in content[0]:
                text_content = content[0]["text"]
                
                # Check if it's an error message
                if text_content.startswith("Error"):
                    return {"error": {"message": text_content}}
                
                # Try to parse as JSON
                try:
                    parsed_content = json.loads(text_content)
                    return {"result": parsed_content}
                except json.JSONDecodeError:
                    return {"result": {"text": text_content}}
        
        return response
    
    def close(self):
        """Close the MCP connection"""
        if self.process:
            self.process.stdin.close()
            self.process.wait()
            self.process = None

def main():
    print("=== Dependency Relationship Creation and Querying Scenario ===\n")
    
    # Initialize MCP client
    server_path = "/Users/junwoobang/mcp/url-db/bin/url-db"
    if not os.path.exists(server_path):
        print("Building URL-DB server...")
        build_result = subprocess.run(["make", "build"], cwd="/Users/junwoobang/mcp/url-db", capture_output=True, text=True)
        if build_result.returncode != 0:
            print(f"Build failed: {build_result.stderr}")
            return
    
    client = MCPClient(server_path)
    
    # Step 1: Setup - Create microservices domain if not exists
    print("Step 1: Setting up microservices domain...")
    domains_response = client.call_tool("list_domains")
    print(f"Current domains: {domains_response}")
    
    if "result" in domains_response:
        domains = [d["name"] for d in domains_response["result"].get("domains", [])]
        if "microservices" not in domains:
            create_domain = client.call_tool("create_domain", {
                "name": "microservices",
                "description": "Microservices architecture components"
            })
            print(f"Created domain: {create_domain}")
    
    # Step 2: Create/ensure nodes exist (api-gateway, user-service, payment-service)
    print("\nStep 2: Creating/ensuring nodes exist...")
    
    nodes = [
        {"name": "api-gateway", "url": "https://api.example.com/gateway"},
        {"name": "user-service", "url": "https://api.example.com/user-service"},
        {"name": "payment-service", "url": "https://api.example.com/payment-service"}
    ]
    
    node_ids = {}
    
    for node in nodes:
        # Try to find existing node first
        find_response = client.call_tool("find_node_by_url", {
            "domain_name": "microservices",
            "url": node["url"]
        })
        
        print(f"Find response for {node['name']}: {find_response}")
        
        if ("result" in find_response and 
            find_response["result"] and 
            find_response["result"] != "null" and
            isinstance(find_response["result"], dict) and
            "composite_id" in find_response["result"]):
            node_ids[node["name"]] = find_response["result"]["composite_id"]
            print(f"Found existing {node['name']}: {find_response['result']['composite_id']}")
        else:
            print(f"Node not found for {node['name']}, creating new node...")
            # Create new node
            create_response = client.call_tool("create_node", {
                "domain_name": "microservices",
                "url": node["url"],
                "title": node["name"],
                "description": f"{node['name']} component"
            })
            
            print(f"Create response for {node['name']}: {create_response}")
            
            if ("result" in create_response and 
                create_response["result"] and 
                isinstance(create_response["result"], dict) and
                "composite_id" in create_response["result"]):
                node_ids[node["name"]] = create_response["result"]["composite_id"]
                print(f"Created {node['name']}: {create_response['result']['composite_id']}")
            else:
                print(f"Failed to create {node['name']}: {create_response}")
                return
    
    print(f"Node IDs: {node_ids}")
    
    # Step 3: Create hard dependency: api-gateway depends on user-service
    print("\nStep 3: Creating hard dependency (api-gateway -> user-service)...")
    hard_dep_response = client.call_tool("create_dependency", {
        "source_composite_id": node_ids["api-gateway"],
        "target_composite_id": node_ids["user-service"],
        "dependency_type": "hard",
        "cascade_delete": True,
        "cascade_update": True,
        "description": "API Gateway requires User Service for authentication"
    })
    print(f"Hard dependency result: {hard_dep_response}")
    
    # Step 4: Create soft dependency: user-service depends on payment-service
    print("\nStep 4: Creating soft dependency (user-service -> payment-service)...")
    soft_dep_response = client.call_tool("create_dependency", {
        "source_composite_id": node_ids["user-service"],
        "target_composite_id": node_ids["payment-service"],
        "dependency_type": "soft",
        "cascade_delete": False,
        "cascade_update": True,
        "description": "User Service integrates with Payment Service for premium features"
    })
    print(f"Soft dependency result: {soft_dep_response}")
    
    # Step 5: Create reference dependency: api-gateway references payment-service
    print("\nStep 5: Creating reference dependency (api-gateway references payment-service)...")
    ref_dep_response = client.call_tool("create_dependency", {
        "source_composite_id": node_ids["api-gateway"],
        "target_composite_id": node_ids["payment-service"],
        "dependency_type": "reference",
        "cascade_delete": False,
        "cascade_update": False,
        "description": "API Gateway references Payment Service for direct payment endpoints",
        "metadata": {
            "relationship": "external API",
            "description": "Payment processing endpoint"
        }
    })
    print(f"Reference dependency result: {ref_dep_response}")
    
    # Step 6: List dependencies for api-gateway (should show 2)
    print("\nStep 6: Listing dependencies for api-gateway...")
    api_deps_response = client.call_tool("list_node_dependencies", {
        "composite_id": node_ids["api-gateway"]
    })
    print(f"API Gateway dependencies: {api_deps_response}")
    
    if "result" in api_deps_response:
        deps = api_deps_response["result"].get("dependencies", [])
        print(f"API Gateway has {len(deps)} dependencies:")
        for dep in deps:
            print(f"  - Type: {dep.get('type')}, Target: {dep.get('target_composite_id')}")
    
    # Step 7: List dependents for user-service (should show api-gateway)
    print("\nStep 7: Listing dependents for user-service...")
    user_dependents_response = client.call_tool("list_node_dependents", {
        "composite_id": node_ids["user-service"]
    })
    print(f"User Service dependents: {user_dependents_response}")
    
    if "result" in user_dependents_response:
        dependents = user_dependents_response["result"].get("dependents", [])
        print(f"User Service has {len(dependents)} dependents:")
        for dep in dependents:
            print(f"  - Type: {dep.get('type')}, Source: {dep.get('source_composite_id')}")
    
    # Step 8: List dependents for payment-service (should show user-service and api-gateway)
    print("\nStep 8: Listing dependents for payment-service...")
    payment_dependents_response = client.call_tool("list_node_dependents", {
        "composite_id": node_ids["payment-service"]
    })
    print(f"Payment Service dependents: {payment_dependents_response}")
    
    if "result" in payment_dependents_response:
        dependents = payment_dependents_response["result"].get("dependents", [])
        print(f"Payment Service has {len(dependents)} dependents:")
        for dep in dependents:
            print(f"  - Type: {dep.get('type')}, Source: {dep.get('source_composite_id')}")
    
    # Validation Summary
    print("\n=== VALIDATION SUMMARY ===")
    
    # Check hard dependency
    hard_dep_created = "result" in hard_dep_response and hard_dep_response["result"]
    print(f"✓ Hard dependency created: {hard_dep_created}")
    
    # Check soft dependency
    soft_dep_created = "result" in soft_dep_response and soft_dep_response["result"]
    print(f"✓ Soft dependency created: {soft_dep_created}")
    
    # Check reference dependency with metadata
    ref_dep_created = "result" in ref_dep_response and ref_dep_response["result"]
    print(f"✓ Reference dependency created: {ref_dep_created}")
    
    # Check api-gateway has 2 dependencies
    api_deps_count = 0
    if "result" in api_deps_response:
        api_deps_count = len(api_deps_response["result"].get("dependencies", []))
    print(f"✓ API Gateway dependencies count: {api_deps_count} (expected: 2)")
    
    # Check user-service has 1 dependent
    user_dependents_count = 0
    if "result" in user_dependents_response:
        user_dependents_count = len(user_dependents_response["result"].get("dependents", []))
    print(f"✓ User Service dependents count: {user_dependents_count} (expected: 1)")
    
    # Check payment-service has 2 dependents
    payment_dependents_count = 0
    if "result" in payment_dependents_response:
        payment_dependents_count = len(payment_dependents_response["result"].get("dependents", []))
    print(f"✓ Payment Service dependents count: {payment_dependents_count} (expected: 2)")
    
    # Overall validation
    all_valid = (
        hard_dep_created and soft_dep_created and ref_dep_created and
        api_deps_count == 2 and user_dependents_count == 1 and payment_dependents_count == 2
    )
    
    print(f"\n🎯 Overall scenario validation: {'✅ PASSED' if all_valid else '❌ FAILED'}")
    
    # Close the client
    client.close()

if __name__ == "__main__":
    main()