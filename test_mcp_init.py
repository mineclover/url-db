#!/usr/bin/env python3
"""
MCP 초기화 및 기능 테스트 스크립트
"""
import json
import subprocess
import sys

def send_mcp_request(request):
    """MCP 요청을 전송하는 함수"""
    try:
        result = subprocess.run(
            ["./bin/url-db", "-mcp-mode=stdio"],
            input=json.dumps(request) + "\n",
            capture_output=True,
            text=True,
            timeout=10
        )
        
        if result.stdout:
            lines = result.stdout.strip().split('\n')
            for line in lines:
                if line.strip():
                    try:
                        return json.loads(line)
                    except json.JSONDecodeError:
                        continue
        return {"error": "No valid JSON response"}
    except Exception as e:
        return {"error": str(e)}

def main():
    print("🚀 MCP 서버 초기화 및 테스트...")
    
    # 1. MCP 서버 초기화
    print("\n1️⃣ MCP 서버 초기화...")
    init_request = {
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
    }
    
    result = send_mcp_request(init_request)
    print(f"초기화 응답: {json.dumps(result, indent=2)}")
    
    # 2. 도구 목록 조회
    print("\n2️⃣ 도구 목록 조회...")
    tools_request = {
        "jsonrpc": "2.0",
        "id": 2,
        "method": "tools/list",
        "params": {}
    }
    
    result = send_mcp_request(tools_request)
    if "result" in result and "tools" in result["result"]:
        tools = result["result"]["tools"]
        print(f"✅ 사용 가능한 도구 수: {len(tools)}")
        for tool in tools[:5]:  # 처음 5개만 출력
            print(f"  - {tool.get('name', 'Unknown')}: {tool.get('description', 'No description')}")
    else:
        print(f"❌ 도구 목록 조회 실패: {result}")
    
    # 3. 서버 정보 조회
    print("\n3️⃣ 서버 정보 조회...")
    server_info_request = {
        "jsonrpc": "2.0",
        "id": 3,
        "method": "tools/call",
        "params": {
            "name": "get_server_info",
            "arguments": {}
        }
    }
    
    result = send_mcp_request(server_info_request)
    if "result" in result:
        print("✅ get_server_info 성공")
        if "content" in result["result"]:
            for content in result["result"]["content"]:
                if content.get("type") == "text":
                    try:
                        data = json.loads(content["text"])
                        print(f"  서버명: {data.get('name', 'Unknown')}")
                        print(f"  버전: {data.get('version', 'Unknown')}")
                    except:
                        print(f"  응답: {content['text']}")
    else:
        print(f"❌ get_server_info 실패: {result}")

    print("\n🎉 MCP 초기화 테스트 완료!")

if __name__ == "__main__":
    main()