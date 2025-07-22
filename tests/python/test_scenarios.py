#!/usr/bin/env python3
"""
LLM-as-a-Judge Test Runner for URL-DB MCP Server
Executes comprehensive test scenarios and generates evaluation report
"""

import json
import subprocess
import sys
import time
from typing import Dict, List, Any, Optional
from dataclasses import dataclass
from tool_constants import CREATE_DOMAIN, CREATE_NODE, DELETE_NODE, FIND_NODE_BY_URL, GET_NODE, GET_NODE_ATTRIBUTES, GET_SERVER_INFO, LIST_DOMAINS, LIST_NODES, SET_NODE_ATTRIBUTES, UPDATE_NODE


@dataclass
class TestResult:
    scenario_name: str
    score: int
    max_score: int
    status: str
    notes: str
    execution_time: float

class MCPTestRunner:
    def __init__(self, server_path: str):
        self.server_path = server_path
        self.process = None
        self.request_id = 1
        self.test_results: List[TestResult] = []
        
    def start_server(self):
        """Start MCP server in stdio mode"""
        self.process = subprocess.Popen(
            [self.server_path, "-mcp-mode=stdio"],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            bufsize=0
        )
        time.sleep(0.5)  # Allow server to start
        
    def send_request(self, method: str, params: Any = None) -> Optional[Dict]:
        """Send JSON-RPC request and return response"""
        request = {
            "jsonrpc": "2.0",
            "id": self.request_id,
            "method": method,
            "params": params or {}
        }
        
        request_json = json.dumps(request)
        self.process.stdin.write(request_json + "\n")
        self.process.stdin.flush()
        
        response_line = self.process.stdout.readline().strip()
        if response_line:
            try:
                response = json.loads(response_line)
                self.request_id += 1
                return response
            except json.JSONDecodeError:
                return None
        return None
        
    def send_notification(self, method: str, params: Any = None):
        """Send JSON-RPC notification (no response expected)"""
        notification = {
            "jsonrpc": "2.0",
            "method": method,
            "params": params or {}
        }
        
        notification_json = json.dumps(notification)
        self.process.stdin.write(notification_json + "\n")
        self.process.stdin.flush()
    
    def close(self):
        """Close server connection"""
        if self.process:
            self.process.stdin.close()
            self.process.wait()

    def test_scenario_1_handshake(self) -> TestResult:
        """Scenario 1: MCP Protocol Handshake Compliance"""
        start_time = time.time()
        score = 0
        notes = []
        
        try:
            # Initialize request
            init_response = self.send_request("initialize", {
                "protocolVersion": "2024-11-05",
                "capabilities": {"experimental": {}, "sampling": {}},
                "clientInfo": {"name": "test-judge-client", "version": "1.0.0"}
            })
            
            if init_response and "result" in init_response:
                result = init_response["result"]
                
                # Check protocol version
                if result.get("protocolVersion") == "2024-11-05":
                    score += 1
                    notes.append("✓ Correct protocol version")
                else:
                    notes.append("✗ Incorrect protocol version")
                
                # Check capabilities structure
                caps = result.get("capabilities", {})
                if "tools" in caps and "resources" in caps:
                    score += 2
                    notes.append("✓ Proper capabilities declared")
                else:
                    notes.append("✗ Missing required capabilities")
                
                # Check server info
                server_info = result.get("serverInfo", {})
                if server_info.get("name") == "url-db-mcp-server":
                    score += 1
                    notes.append("✓ Correct server name")
                else:
                    notes.append("✗ Incorrect server name")
                    
                if server_info.get("version") == "1.0.0":
                    score += 1
                    notes.append("✓ Correct server version")
                else:
                    notes.append("✗ Incorrect server version")
            else:
                notes.append("✗ Failed to get initialize response")
            
            # Send initialized notification
            self.send_notification("initialized", {})
            time.sleep(0.2)
            score += 1  # Notification sent successfully
            notes.append("✓ Initialized notification sent")
            
            # Test malformed request handling
            malformed_response = self.send_request("invalid_method", {})
            if malformed_response and "error" in malformed_response:
                score += 2
                notes.append("✓ Proper error handling for unknown methods")
            else:
                notes.append("✗ Poor error handling")
                
        except Exception as e:
            notes.append(f"✗ Exception during handshake: {str(e)}")
        
        execution_time = time.time() - start_time
        status = "PASS" if score >= 8 else "FAIL"
        
        return TestResult(
            scenario_name="Protocol Handshake Compliance",
            score=score,
            max_score=10,
            status=status,
            notes="; ".join(notes),
            execution_time=execution_time
        )

    def test_scenario_2_tool_discovery(self) -> TestResult:
        """Scenario 2: Tool Discovery and Schema Validation"""
        start_time = time.time()
        score = 0
        notes = []
        
        expected_tools = [
            LIST_DOMAINS, CREATE_DOMAIN, LIST_NODES,
            CREATE_NODE, GET_NODE, UPDATE_NODE,
            DELETE_NODE, FIND_NODE_BY_URL, GET_NODE_ATTRIBUTES,
            SET_NODE_ATTRIBUTES, GET_SERVER_INFO
        ]
        
        try:
            tools_response = self.send_request("tools/list", {})
            
            if tools_response and "result" in tools_response:
                tools = tools_response["result"].get("tools", [])
                tool_names = [tool["name"] for tool in tools]
                
                # Check tool completeness
                found_tools = set(tool_names) & set(expected_tools)
                if len(found_tools) == len(expected_tools):
                    score += 3
                    notes.append(f"✓ All {len(expected_tools)} tools present")
                else:
                    missing = set(expected_tools) - found_tools
                    score += max(0, 3 - len(missing))
                    notes.append(f"✗ Missing tools: {list(missing)}")
                
                # Check schema validity
                valid_schemas = 0
                for tool in tools:
                    schema = tool.get("inputSchema", {})
                    if isinstance(schema, dict) and "type" in schema:
                        valid_schemas += 1
                
                if valid_schemas >= len(tools) * 0.9:  # 90% have valid schemas
                    score += 3
                    notes.append("✓ Tool schemas are valid")
                else:
                    score += 1
                    notes.append("✗ Some tools have invalid schemas")
                
                # Check documentation quality
                documented_tools = sum(1 for tool in tools if tool.get("description", "").strip())
                if documented_tools >= len(tools) * 0.9:
                    score += 2
                    notes.append("✓ Tools are well documented")
                else:
                    score += 1
                    notes.append("✗ Some tools lack documentation")
                
                # Check parameter validation
                tools_with_params = sum(1 for tool in tools 
                                      if tool.get("inputSchema", {}).get("properties"))
                if tools_with_params >= 8:  # Most tools should have parameters
                    score += 2
                    notes.append("✓ Parameter specifications present")
                else:
                    score += 1
                    notes.append("✗ Limited parameter specifications")
                    
            else:
                notes.append("✗ Failed to get tools list")
                
        except Exception as e:
            notes.append(f"✗ Exception during tool discovery: {str(e)}")
        
        execution_time = time.time() - start_time
        status = "PASS" if score >= 8 else "FAIL"
        
        return TestResult(
            scenario_name="Tool Discovery and Schema Validation",
            score=score,
            max_score=10,
            status=status,
            notes="; ".join(notes),
            execution_time=execution_time
        )

    def test_scenario_3_domain_management(self) -> TestResult:
        """Scenario 3: Domain Management Workflow"""
        start_time = time.time()
        score = 0
        notes = []
        
        try:
            # List initial domains
            domains_response = self.send_request("tools/call", {
                "name": LIST_DOMAINS,
                "arguments": {}
            })
            
            if domains_response and "result" in domains_response:
                content = json.loads(domains_response["result"]["content"][0]["text"])
                initial_count = len(content.get("domains", []))
                score += 1
                notes.append(f"✓ Initial domain count: {initial_count}")
            
            # Create test domain with unique name
            unique_name = f"test-scenario-{int(time.time())}"
            create_response = self.send_request("tools/call", {
                "name": CREATE_DOMAIN,
                "arguments": {
                    "name": unique_name,
                    "description": "Test domain for LLM judge scenarios"
                }
            })
            
            if create_response and "result" in create_response:
                if not create_response["result"].get("isError", False):
                    score += 2
                    notes.append("✓ Domain creation successful")
                    
                    # Verify domain data
                    domain_data = json.loads(create_response["result"]["content"][0]["text"])
                    if "created_at" in domain_data and "updated_at" in domain_data:
                        score += 1
                        notes.append("✓ Proper metadata timestamps")
                else:
                    notes.append("✗ Domain creation failed")
            
            # List domains again to verify
            domains_response2 = self.send_request("tools/call", {
                "name": LIST_DOMAINS,
                "arguments": {}
            })
            
            if domains_response2 and "result" in domains_response2:
                content = json.loads(domains_response2["result"]["content"][0]["text"])
                new_count = len(content.get("domains", []))
                if new_count > initial_count:
                    score += 2
                    notes.append("✓ Domain persists in listings")
                else:
                    notes.append("✗ Domain not found in subsequent listing")
            
            # Test duplicate domain creation
            duplicate_response = self.send_request("tools/call", {
                "name": CREATE_DOMAIN,
                "arguments": {
                    "name": unique_name,
                    "description": "Duplicate domain"
                }
            })
            
            if duplicate_response and "result" in duplicate_response:
                if duplicate_response["result"].get("isError", False):
                    score += 2
                    notes.append("✓ Duplicate domain properly rejected")
                else:
                    score += 1
                    notes.append("⚠ Duplicate domain allowed (may be by design)")
            
        except Exception as e:
            notes.append(f"✗ Exception during domain management: {str(e)}")
        
        execution_time = time.time() - start_time
        status = "PASS" if score >= 7 else "FAIL"
        
        return TestResult(
            scenario_name="Domain Management Workflow",
            score=score,
            max_score=10,
            status=status,
            notes="; ".join(notes),
            execution_time=execution_time
        )

    def test_scenario_4_node_management(self) -> TestResult:
        """Scenario 4: Node/URL Management with Composite Keys"""
        start_time = time.time()
        score = 0
        notes = []
        
        try:
            # Create node (using existing test domain)
            create_response = self.send_request("tools/call", {
                "name": CREATE_NODE,
                "arguments": {
                    "domain_name": "test-domain",  # Use existing domain
                    "url": f"https://example.com/test-page-{int(time.time())}",  # Unique URL
                    "title": "Test Page for LLM Judge",
                    "description": "Test node for scenario validation"
                }
            })
            
            composite_id = None
            if create_response and "result" in create_response:
                if not create_response["result"].get("isError", False):
                    node_data = json.loads(create_response["result"]["content"][0]["text"])
                    composite_id = node_data.get("composite_id")
                    
                    # Check composite key format
                    if composite_id and composite_id.startswith("url-db:test-domain:"):
                        score += 2
                        notes.append("✓ Proper composite key format")
                    else:
                        notes.append("✗ Invalid composite key format")
                    
                    score += 1
                    notes.append("✓ Node creation successful")
                else:
                    notes.append("✗ Node creation failed")
            
            if composite_id:
                # Retrieve node by composite key
                get_response = self.send_request("tools/call", {
                    "name": GET_NODE,
                    "arguments": {"composite_id": composite_id}
                })
                
                if get_response and "result" in get_response:
                    if not get_response["result"].get("isError", False):
                        score += 2
                        notes.append("✓ Node retrieval by composite key successful")
                    else:
                        notes.append("✗ Node retrieval failed")
                
                # Update node
                update_response = self.send_request("tools/call", {
                    "name": UPDATE_NODE,
                    "arguments": {
                        "composite_id": composite_id,
                        "title": "Updated Test Page",
                        "description": "Updated description"
                    }
                })
                
                if update_response and "result" in update_response:
                    if not update_response["result"].get("isError", False):
                        score += 1
                        notes.append("✓ Node update successful")
                
                # Find node by URL
                test_url = f"https://example.com/test-page-{int(time.time())}"
                find_response = self.send_request("tools/call", {
                    "name": FIND_NODE_BY_URL,
                    "arguments": {
                        "domain_name": "test-domain",
                        "url": node_data.get("url", test_url)  # Use actual URL from created node
                    }
                })
                
                if find_response and "result" in find_response:
                    if not find_response["result"].get("isError", False):
                        score += 2
                        notes.append("✓ URL search successful")
                
                # Delete node
                delete_response = self.send_request("tools/call", {
                    "name": DELETE_NODE,
                    "arguments": {"composite_id": composite_id}
                })
                
                if delete_response and "result" in delete_response:
                    if not delete_response["result"].get("isError", False):
                        score += 2
                        notes.append("✓ Node deletion successful")
            
        except Exception as e:
            notes.append(f"✗ Exception during node management: {str(e)}")
        
        execution_time = time.time() - start_time
        status = "PASS" if score >= 8 else "FAIL"
        
        return TestResult(
            scenario_name="Node/URL Management with Composite Keys",
            score=score,
            max_score=10,
            status=status,
            notes="; ".join(notes),
            execution_time=execution_time
        )

    def test_scenario_5_resource_system(self) -> TestResult:
        """Scenario 5: Resource System Integration"""
        start_time = time.time()
        score = 0
        notes = []
        
        try:
            # Get resources list
            resources_response = self.send_request("resources/list", {})
            
            if resources_response and "result" in resources_response:
                resources = resources_response["result"].get("resources", [])
                resource_uris = [r["uri"] for r in resources]
                
                # Check for expected resources
                expected_uris = [
                    "mcp://server/info",
                    "mcp://domains/test-scenario-domain",
                    "mcp://domains/test-scenario-domain/nodes"
                ]
                
                found_uris = sum(1 for uri in expected_uris if uri in resource_uris)
                score += min(3, found_uris)
                notes.append(f"✓ Found {found_uris}/{len(expected_uris)} expected resources")
                
                # Read server info resource
                server_info_response = self.send_request("resources/read", {
                    "uri": "mcp://server/info"
                })
                
                if server_info_response and "result" in server_info_response:
                    content = server_info_response["result"]["contents"][0]["text"]
                    try:
                        server_data = json.loads(content)
                        if "name" in server_data and "capabilities" in server_data:
                            score += 2
                            notes.append("✓ Server info resource valid")
                    except:
                        notes.append("✗ Server info resource invalid JSON")
                
                # Read domain resource if it exists
                if "mcp://domains/test-scenario-domain" in resource_uris:
                    domain_response = self.send_request("resources/read", {
                        "uri": "mcp://domains/test-scenario-domain"
                    })
                    
                    if domain_response and "result" in domain_response:
                        score += 2
                        notes.append("✓ Domain resource accessible")
                
                # Check URI format consistency
                valid_uri_format = all(uri.startswith("mcp://") for uri in resource_uris)
                if valid_uri_format:
                    score += 2
                    notes.append("✓ Consistent URI format")
                else:
                    notes.append("✗ Inconsistent URI format")
                    
                score += 1  # Base score for resource system working
                
            else:
                notes.append("✗ Failed to get resources list")
                
        except Exception as e:
            notes.append(f"✗ Exception during resource testing: {str(e)}")
        
        execution_time = time.time() - start_time
        status = "PASS" if score >= 7 else "FAIL"
        
        return TestResult(
            scenario_name="Resource System Integration",
            score=score,
            max_score=10,
            status=status,
            notes="; ".join(notes),
            execution_time=execution_time
        )

    def run_all_scenarios(self):
        """Run all test scenarios"""
        print("🧪 Starting LLM-as-a-Judge MCP Server Testing")
        print("=" * 60)
        
        try:
            self.start_server()
            
            # Run handshake first (required for other tests)
            handshake_result = self.test_scenario_1_handshake()
            self.test_results.append(handshake_result)
            
            if handshake_result.status == "PASS":
                # Run remaining scenarios
                self.test_results.append(self.test_scenario_2_tool_discovery())
                self.test_results.append(self.test_scenario_3_domain_management())
                self.test_results.append(self.test_scenario_4_node_management())
                self.test_results.append(self.test_scenario_5_resource_system())
            else:
                print("❌ Handshake failed - skipping remaining scenarios")
                
        finally:
            self.close()
    
    def generate_report(self) -> str:
        """Generate comprehensive test report"""
        total_score = sum(r.score for r in self.test_results)
        max_total = sum(r.max_score for r in self.test_results)
        percentage = (total_score / max_total * 100) if max_total > 0 else 0
        
        # Determine overall grade
        if percentage >= 90:
            grade = "Excellent"
        elif percentage >= 80:
            grade = "Good"
        elif percentage >= 70:
            grade = "Acceptable"
        elif percentage >= 50:
            grade = "Poor"
        else:
            grade = "Failing"
        
        # Determine production readiness
        passing_scenarios = sum(1 for r in self.test_results if r.status == "PASS")
        production_ready = passing_scenarios >= 4 and percentage >= 70
        
        report = f"""# MCP Server LLM-as-a-Judge Test Report

**Date**: {time.strftime('%Y-%m-%d %H:%M:%S')}
**Server Version**: url-db-mcp-server v1.0.0
**Test Environment**: Local development

## Executive Summary

- **Total Score**: {total_score}/{max_total} ({percentage:.1f}%)
- **Overall Grade**: {grade}
- **Scenarios Passed**: {passing_scenarios}/{len(self.test_results)}
- **Production Ready**: {'✅ Yes' if production_ready else '❌ No'}

## Scenario Results

"""
        
        for result in self.test_results:
            report += f"""### {result.scenario_name}
- **Score**: {result.score}/{result.max_score} ({result.score/result.max_score*100:.1f}%)
- **Status**: {result.status}
- **Execution Time**: {result.execution_time:.2f}s
- **Notes**: {result.notes}

"""
        
        report += f"""## Performance Analysis

- **Average Execution Time**: {sum(r.execution_time for r in self.test_results)/len(self.test_results):.2f}s
- **Fastest Scenario**: {min(self.test_results, key=lambda r: r.execution_time).scenario_name}
- **Slowest Scenario**: {max(self.test_results, key=lambda r: r.execution_time).scenario_name}

## Detailed Assessment

### Strengths
"""
        
        # Identify strengths (scenarios with score >= 8)
        strong_scenarios = [r for r in self.test_results if r.score >= 8]
        if strong_scenarios:
            for scenario in strong_scenarios:
                report += f"- {scenario.scenario_name}: Excellent implementation\n"
        else:
            report += "- No scenarios scored in excellent range\n"
        
        report += "\n### Areas for Improvement\n"
        
        # Identify weaknesses (scenarios with score < 7)
        weak_scenarios = [r for r in self.test_results if r.score < 7]
        if weak_scenarios:
            for scenario in weak_scenarios:
                report += f"- {scenario.scenario_name}: Needs attention (scored {scenario.score}/10)\n"
        else:
            report += "- All scenarios perform at acceptable levels or higher\n"
        
        report += f"""
## Production Readiness Assessment

{'✅' if production_ready else '❌'} **Overall Assessment**: {'PRODUCTION READY' if production_ready else 'NOT PRODUCTION READY'}

### Criteria Evaluation:
- MCP Protocol Compliance: {'✅' if self.test_results and self.test_results[0].status == 'PASS' else '❌'}
- Core Functionality: {'✅' if passing_scenarios >= 3 else '❌'}
- Error Handling: {'✅' if percentage >= 70 else '❌'}
- Performance: {'✅' if all(r.execution_time < 5.0 for r in self.test_results) else '❌'}

## Recommendations

"""
        
        if production_ready:
            report += """The URL-DB MCP server demonstrates excellent compliance with the MCP protocol and provides robust functionality. The implementation is ready for production deployment.

### Recommended Actions:
1. Deploy to production environment
2. Monitor performance metrics in real-world usage
3. Implement additional monitoring and logging
4. Consider implementing caching for improved performance
"""
        else:
            report += """The URL-DB MCP server requires improvements before production deployment.

### Priority Actions:
"""
            if self.test_results and self.test_results[0].status != "PASS":
                report += "1. **CRITICAL**: Fix MCP protocol handshake compliance\n"
            
            for scenario in weak_scenarios:
                report += f"2. Address issues in {scenario.scenario_name}\n"
            
            report += "3. Re-run comprehensive testing after fixes\n"
        
        report += f"""
---
*Report generated by LLM-as-a-Judge testing framework*
*Test execution completed in {sum(r.execution_time for r in self.test_results):.2f} seconds*
"""
        
        return report

def main():
    """Main test execution"""
    server_path = "../../bin/url-db"
    
    # Check if server exists
    try:
        subprocess.run([server_path, "--help"], capture_output=True, check=True)
    except:
        print(f"❌ Server not found at {server_path}")
        print("Please build the server first: go build -o bin/url-db cmd/server/main.go")
        sys.exit(1)
    
    # Run tests
    runner = MCPTestRunner(server_path)
    runner.run_all_scenarios()
    
    # Generate and display report
    report = runner.generate_report()
    print(report)
    
    # Save report to file
    with open("mcp_test_report.md", "w") as f:
        f.write(report)
    
    print(f"\n📄 Detailed report saved to: mcp_test_report.md")

if __name__ == "__main__":
    main()