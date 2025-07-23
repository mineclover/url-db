#!/usr/bin/env python3
"""
MCP 핵심 기능 테스트 스크립트
"""
import json
import subprocess
import sys

def test_mcp_tool(tool_name, params=None):
    """MCP 도구를 테스트하는 함수"""
    if params is None:
        params = {}
    
    request = {
        "jsonrpc": "2.0",
        "id": 1,
        "method": "tools/call",
        "params": {
            "name": tool_name,
            "arguments": params
        }
    }
    
    try:
        # MCP 서버에 요청 전송
        result = subprocess.run(
            ["./bin/url-db", "-mcp-mode=stdio"],
            input=json.dumps(request) + "\n",
            capture_output=True,
            text=True,
            timeout=5
        )
        
        if result.stdout:
            response = json.loads(result.stdout.strip())
            return response
        else:
            return {"error": f"No output from {tool_name}"}
            
    except subprocess.TimeoutExpired:
        return {"error": f"Timeout for {tool_name}"}
    except json.JSONDecodeError as e:
        return {"error": f"JSON decode error for {tool_name}: {e}"}
    except Exception as e:
        return {"error": f"Error testing {tool_name}: {e}"}

def main():
    print("🚀 MCP 핵심 기능 테스트 시작...")
    
    # 1. 서버 정보 조회
    print("\n1️⃣ 서버 정보 조회...")
    result = test_mcp_tool("get_server_info")
    if "error" not in result:
        print("✅ get_server_info 성공")
    else:
        print(f"❌ get_server_info 실패: {result}")
    
    # 2. 도메인 생성
    print("\n2️⃣ 도메인 생성...")
    result = test_mcp_tool("create_domain", {
        "name": "test-domain",
        "description": "테스트 도메인"
    })
    if "error" not in result:
        print("✅ create_domain 성공")
    else:
        print(f"❌ create_domain 실패: {result}")
    
    # 3. 도메인 목록 조회
    print("\n3️⃣ 도메인 목록 조회...")
    result = test_mcp_tool("list_domains")
    if "error" not in result:
        print("✅ list_domains 성공")
    else:
        print(f"❌ list_domains 실패: {result}")
    
    # 4. 노드 생성
    print("\n4️⃣ 노드 생성...")
    result = test_mcp_tool("create_node", {
        "domain_name": "test-domain",
        "url": "https://example.com",
        "title": "예제 사이트",
        "description": "테스트용 URL"
    })
    if "error" not in result:
        print("✅ create_node 성공")
        node_id = "url-db:test-domain:1"  # 첫 번째 노드 ID
    else:
        print(f"❌ create_node 실패: {result}")
        return
    
    # 5. 노드 조회
    print("\n5️⃣ 노드 조회...")
    result = test_mcp_tool("get_node", {
        "composite_id": node_id
    })
    if "error" not in result:
        print("✅ get_node 성공")
    else:
        print(f"❌ get_node 실패: {result}")
    
    # 6. 노드 목록 조회
    print("\n6️⃣ 노드 목록 조회...")
    result = test_mcp_tool("list_nodes", {
        "domain_name": "test-domain"
    })
    if "error" not in result:
        print("✅ list_nodes 성공")
    else:
        print(f"❌ list_nodes 실패: {result}")
    
    print("\n🎉 MCP 핵심 기능 테스트 완료!")

if __name__ == "__main__":
    main()