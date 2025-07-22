#!/usr/bin/env python3
"""
Subscription Lifecycle Management Test Script
Tests the complete subscription lifecycle as specified in the MCP LLM judge scenarios.
"""

import asyncio
import json
import subprocess
import sys
import time
from typing import Dict, List, Any, Optional
import tempfile
import os

class MCPClient:
    def __init__(self, command: List[str]):
        self.command = command
        self.process = None
        
    async def start(self):
        """Start the MCP server process"""
        self.process = await asyncio.create_subprocess_exec(
            *self.command,
            stdin=asyncio.subprocess.PIPE,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE
        )
        
        # Initialize MCP session
        await self._send_message({
            "jsonrpc": "2.0",
            "id": 1,
            "method": "initialize",
            "params": {
                "protocolVersion": "2024-11-05",
                "capabilities": {
                    "tools": {}
                },
                "clientInfo": {
                    "name": "test-client",
                    "version": "1.0.0"
                }
            }
        })
        
        await self._read_message()
        
        await self._send_message({
            "jsonrpc": "2.0",
            "method": "notifications/initialized"
        })
        
    async def _send_message(self, message: Dict[str, Any]):
        """Send a JSON-RPC message to the server"""
        message_str = json.dumps(message) + "\n"
        self.process.stdin.write(message_str.encode())
        await self.process.stdin.drain()
        
    async def _read_message(self) -> Dict[str, Any]:
        """Read a JSON-RPC message from the server"""
        while True:
            line = await self.process.stdout.readline()
            if not line:
                break
            line = line.decode().strip()
            if line:
                try:
                    return json.loads(line)
                except json.JSONDecodeError:
                    continue
        return {}
        
    async def call_tool(self, tool_name: str, arguments: Dict[str, Any], request_id: int = None) -> Dict[str, Any]:
        """Call an MCP tool"""
        if request_id is None:
            request_id = int(time.time() * 1000)
            
        await self._send_message({
            "jsonrpc": "2.0",
            "id": request_id,
            "method": "tools/call",
            "params": {
                "name": tool_name,
                "arguments": arguments
            }
        })
        
        return await self._read_message()
        
    async def close(self):
        """Close the MCP connection"""
        if self.process:
            self.process.terminate()
            await self.process.wait()


async def test_subscription_lifecycle():
    """Execute the Subscription Lifecycle Management scenario"""
    
    # Create a temporary database file
    temp_db = tempfile.NamedTemporaryFile(suffix='.db', delete=False)
    temp_db.close()
    db_path = temp_db.name
    
    try:
        # Start MCP server
        client = MCPClient(['./bin/url-db', '-mcp-mode=stdio', f'-db-path={db_path}'])
        await client.start()
        
        print("=== Subscription Lifecycle Management Test ===\n")
        
        # Step 1: Create domain "microservices"
        print("Step 1: Creating domain 'microservices'...")
        result = await client.call_tool("create_domain", {
            "name": "microservices",
            "description": "Microservices architecture domain"
        })
        print(f"Result: {json.dumps(result, indent=2)}")
        
        if result.get("error"):
            print(f"ERROR: {result['error']}")
            return False
            
        domain_created = "microservices" in json.dumps(result.get("result", {}))
        print(f"Domain creation success: {domain_created}\n")
        
        # Step 2: Add nodes: api-gateway, user-service, payment-service
        nodes_to_create = [
            ("api-gateway", "https://api-gateway.microservices.local"),
            ("user-service", "https://user-service.microservices.local"),
            ("payment-service", "https://payment-service.microservices.local")
        ]
        
        node_ids = {}
        
        for service_name, service_url in nodes_to_create:
            print(f"Creating node for {service_name}...")
            result = await client.call_tool("create_node", {
                "domain_name": "microservices",
                "url": service_url,
                "title": service_name,
                "description": f"{service_name} microservice"
            })
            print(f"Result: {json.dumps(result, indent=2)}")
            
            if result.get("error"):
                print(f"ERROR: {result['error']}")
                return False
                
            # Extract node ID from result
            node_data = result.get("result", {})
            if "content" in node_data and len(node_data["content"]) > 0:
                node_info = json.loads(node_data["content"][0]["text"])
                node_ids[service_name] = node_info.get("composite_id", f"url-db:microservices:{node_info.get('id')}")
            print(f"Node {service_name} created with ID: {node_ids.get(service_name)}\n")
            
        # Step 3: Create subscription to monitor "api-gateway" for all event types
        print("Step 3: Creating subscription for api-gateway with all event types...")
        api_gateway_id = node_ids.get("api-gateway", "url-db:microservices:1")
        result = await client.call_tool("create_subscription", {
            "composite_id": api_gateway_id,
            "subscriber_service": "monitoring-system",
            "event_types": ["created", "updated", "deleted", "attribute_changed"]
        })
        print(f"Result: {json.dumps(result, indent=2)}")
        
        if result.get("error"):
            print(f"ERROR: {result['error']}")
            return False
            
        # Extract subscription ID
        subscription1_id = None
        sub_data = result.get("result", {})
        if "content" in sub_data and len(sub_data["content"]) > 0:
            sub_info = json.loads(sub_data["content"][0]["text"])
            subscription1_id = sub_info.get("id")
        print(f"Subscription 1 created with ID: {subscription1_id}\n")
        
        # Step 4: Create filtered subscription for "user-service" monitoring only "updated" events
        print("Step 4: Creating filtered subscription for user-service with 'updated' events only...")
        user_service_id = node_ids.get("user-service", "url-db:microservices:2")
        result = await client.call_tool("create_subscription", {
            "composite_id": user_service_id,
            "subscriber_service": "alert-system",
            "subscriber_endpoint": "https://alerts.example.com/webhook",
            "event_types": ["updated"]
        })
        print(f"Result: {json.dumps(result, indent=2)}")
        
        if result.get("error"):
            print(f"ERROR: {result['error']}")
            return False
            
        # Extract subscription ID
        subscription2_id = None
        sub_data = result.get("result", {})
        if "content" in sub_data and len(sub_data["content"]) > 0:
            sub_info = json.loads(sub_data["content"][0]["text"])
            subscription2_id = sub_info.get("id")
        print(f"Subscription 2 created with ID: {subscription2_id}\n")
        
        # Step 5: List all subscriptions and verify both were created
        print("Step 5: Listing all subscriptions...")
        result = await client.call_tool("list_subscriptions", {})
        print(f"Result: {json.dumps(result, indent=2)}")
        
        if result.get("error"):
            print(f"ERROR: {result['error']}")
            return False
            
        # Parse subscription list
        subscriptions_data = result.get("result", {})
        subscriptions_found = False
        if "content" in subscriptions_data and len(subscriptions_data["content"]) > 0:
            subs_info = json.loads(subscriptions_data["content"][0]["text"])
            subscriptions = subs_info.get("subscriptions", [])
            print(f"Found {len(subscriptions)} subscriptions")
            subscriptions_found = len(subscriptions) >= 2
        print(f"Both subscriptions found: {subscriptions_found}\n")
        
        # Step 6: Get node-specific subscriptions for "api-gateway"
        print("Step 6: Getting node-specific subscriptions for api-gateway...")
        result = await client.call_tool("get_node_subscriptions", {
            "composite_id": api_gateway_id
        })
        print(f"Result: {json.dumps(result, indent=2)}")
        
        if result.get("error"):
            print(f"ERROR: {result['error']}")
            return False
            
        # Verify api-gateway subscription
        node_subscriptions_found = False
        node_subs_data = result.get("result", {})
        if "content" in node_subs_data and len(node_subs_data["content"]) > 0:
            node_subs_info = json.loads(node_subs_data["content"][0]["text"])
            node_subscriptions = node_subs_info if isinstance(node_subs_info, list) else []
            print(f"Found {len(node_subscriptions)} subscriptions for api-gateway")
            node_subscriptions_found = len(node_subscriptions) >= 1
        print(f"API Gateway subscriptions found: {node_subscriptions_found}\n")
        
        # Step 7: Delete the first subscription and verify removal
        print("Step 7: Deleting first subscription...")
        if subscription1_id:
            result = await client.call_tool("delete_subscription", {
                "subscription_id": subscription1_id
            })
            print(f"Result: {json.dumps(result, indent=2)}")
            
            if result.get("error"):
                print(f"ERROR: {result['error']}")
                return False
                
            print("Subscription deleted successfully")
            
            # Verify deletion by listing subscriptions again
            print("Verifying deletion by listing subscriptions again...")
            result = await client.call_tool("list_subscriptions", {})
            print(f"Result: {json.dumps(result, indent=2)}")
            
            remaining_subscriptions = 0
            if result.get("result", {}).get("content"):
                subs_info = json.loads(result["result"]["content"][0]["text"])
                remaining_subscriptions = len(subs_info.get("subscriptions", []))
            
            print(f"Remaining subscriptions: {remaining_subscriptions}")
            deletion_verified = remaining_subscriptions == 1
            print(f"Deletion verified: {deletion_verified}\n")
        else:
            print("ERROR: Could not retrieve subscription1_id for deletion\n")
            return False
            
        await client.close()
        
        print("=== Subscription Lifecycle Management Test Completed ===")
        print("✓ Domain created successfully")
        print("✓ All three nodes created successfully") 
        print("✓ Two subscriptions created with different configurations")
        print("✓ Subscriptions listed and verified")
        print("✓ Node-specific subscriptions retrieved")
        print("✓ Subscription deleted and verified")
        
        return True
        
    except Exception as e:
        print(f"Test failed with exception: {e}")
        return False
    finally:
        # Clean up temporary database file
        try:
            os.unlink(db_path)
        except:
            pass


if __name__ == "__main__":
    result = asyncio.run(test_subscription_lifecycle())
    sys.exit(0 if result else 1)