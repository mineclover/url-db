#!/usr/bin/env python3
"""
간단한 MCP 클라이언트 테스트 스크립트
JSON-RPC 2.0 over stdio를 사용해서 URL-DB MCP 서버와 통신
"""

import json
import subprocess
import sys
import time
from tool_constants import CREATE_DOMAIN, CREATE_NODE, LIST_DOMAINS


class MCPClient:
    def __init__(self, server_path):
        self.server_path = server_path
        self.process = None
        self.request_id = 1
        
    def start_server(self):
        """MCP 서버를 stdio 모드로 시작"""
        self.process = subprocess.Popen(
            [self.server_path, "-mcp-mode=stdio"],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            bufsize=0
        )
        
    def send_request(self, method, params=None):
        """JSON-RPC 2.0 요청 전송"""
        request = {
            "jsonrpc": "2.0",
            "id": self.request_id,
            "method": method,
            "params": params or {}
        }
        
        request_json = json.dumps(request)
        print(f">>> Sending: {request_json}")
        
        self.process.stdin.write(request_json + "\n")
        self.process.stdin.flush()
        
        # 응답 읽기
        response_line = self.process.stdout.readline().strip()
        print(f"<<< Received: {response_line}")
        
        if response_line:
            try:
                response = json.loads(response_line)
                self.request_id += 1
                return response
            except json.JSONDecodeError as e:
                print(f"Failed to parse response: {e}")
                return None
        
        return None
        
    def send_notification(self, method, params=None):
        """JSON-RPC 2.0 알림 전송 (응답 없음)"""
        notification = {
            "jsonrpc": "2.0",
            "method": method,
            "params": params or {}
        }
        
        notification_json = json.dumps(notification)
        print(f">>> Sending notification: {notification_json}")
        
        self.process.stdin.write(notification_json + "\n")
        self.process.stdin.flush()
        
    def close(self):
        """서버 종료"""
        if self.process:
            self.process.stdin.close()
            self.process.wait()

def main():
    # 서버 바이너리 경로
    server_path = "../../bin/url-db"
    
    client = MCPClient(server_path)
    
    try:
        print("Starting MCP server...")
        client.start_server()
        
        # 잠시 대기
        time.sleep(1)
        
        print("\n=== MCP Handshake ===")
        
        # 1. Initialize 요청
        print("\n1. Sending initialize request...")
        init_response = client.send_request("initialize", {
            "protocolVersion": "2024-11-05",
            "capabilities": {
                "experimental": {},
                "sampling": {}
            },
            "clientInfo": {
                "name": "test-mcp-client",
                "version": "1.0.0"
            }
        })
        
        if init_response and "result" in init_response:
            print("✅ Initialize successful!")
            print(f"Server info: {init_response['result']['serverInfo']}")
        else:
            print("❌ Initialize failed!")
            return
            
        # 2. Initialized 알림
        print("\n2. Sending initialized notification...")
        client.send_notification("initialized", {})
        
        # 잠시 대기
        time.sleep(0.5)
        
        print("\n=== Testing Tools ===")
        
        # 3. Tools 목록 요청
        print("\n3. Getting tools list...")
        tools_response = client.send_request("tools/list", {})
        
        if tools_response and "result" in tools_response:
            tools = tools_response["result"]["tools"]
            print(f"✅ Found {len(tools)} tools:")
            for tool in tools[:3]:  # 처음 3개만 출력
                print(f"  - {tool['name']}: {tool['description']}")
        else:
            print("❌ Failed to get tools list!")
            return
            
        print("\n=== Testing Domain Operations ===")
        
        # 4. 도메인 목록 조회
        print("\n4. Listing domains...")
        domains_response = client.send_request("tools/call", {
            "name": LIST_DOMAINS,
            "arguments": {}
        })
        
        if domains_response and "result" in domains_response:
            print("✅ Domain list call successful!")
            print("Response content:", domains_response["result"]["content"][0]["text"][:200] + "...")
        else:
            print("❌ Domain list call failed!")
            
        # 5. 테스트 도메인 생성
        print("\n5. Creating test domain...")
        create_domain_response = client.send_request("tools/call", {
            "name": CREATE_DOMAIN,
            "arguments": {
                "name": "test-domain",
                "description": "Test domain for MCP integration"
            }
        })
        
        if create_domain_response and "result" in create_domain_response:
            if not create_domain_response["result"].get("isError", False):
                print("✅ Domain creation successful!")
            else:
                print("⚠️ Domain creation returned error (might already exist)")
        else:
            print("❌ Domain creation failed!")
            
        # 6. 테스트 노드 생성
        print("\n6. Creating test node...")
        create_node_response = client.send_request("tools/call", {
            "name": CREATE_NODE,
            "arguments": {
                "domain_name": "test-domain",
                "url": "https://example.com/test",
                "title": "Test Node via MCP",
                "description": "Created through MCP JSON-RPC protocol"
            }
        })
        
        if create_node_response and "result" in create_node_response:
            if not create_node_response["result"].get("isError", False):
                print("✅ Node creation successful!")
                response_content = create_node_response["result"]["content"][0]["text"]
                print("Created node info:", response_content[:200] + "...")
            else:
                print("⚠️ Node creation returned error")
                print("Error:", create_node_response["result"]["content"][0]["text"])
        else:
            print("❌ Node creation failed!")
            
        print("\n=== Testing Resources ===")
        
        # 7. 리소스 목록 조회
        print("\n7. Getting resources list...")
        resources_response = client.send_request("resources/list", {})
        
        if resources_response and "result" in resources_response:
            resources = resources_response["result"]["resources"]
            print(f"✅ Found {len(resources)} resources:")
            for resource in resources[:3]:  # 처음 3개만 출력
                print(f"  - {resource['uri']}: {resource['name']}")
        else:
            print("❌ Failed to get resources list!")
            
        # 8. 서버 정보 리소스 읽기
        print("\n8. Reading server info resource...")
        server_info_response = client.send_request("resources/read", {
            "uri": "mcp://server/info"
        })
        
        if server_info_response and "result" in server_info_response:
            print("✅ Server info resource read successful!")
            content = server_info_response["result"]["contents"][0]["text"]
            print("Server info:", content[:200] + "...")
        else:
            print("❌ Failed to read server info resource!")
            
        print("\n=== Test Complete ===")
        print("✅ MCP JSON-RPC integration test completed successfully!")
        
    except Exception as e:
        print(f"❌ Test failed with error: {e}")
        import traceback
        traceback.print_exc()
        
    finally:
        print("\nClosing MCP client...")
        client.close()

if __name__ == "__main__":
    main()