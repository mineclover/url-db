#!/usr/bin/env python3
"""
Scenario 1.3: Automatic Event Generation and Tracking
From MCP LLM judge scenarios file
"""

import json
import subprocess
import sys
import time

def send_mcp_request(method, params=None):
    """Send MCP JSON-RPC 2.0 request"""
    request = {
        "jsonrpc": "2.0",
        "id": 1,
        "method": method,
        "params": params or {}
    }
    
    try:
        # Send request to MCP server via stdin
        process = subprocess.Popen(
            ["./bin/url-db", "-mcp-mode=stdio", "-db-path=test_external_deps.db"],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )
        
        stdout, stderr = process.communicate(json.dumps(request) + "\n")
        
        if stderr:
            print(f"STDERR: {stderr}")
        
        if stdout:
            return json.loads(stdout.strip())
        
    except Exception as e:
        print(f"Error sending request: {e}")
        return None

def test_event_generation():
    """Execute Scenario 1.3: Automatic Event Generation and Tracking"""
    print("=== Scenario 1.3: Automatic Event Generation and Tracking ===")
    
    results = []
    
    # Step 1: Create new node "database-service" in "microservices" domain
    print("\n1. Creating database-service node...")
    response = send_mcp_request("create_node", {
        "domain_name": "microservices",
        "url": "https://example.com/database-service",
        "title": "Database Service",
        "description": "Database service for the microservices architecture"
    })
    
    if response and "result" in response:
        node_id = response["result"].get("composite_id")
        print(f"✅ Node created: {node_id}")
        results.append(("Create Node", True, node_id))
    else:
        print(f"❌ Failed to create node: {response}")
        results.append(("Create Node", False, str(response)))
        return results
    
    # Step 2: Immediately check for "node.created" event
    print("\n2. Checking for node.created event...")
    response = send_mcp_request("get_node_events", {
        "composite_id": node_id
    })
    
    if response and "result" in response:
        events = response["result"].get("events", [])
        created_events = [e for e in events if e.get("event_type") == "node.created"]
        if created_events:
            print(f"✅ Found node.created event: {created_events[0]}")
            results.append(("Node Created Event", True, f"Found {len(created_events)} created events"))
        else:
            print(f"❌ No node.created event found. Events: {events}")
            results.append(("Node Created Event", False, f"No created events in {len(events)} total events"))
    else:
        print(f"❌ Failed to get events: {response}")
        results.append(("Node Created Event", False, str(response)))
    
    # Step 3: Update the node title to "Primary Database Service"
    print("\n3. Updating node title...")
    response = send_mcp_request("update_node", {
        "composite_id": node_id,
        "title": "Primary Database Service"
    })
    
    if response and "result" in response:
        print("✅ Node title updated")
        results.append(("Update Node Title", True, "Title updated to Primary Database Service"))
    else:
        print(f"❌ Failed to update node: {response}")
        results.append(("Update Node Title", False, str(response)))
    
    # Step 4: Check for "node.updated" event with before/after data
    print("\n4. Checking for node.updated event after title change...")
    response = send_mcp_request("get_node_events", {
        "composite_id": node_id
    })
    
    if response and "result" in response:
        events = response["result"].get("events", [])
        updated_events = [e for e in events if e.get("event_type") == "node.updated"]
        if updated_events:
            latest_update = updated_events[-1]  # Get most recent update
            print(f"✅ Found node.updated event: {latest_update}")
            results.append(("Node Updated Event 1", True, f"Found update event with before/after data"))
        else:
            print(f"❌ No node.updated event found. Events: {events}")
            results.append(("Node Updated Event 1", False, f"No updated events found"))
    else:
        print(f"❌ Failed to get events: {response}")
        results.append(("Node Updated Event 1", False, str(response)))
    
    # Step 5: Update description to "Main PostgreSQL database instance"
    print("\n5. Updating node description...")
    response = send_mcp_request("update_node", {
        "composite_id": node_id,
        "description": "Main PostgreSQL database instance"
    })
    
    if response and "result" in response:
        print("✅ Node description updated")
        results.append(("Update Node Description", True, "Description updated"))
    else:
        print(f"❌ Failed to update description: {response}")
        results.append(("Update Node Description", False, str(response)))
    
    # Step 6: Verify second "node.updated" event was created
    print("\n6. Checking for second node.updated event...")
    response = send_mcp_request("get_node_events", {
        "composite_id": node_id
    })
    
    if response and "result" in response:
        events = response["result"].get("events", [])
        updated_events = [e for e in events if e.get("event_type") == "node.updated"]
        if len(updated_events) >= 2:
            print(f"✅ Found {len(updated_events)} node.updated events")
            results.append(("Node Updated Event 2", True, f"Found {len(updated_events)} update events"))
        else:
            print(f"❌ Expected at least 2 update events, found {len(updated_events)}")
            results.append(("Node Updated Event 2", False, f"Only {len(updated_events)} update events"))
    else:
        print(f"❌ Failed to get events: {response}")
        results.append(("Node Updated Event 2", False, str(response)))
    
    # Step 7: Get all pending events to see unprocessed events
    print("\n7. Getting all pending events...")
    response = send_mcp_request("get_pending_events", {})
    
    if response and "result" in response:
        pending_events = response["result"].get("events", [])
        print(f"✅ Found {len(pending_events)} pending events")
        results.append(("Pending Events", True, f"Found {len(pending_events)} pending events"))
        
        # Step 8: Process one event and verify it's marked as processed
        if pending_events:
            print("\n8. Processing one event...")
            event_id = pending_events[0].get("id")
            response = send_mcp_request("process_event", {
                "event_id": event_id
            })
            
            if response and "result" in response:
                print(f"✅ Event {event_id} processed")
                results.append(("Process Event", True, f"Event {event_id} processed"))
            else:
                print(f"❌ Failed to process event: {response}")
                results.append(("Process Event", False, str(response)))
        else:
            print("⚠️ No pending events to process")
            results.append(("Process Event", False, "No pending events available"))
    else:
        print(f"❌ Failed to get pending events: {response}")
        results.append(("Pending Events", False, str(response)))
    
    # Step 9: Delete the node and check for "node.deleted" event
    print("\n9. Deleting node and checking for delete event...")
    response = send_mcp_request("delete_node", {
        "composite_id": node_id
    })
    
    if response and "result" in response:
        print("✅ Node deleted")
        results.append(("Delete Node", True, "Node successfully deleted"))
        
        # Check for delete event
        response = send_mcp_request("get_node_events", {
            "composite_id": node_id
        })
        
        if response and "result" in response:
            events = response["result"].get("events", [])
            deleted_events = [e for e in events if e.get("event_type") == "node.deleted"]
            if deleted_events:
                print(f"✅ Found node.deleted event: {deleted_events[0]}")
                results.append(("Node Deleted Event", True, "Delete event found"))
            else:
                print(f"❌ No delete event found. Events: {events}")
                results.append(("Node Deleted Event", False, "No delete event found"))
        else:
            print(f"❌ Failed to get events after deletion: {response}")
            results.append(("Node Deleted Event", False, "Could not retrieve events after deletion"))
    else:
        print(f"❌ Failed to delete node: {response}")
        results.append(("Delete Node", False, str(response)))
    
    return results

def print_results(results):
    """Print test results summary"""
    print("\n" + "="*60)
    print("EVENT GENERATION AND TRACKING TEST RESULTS")
    print("="*60)
    
    passed = sum(1 for _, success, _ in results if success)
    total = len(results)
    
    for test_name, success, details in results:
        status = "✅ PASSED" if success else "❌ FAILED"
        print(f"{status}: {test_name}")
        if details and not success:
            print(f"         Details: {details}")
    
    print(f"\nTEST SUMMARY: {passed}/{total} tests passed")
    
    if passed == total:
        print("🎉 ALL TESTS PASSED - Event system working correctly!")
    else:
        print("⚠️ SOME TESTS FAILED - Event system needs attention")

if __name__ == "__main__":
    results = test_event_generation()
    print_results(results)