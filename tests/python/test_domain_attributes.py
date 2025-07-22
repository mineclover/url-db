#!/usr/bin/env python3
"""
도메인 속성 관리 기능 테스트 스크립트
"""

import json
import subprocess
import sys
import time

def send_mcp_request(method, params=None):
    """MCP 요청을 보내고 응답을 받습니다."""
    request = {
        "jsonrpc": "2.0",
        "id": 1,
        "method": method,
        "params": params or {}
    }
    
    # 서버에 요청 전송
    try:
        result = subprocess.run(
            ["./bin/url-db", "-mcp-mode=stdio", "-db-path=test.db", "-tool-name=test-tool"],
            input=json.dumps(request) + "\n",
            text=True,
            capture_output=True,
            timeout=10
        )
        
        if result.returncode != 0:
            print(f"서버 실행 오류: {result.stderr}")
            return None
            
        # 응답 파싱
        try:
            response = json.loads(result.stdout.strip())
            return response
        except json.JSONDecodeError as e:
            print(f"응답 파싱 오류: {e}")
            print(f"원본 응답: {result.stdout}")
            return None
            
    except subprocess.TimeoutExpired:
        print("요청 시간 초과")
        return None
    except Exception as e:
        print(f"요청 오류: {e}")
        return None

def test_domain_attributes():
    """도메인 속성 관리 기능을 테스트합니다."""
    print("=== 도메인 속성 관리 기능 테스트 ===\n")
    
    # 1. 도메인 생성
    print("1. 도메인 생성...")
    response = send_mcp_request("create_mcp_domain", {
        "name": "test-projects",
        "description": "테스트 프로젝트 도메인"
    })
    
    if response and "result" in response:
        print("✅ 도메인 생성 성공")
    else:
        print("❌ 도메인 생성 실패")
        print(f"응답: {response}")
        return
    
    # 2. 도메인 속성 생성
    print("\n2. 도메인 속성 생성...")
    response = send_mcp_request("create_mcp_domain_attribute", {
        "domain_name": "test-projects",
        "name": "status",
        "type": "tag",
        "description": "프로젝트 상태"
    })
    
    if response and "result" in response:
        print("✅ 속성 생성 성공")
        attribute_id = response["result"]["composite_id"]
        print(f"생성된 속성 ID: {attribute_id}")
    else:
        print("❌ 속성 생성 실패")
        print(f"응답: {response}")
        return
    
    # 3. 도메인 속성 목록 조회
    print("\n3. 도메인 속성 목록 조회...")
    response = send_mcp_request("list_mcp_domain_attributes", {
        "domain_name": "test-projects"
    })
    
    if response and "result" in response:
        print("✅ 속성 목록 조회 성공")
        attributes = response["result"]["attributes"]
        print(f"속성 개수: {len(attributes)}")
        for attr in attributes:
            print(f"  - {attr['name']} ({attr['type']}): {attr['description']}")
    else:
        print("❌ 속성 목록 조회 실패")
        print(f"응답: {response}")
        return
    
    # 4. 개별 속성 조회
    print("\n4. 개별 속성 조회...")
    response = send_mcp_request("get_mcp_domain_attribute", {
        "composite_id": attribute_id
    })
    
    if response and "result" in response:
        print("✅ 개별 속성 조회 성공")
        attr = response["result"]
        print(f"  - 이름: {attr['name']}")
        print(f"  - 타입: {attr['type']}")
        print(f"  - 설명: {attr['description']}")
    else:
        print("❌ 개별 속성 조회 실패")
        print(f"응답: {response}")
        return
    
    # 5. 속성 업데이트
    print("\n5. 속성 업데이트...")
    response = send_mcp_request("update_mcp_domain_attribute", {
        "composite_id": attribute_id,
        "description": "업데이트된 프로젝트 상태 설명"
    })
    
    if response and "result" in response:
        print("✅ 속성 업데이트 성공")
        attr = response["result"]
        print(f"  - 업데이트된 설명: {attr['description']}")
    else:
        print("❌ 속성 업데이트 실패")
        print(f"응답: {response}")
        return
    
    # 6. 추가 속성 생성
    print("\n6. 추가 속성 생성...")
    response = send_mcp_request("create_mcp_domain_attribute", {
        "domain_name": "test-projects",
        "name": "priority",
        "type": "ordered_tag",
        "description": "프로젝트 우선순위"
    })
    
    if response and "result" in response:
        print("✅ 추가 속성 생성 성공")
        priority_id = response["result"]["composite_id"]
    else:
        print("❌ 추가 속성 생성 실패")
        print(f"응답: {response}")
        return
    
    # 7. 최종 속성 목록 확인
    print("\n7. 최종 속성 목록 확인...")
    response = send_mcp_request("list_mcp_domain_attributes", {
        "domain_name": "test-projects"
    })
    
    if response and "result" in response:
        print("✅ 최종 속성 목록 조회 성공")
        attributes = response["result"]["attributes"]
        print(f"총 속성 개수: {len(attributes)}")
        for attr in attributes:
            print(f"  - {attr['name']} ({attr['type']}): {attr['description']}")
    else:
        print("❌ 최종 속성 목록 조회 실패")
        print(f"응답: {response}")
        return
    
    print("\n=== 모든 테스트 완료 ===")
    print("✅ 도메인 속성 관리 기능이 정상적으로 작동합니다!")

if __name__ == "__main__":
    test_domain_attributes() 