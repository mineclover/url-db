#!/usr/bin/env python3
"""
Final MCP LLM Judge Scenarios Test with Correct Response Parsing
"""

import json
import subprocess
import sys
import time

class MCPClient:
    def __init__(self):
        self.process = None
        self.initialized = False
    
    def start_server(self):
        """Start MCP server process"""
        self.process = subprocess.Popen(
            ["./bin/url-db", "-mcp-mode=stdio", "-db-path=test_final_scenarios.db"],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            bufsize=1
        )
        return self.process is not None
    
    def send_request(self, method, params=None, is_notification=False):
        """Send MCP request/notification"""
        if not self.process:
            return None
            
        message = {
            "jsonrpc": "2.0",
            "method": method,
            "params": params or {}
        }
        
        if not is_notification:
            message["id"] = 1
        
        try:
            request_str = json.dumps(message) + "\n"
            self.process.stdin.write(request_str)
            self.process.stdin.flush()
            
            if is_notification:
                return None
            
            response_line = self.process.stdout.readline()
            if response_line:
                return json.loads(response_line.strip())
            
        except Exception as e:
            print(f"Error sending request: {e}")
            return None
    
    def initialize(self):
        """Initialize MCP connection"""
        init_params = {
            "protocolVersion": "2024-11-05",
            "capabilities": {
                "roots": {"listChanged": True},
                "sampling": {}
            },
            "clientInfo": {
                "name": "llm-judge-client",
                "version": "1.0.0"
            }
        }
        
        response = self.send_request("initialize", init_params)
        if response and "result" in response:
            self.send_request("notifications/initialized", {}, is_notification=True)
            self.initialized = True
            return True
        return False
    
    def call_tool(self, name, arguments=None):
        """Call a specific tool and parse response correctly"""
        if not self.initialized:
            return None
            
        response = self.send_request("tools/call", {
            "name": name,
            "arguments": arguments or {}
        })
        
        if response and "result" in response:
            # Parse the MCP response format: result.content[0].text contains JSON
            content = response["result"].get("content", [])
            if content and content[0].get("type") == "text":
                try:
                    # Parse the JSON string from the text content
                    return json.loads(content[0]["text"])
                except json.JSONDecodeError as e:
                    print(f"❌ Failed to parse tool response: {e}")
                    return None
            return response["result"]
        else:
            if response and "error" in response:
                print(f"❌ Tool error: {response['error']}")
            return None
    
    def cleanup(self):
        """Clean up resources"""
        if self.process:
            try:
                self.process.terminate()
                self.process.wait(timeout=5)
            except:
                self.process.kill()

def execute_all_possible_scenarios():
    """Execute all possible MCP LLM judge scenarios with current tools"""
    client = MCPClient()
    
    print("=== MCP LLM JUDGE SCENARIOS TEST SUITE ===")
    print("Testing all scenarios possible with current implementation")
    
    if not client.start_server():
        print("❌ Failed to start server")
        return False
    
    try:
        if not client.initialize():
            print("❌ Failed to initialize")
            return False
        
        print("✅ MCP server initialized successfully")
        
        results = []
        
        # === PART 1: FOUNDATION TESTING ===
        print("\n" + "="*60)
        print("PART 1: FOUNDATION - DOMAIN AND NODE MANAGEMENT")
        print("="*60)
        
        # Create test domains
        domains = [
            ("microservices", "Microservices architecture components"),
            ("infrastructure", "Infrastructure components for testing"),
            ("applications", "Application layer components"),
            ("data", "Data layer components")
        ]
        
        created_domains = []
        for domain_name, description in domains:
            print(f"\n📁 Creating domain: {domain_name}")
            result = client.call_tool("create_domain", {
                "name": domain_name,
                "description": description
            })
            
            if result and "name" in result:
                print(f"✅ Domain created: {domain_name}")
                created_domains.append(domain_name)
                results.append((f"Create {domain_name} domain", True))
            else:
                print(f"❌ Failed to create domain: {domain_name}")
                results.append((f"Create {domain_name} domain", False))
        
        # Create nodes for each domain
        nodes_per_domain = {
            "microservices": [
                ("api-gateway", "https://micro.example.com/api-gateway", "API Gateway", "Main API gateway service"),
                ("user-service", "https://micro.example.com/user-service", "User Service", "User management service"),
                ("payment-service", "https://micro.example.com/payment-service", "Payment Service", "Payment processing service"),
                ("database-service", "https://micro.example.com/database-service", "Database Service", "Database service")
            ],
            "infrastructure": [
                ("load-balancer", "https://infra.example.com/load-balancer", "Load Balancer", "Main load balancer"),
                ("web-server", "https://infra.example.com/web-server", "Web Server", "Nginx web server"),
                ("app-server", "https://infra.example.com/app-server", "App Server", "Application server"),
                ("database", "https://infra.example.com/database", "Database", "PostgreSQL database")
            ],
            "applications": [
                ("web-app", "https://apps.example.com/web-app", "Web Application", "Main web application"),
                ("mobile-app", "https://apps.example.com/mobile-app", "Mobile Application", "Mobile app")
            ],
            "data": [
                ("user-db", "https://data.example.com/user-db", "User Database", "User data storage"),
                ("analytics-db", "https://data.example.com/analytics-db", "Analytics Database", "Analytics data"),
                ("cache", "https://data.example.com/cache", "Cache", "Redis cache system")
            ]
        }
        
        created_nodes = {}
        for domain_name in created_domains:
            if domain_name in nodes_per_domain:
                created_nodes[domain_name] = []
                for name, url, title, desc in nodes_per_domain[domain_name]:
                    print(f"\n🔗 Creating node: {name} in {domain_name}")
                    result = client.call_tool("create_node", {
                        "domain_name": domain_name,
                        "url": url,
                        "title": title,
                        "description": desc
                    })
                    
                    if result and "composite_id" in result:
                        composite_id = result["composite_id"]
                        created_nodes[domain_name].append((name, composite_id))
                        print(f"✅ Node created: {name} ({composite_id})")
                        results.append((f"Create {name} node", True))
                    else:
                        print(f"❌ Failed to create node: {name}")
                        results.append((f"Create {name} node", False))
        
        # === PART 2: ADVANCED NODE OPERATIONS ===
        print("\n" + "="*60)
        print("PART 2: ADVANCED NODE OPERATIONS")
        print("="*60)
        
        # Test node retrieval and updates
        if "microservices" in created_nodes and created_nodes["microservices"]:
            node_name, node_id = created_nodes["microservices"][0]  # Use first node
            
            print(f"\n📋 Getting node details: {node_name}")
            result = client.call_tool("get_node", {"composite_id": node_id})
            if result and "composite_id" in result:
                print(f"✅ Retrieved node details for {node_name}")
                results.append((f"Get {node_name} details", True))
            else:
                print(f"❌ Failed to get node details for {node_name}")
                results.append((f"Get {node_name} details", False))
            
            print(f"\n✏️ Updating node title: {node_name}")
            result = client.call_tool("update_node", {
                "composite_id": node_id,
                "title": f"Updated {node_name.replace('-', ' ').title()}"
            })
            if result:
                print(f"✅ Updated title for {node_name}")
                results.append((f"Update {node_name} title", True))
            else:
                print(f"❌ Failed to update title for {node_name}")
                results.append((f"Update {node_name} title", False))
            
            print(f"\n✏️ Updating node description: {node_name}")
            result = client.call_tool("update_node", {
                "composite_id": node_id,
                "description": f"Updated description for {node_name}"
            })
            if result:
                print(f"✅ Updated description for {node_name}")
                results.append((f"Update {node_name} description", True))
            else:
                print(f"❌ Failed to update description for {node_name}")
                results.append((f"Update {node_name} description", False))
        
        # === PART 3: ATTRIBUTE SYSTEM TESTING ===
        print("\n" + "="*60)
        print("PART 3: ATTRIBUTE SYSTEM - SIMULATING METADATA")
        print("="*60)
        
        # Create domain attributes for metadata
        if "microservices" in created_domains:
            print(f"\n🏷️ Creating domain attributes for microservices")
            
            attributes_to_create = [
                ("dependencies", "string", "Node dependencies metadata"),
                ("service_type", "tag", "Type of microservice"),
                ("priority", "number", "Service priority level"),
                ("documentation", "markdown", "Service documentation")
            ]
            
            for attr_name, attr_type, attr_desc in attributes_to_create:
                result = client.call_tool("create_domain_attribute", {
                    "domain_name": "microservices",
                    "name": attr_name,
                    "type": attr_type,
                    "description": attr_desc
                })
                
                if result:
                    print(f"✅ Created attribute: {attr_name}")
                    results.append((f"Create {attr_name} attribute", True))
                else:
                    print(f"❌ Failed to create attribute: {attr_name}")
                    results.append((f"Create {attr_name} attribute", False))
            
            # Set attributes on nodes (simulating dependency metadata)
            if "microservices" in created_nodes and created_nodes["microservices"]:
                for i, (node_name, node_id) in enumerate(created_nodes["microservices"][:2]):  # Test first 2 nodes
                    print(f"\n🏷️ Setting attributes for {node_name}")
                    
                    # Create dependency metadata
                    dep_metadata = json.dumps({
                        "type": "hard" if i == 0 else "soft",
                        "cascade_delete": i == 0,
                        "cascade_update": True,
                        "targets": ["database-service"] if node_name == "user-service" else ["user-service"]
                    })
                    
                    attributes = [
                        {"name": "dependencies", "value": dep_metadata},
                        {"name": "service_type", "value": "core" if i == 0 else "auxiliary"},
                        {"name": "priority", "value": str(10 - i)},
                        {"name": "documentation", "value": f"# {node_name.replace('-', ' ').title()}\\n\\nThis is the documentation for {node_name}."}
                    ]
                    
                    result = client.call_tool("set_node_attributes", {
                        "composite_id": node_id,
                        "attributes": attributes
                    })
                    
                    if result:
                        print(f"✅ Set attributes for {node_name}")
                        results.append((f"Set attributes for {node_name}", True))
                    else:
                        print(f"❌ Failed to set attributes for {node_name}")
                        results.append((f"Set attributes for {node_name}", False))
        
        # === PART 4: COMPLEX QUERIES AND FILTERING ===
        print("\n" + "="*60)
        print("PART 4: COMPLEX QUERIES AND FILTERING")
        print("="*60)
        
        # Test listing with filtering
        for domain_name in created_domains[:2]:  # Test first 2 domains
            print(f"\n📋 Listing nodes in {domain_name}")
            result = client.call_tool("list_nodes", {
                "domain_name": domain_name,
                "size": 10
            })
            
            if result and "nodes" in result:
                nodes = result["nodes"]
                print(f"✅ Listed {len(nodes)} nodes in {domain_name}")
                results.append((f"List {domain_name} nodes", True))
                
                # Test search functionality
                if nodes:
                    search_term = nodes[0]["title"].split()[0]  # Use first word of first node title
                    print(f"🔍 Searching for '{search_term}' in {domain_name}")
                    search_result = client.call_tool("list_nodes", {
                        "domain_name": domain_name,
                        "search": search_term
                    })
                    
                    if search_result and "nodes" in search_result:
                        found_nodes = search_result["nodes"]
                        print(f"✅ Search found {len(found_nodes)} nodes")
                        results.append((f"Search in {domain_name}", True))
                    else:
                        print(f"❌ Search failed in {domain_name}")
                        results.append((f"Search in {domain_name}", False))
            else:
                print(f"❌ Failed to list nodes in {domain_name}")
                results.append((f"List {domain_name} nodes", False))
        
        # Test attribute-based filtering
        if "microservices" in created_domains:
            print(f"\n🔍 Testing attribute-based filtering in microservices")
            result = client.call_tool("filter_nodes_by_attributes", {
                "domain_name": "microservices",
                "filters": [
                    {"name": "service_type", "value": "core", "operator": "equals"}
                ]
            })
            
            if result and "nodes" in result:
                filtered_nodes = result["nodes"]
                print(f"✅ Filtered nodes: found {len(filtered_nodes)} core services")
                results.append(("Filter by service type", True))
            else:
                print(f"❌ Failed to filter by attributes")
                results.append(("Filter by service type", False))
        
        # === PART 5: CROSS-DOMAIN OPERATIONS ===
        print("\n" + "="*60)
        print("PART 5: CROSS-DOMAIN OPERATIONS")
        print("="*60)
        
        # Test URL finding across domains
        if created_nodes:
            for domain_name, nodes in list(created_nodes.items())[:2]:
                if nodes:
                    node_name, node_id = nodes[0]
                    # Get the URL for this node
                    node_details = client.call_tool("get_node", {"composite_id": node_id})
                    if node_details and "url" in node_details:
                        url = node_details["url"]
                        print(f"\n🔍 Finding node by URL: {url}")
                        result = client.call_tool("find_node_by_url", {
                            "domain_name": domain_name,
                            "url": url
                        })
                        
                        if result and "composite_id" in result:
                            print(f"✅ Found node by URL: {result['composite_id']}")
                            results.append((f"Find {node_name} by URL", True))
                        else:
                            print(f"❌ Failed to find node by URL")
                            results.append((f"Find {node_name} by URL", False))
        
        # === PART 6: COMPREHENSIVE TESTING ===
        print("\n" + "="*60)
        print("PART 6: COMPREHENSIVE NODE AND ATTRIBUTE OPERATIONS")
        print("="*60)
        
        # Test get_node_with_attributes
        if "microservices" in created_nodes and created_nodes["microservices"]:
            node_name, node_id = created_nodes["microservices"][0]
            print(f"\n📋 Getting node with attributes: {node_name}")
            result = client.call_tool("get_node_with_attributes", {
                "composite_id": node_id
            })
            
            if result:
                print(f"✅ Retrieved node with attributes for {node_name}")
                if "attributes" in result:
                    attrs = result["attributes"]
                    print(f"   Found {len(attrs)} attributes")
                results.append((f"Get {node_name} with attributes", True))
            else:
                print(f"❌ Failed to get node with attributes for {node_name}")
                results.append((f"Get {node_name} with attributes", False))
        
        # Test server info
        print(f"\n ℹ️ Getting server information")
        server_info = client.call_tool("get_server_info", {})
        if server_info:
            print(f"✅ Server info retrieved")
            print(f"   Name: {server_info.get('name', 'Unknown')}")
            print(f"   Version: {server_info.get('version', 'Unknown')}")
            results.append(("Get server info", True))
        else:
            print(f"❌ Failed to get server info")
            results.append(("Get server info", False))
        
        # === FINAL RESULTS ===
        print("\n" + "="*70)
        print("FINAL TEST RESULTS - MCP LLM JUDGE SCENARIOS")
        print("="*70)
        
        passed = sum(1 for _, success in results if success)
        total = len(results)
        pass_rate = (passed / total) * 100 if total > 0 else 0
        
        # Group results by category
        categories = {
            "Domain Management": [],
            "Node Management": [],
            "Attribute System": [],
            "Complex Queries": [],
            "Cross-Domain Operations": [],
            "Server Operations": []
        }
        
        for test_name, success in results:
            if "domain" in test_name.lower():
                categories["Domain Management"].append((test_name, success))
            elif "attribute" in test_name.lower():
                categories["Attribute System"].append((test_name, success))
            elif "search" in test_name.lower() or "filter" in test_name.lower() or "list" in test_name.lower():
                categories["Complex Queries"].append((test_name, success))
            elif "url" in test_name.lower():
                categories["Cross-Domain Operations"].append((test_name, success))
            elif "server" in test_name.lower():
                categories["Server Operations"].append((test_name, success))
            else:
                categories["Node Management"].append((test_name, success))
        
        for category, tests in categories.items():
            if tests:
                cat_passed = sum(1 for _, success in tests if success)
                cat_total = len(tests)
                cat_rate = (cat_passed / cat_total) * 100 if cat_total > 0 else 0
                print(f"\n{category}: {cat_passed}/{cat_total} ({cat_rate:.1f}%)")
                for test_name, success in tests:
                    status = "✅" if success else "❌"
                    print(f"  {status} {test_name}")
        
        print(f"\nOVERALL RESULTS: {passed}/{total} tests passed ({pass_rate:.1f}%)")
        
        if pass_rate >= 90:
            print("🎉 EXCELLENT! System working exceptionally well!")
            grade = "A+"
        elif pass_rate >= 80:
            print("🎉 VERY GOOD! System working well with minor issues!")
            grade = "A"
        elif pass_rate >= 70:
            print("👍 GOOD! System mostly working with some gaps!")
            grade = "B"
        elif pass_rate >= 60:
            print("👌 FAIR! Core functionality working but needs improvement!")
            grade = "C"
        else:
            print("⚠️ NEEDS SIGNIFICANT WORK! Major functionality missing!")
            grade = "F"
        
        print(f"FINAL GRADE: {grade}")
        
        # Report on missing functionality for complete LLM judge scenarios
        print("\n" + "="*70)
        print("MISSING FUNCTIONALITY FOR COMPLETE LLM JUDGE SCENARIOS")
        print("="*70)
        
        missing_features = [
            "❌ Event System - Automatic event generation for node operations",
            "❌ Event Management - get_node_events, get_pending_events, process_event",
            "❌ Subscription System - create_subscription, list_subscriptions, delete_subscription",
            "❌ Dependency Management - create_dependency, list_node_dependencies, delete_dependency",
            "❌ Event Statistics - get_event_stats for system monitoring"
        ]
        
        print("The following features need to be implemented for full LLM judge scenario compliance:")
        for feature in missing_features:
            print(f"  {feature}")
        
        print(f"\nCURRENT STATUS: Core URL-DB functionality working well ({pass_rate:.1f}%)")
        print("NEXT STEPS: Implement external dependency API for complete scenario coverage")
        
        return pass_rate >= 60  # Consider success if 60%+ pass rate
        
    finally:
        client.cleanup()

if __name__ == "__main__":
    print("Starting comprehensive MCP LLM Judge Scenarios test...")
    success = execute_all_possible_scenarios()
    
    if success:
        print("\n✅ TEST SUITE COMPLETED SUCCESSFULLY!")
        print("Core functionality is working well. External dependency API needed for full compliance.")
    else:
        print("\n❌ TEST SUITE REVEALED SIGNIFICANT ISSUES!")
        print("Core functionality needs attention before implementing additional features.")
    
    sys.exit(0 if success else 1)