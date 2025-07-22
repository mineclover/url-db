#!/usr/bin/env python3
"""
Debug domain creation issue
"""

import json
import subprocess
import sys
import time

def debug_domain_creation():
    server_path = "./bin/url-db"
    
    process = subprocess.Popen(
        [server_path, "-mcp-mode=stdio"],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
        bufsize=0
    )
    
    def send_request(method, params=None, request_id=1):
        request = {
            "jsonrpc": "2.0",
            "id": request_id,
            "method": method,
            "params": params or {}
        }
        
        request_json = json.dumps(request)
        print(f">>> Sending: {request_json}")
        
        process.stdin.write(request_json + "\n")
        process.stdin.flush()
        
        response_line = process.stdout.readline().strip()
        print(f"<<< Received: {response_line}")
        
        if response_line:
            try:
                return json.loads(response_line)
            except json.JSONDecodeError as e:
                print(f"JSON decode error: {e}")
                return None
        return None
    
    try:
        # Initialize
        print("=== Initializing ===")
        init_response = send_request("initialize", {
            "protocolVersion": "2024-11-05",
            "capabilities": {},
            "clientInfo": {"name": "debug", "version": "1.0.0"}
        })
        
        if init_response and "result" in init_response:
            print("✓ Initialize successful")
        else:
            print("✗ Initialize failed")
            return
            
        # Send initialized
        print("\n=== Sending initialized ===")
        notification = {"jsonrpc": "2.0", "method": "initialized", "params": {}}
        notification_json = json.dumps(notification)
        print(f">>> Sending: {notification_json}")
        process.stdin.write(notification_json + "\n")
        process.stdin.flush()
        time.sleep(0.2)
        
        # List existing domains
        print("\n=== Listing existing domains ===")
        domains_response = send_request("tools/call", {
            "name": "list_mcp_domains",
            "arguments": {}
        }, 2)
        
        if domains_response:
            print("Raw domain response:", domains_response)
            if "result" in domains_response and not domains_response["result"].get("isError", False):
                content = json.loads(domains_response["result"]["content"][0]["text"])
                print(f"Existing domains: {content}")
            else:
                print("Error in domain listing:", domains_response["result"])
        
        # Try to create a domain
        print("\n=== Creating test domain ===")
        create_response = send_request("tools/call", {
            "name": "create_mcp_domain",
            "arguments": {
                "name": "debug-test-domain",
                "description": "Debug test domain"
            }
        }, 3)
        
        if create_response:
            print("Raw create response:", create_response)
            if "result" in create_response:
                if create_response["result"].get("isError", False):
                    print("Domain creation error:", create_response["result"]["content"][0]["text"])
                else:
                    print("Domain created successfully:", create_response["result"]["content"][0]["text"])
            else:
                print("No result in create response")
        
    finally:
        process.stdin.close()
        process.wait()

if __name__ == "__main__":
    debug_domain_creation()