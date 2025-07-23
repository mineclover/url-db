# MCP 도구 테스트 가이드

## 개요
이 문서는 URL-DB MCP 도구들의 테스트 방법과 절차를 설명합니다. **모든 MCP 도구는 이 가이드에 따라 체계적으로 테스트되어야 합니다.** 새로운 도구가 추가되거나 기존 도구가 수정될 때마다 이 가이드를 참조하여 전체 도구 세트를 테스트해야 합니다.

## 테스트 원칙

### 1. 체계적 접근
- **순차적 테스트**: CRUD 작업 순서에 따라 테스트
- **데이터 무결성**: 생성 → 조회 → 수정 → 삭제 순서 준수
- **정리 작업**: 테스트 후 데이터 정리 필수

### 2. 테스트 시나리오
- **기본 기능**: 각 도구의 기본 동작 확인
- **에러 처리**: 잘못된 입력에 대한 적절한 응답
- **경계 조건**: 빈 값, 최대값 등 경계 조건 테스트
- **통합 테스트**: 여러 도구 간 연동 테스트

## 테스트 절차

### Phase 1: 시스템 상태 확인
```bash
# 1. 서버 정보 조회
mcp_url-db_get_server_info

# 2. 기존 도메인 목록 확인
mcp_url-db_list_domains
```

### Phase 2: 도메인 관리 테스트
```bash
# 1. 도메인 생성
mcp_url-db_create_domain
- name: "test-domain"
- description: "테스트용 도메인"

# 2. 도메인 속성 생성
mcp_url-db_create_domain_attribute
- domain_name: "test-domain"
- name: "category"
- type: "tag"
- description: "카테고리 분류"

# 3. 도메인 속성 목록 조회
mcp_url-db_list_domain_attributes
- domain_name: "test-domain"

# 4. 도메인 속성 상세 조회
mcp_url-db_get_domain_attribute
- composite_id: "url-db:test-domain:attr-{id}"

# 5. 도메인 속성 업데이트
mcp_url-db_update_domain_attribute
- composite_id: "url-db:test-domain:attr-{id}"
- description: "업데이트된 설명"
```

### Phase 3: 노드 관리 테스트
```bash
# 1. 노드 생성
mcp_url-db_create_node
- domain_name: "test-domain"
- url: "https://example.com/test"
- title: "테스트 URL"
- description: "테스트용 URL"

# 2. 노드 목록 조회
mcp_url-db_list_nodes
- domain_name: "test-domain"
- page: 1
- size: 10

# 3. 특정 노드 조회
mcp_url-db_get_node
- composite_id: "url-db:test-domain:{id}"

# 4. URL로 노드 검색
mcp_url-db_find_node_by_url
- domain_name: "test-domain"
- url: "https://example.com/test"
```

### Phase 4: 속성 관리 테스트
```bash
# 1. 노드에 속성 설정
mcp_url-db_set_node_attributes
- composite_id: "url-db:test-domain:{id}"
- attributes: [{"name": "category", "value": "test"}]

# 2. 노드 속성 조회
mcp_url-db_get_node_attributes
- composite_id: "url-db:test-domain:{id}"

# 3. 속성으로 노드 필터링
mcp_url-db_filter_nodes_by_attributes
- domain_name: "test-domain"
- filters: [{"name": "category", "value": "test", "operator": "equals"}]

# 4. 노드와 속성 함께 조회
mcp_url-db_get_node_with_attributes
- composite_id: "url-db:test-domain:{id}"
```

### Phase 5: 업데이트 테스트
```bash
# 1. 노드 업데이트
mcp_url-db_update_node
- composite_id: "url-db:test-domain:{id}"
- title: "업데이트된 제목"
- description: "업데이트된 설명"
```

### Phase 6: 정리 작업
```bash
# 1. 노드 삭제
mcp_url-db_delete_node
- composite_id: "url-db:test-domain:{id}"

# 2. 도메인 속성 삭제
mcp_url-db_delete_domain_attribute
- composite_id: "url-db:test-domain:attr-{id}"

# 3. 정리 확인
mcp_url-db_list_nodes
mcp_url-db_list_domain_attributes
```

### Phase 7: 전체 도구 세트 검증
```bash
# 1. 등록된 모든 MCP 도구 확인
# 시스템에서 사용 가능한 모든 도구 목록을 확인하고 각각 테스트

# 2. 도구별 기본 기능 테스트
# 각 도구에 대해 다음 사항들을 확인:
# - 기본 파라미터로 정상 동작
# - 필수 파라미터 누락 시 적절한 에러 반환
# - 잘못된 파라미터 형식 시 적절한 에러 반환
# - 응답 형식이 예상과 일치

# 3. 도구 간 연동 테스트
# - 생성 도구 → 조회 도구 → 수정 도구 → 삭제 도구 순서로 테스트
# - 연관된 도구들 간의 데이터 일관성 확인
```

## 테스트 체크리스트

### ✅ 전체 도구 세트 테스트
- [ ] **등록된 모든 도구 확인** - 시스템에서 사용 가능한 모든 MCP 도구 목록 확인
- [ ] **도구별 기본 기능** - 각 도구의 기본 동작 테스트
- [ ] **도구별 에러 처리** - 각 도구의 에러 상황 처리 테스트
- [ ] **도구 간 연동** - 연관된 도구들 간의 데이터 일관성 테스트
- [ ] **응답 형식 검증** - 각 도구의 응답 형식이 예상과 일치하는지 확인

### ✅ 기본 기능 테스트
- [ ] 서버 정보 조회
- [ ] 도메인 생성/조회
- [ ] 노드 생성/조회/수정/삭제
- [ ] 속성 생성/설정/조회
- [ ] 필터링 기능
- [ ] 검색 기능

### ✅ 에러 처리 테스트
- [ ] 필수 파라미터 누락 시 에러
- [ ] 잘못된 형식의 입력 시 에러
- [ ] 존재하지 않는 리소스 조회 시 에러
- [ ] 중복 데이터 생성 시 에러

### ✅ 데이터 무결성 테스트
- [ ] 생성된 데이터가 정상적으로 조회됨
- [ ] 업데이트된 데이터가 반영됨
- [ ] 삭제된 데이터가 목록에서 제외됨
- [ ] 연관 데이터의 일관성 유지

### ✅ 성능 테스트
- [ ] 대량 데이터 처리 시 응답 시간
- [ ] 페이지네이션 정상 작동
- [ ] 검색 기능의 정확성

## 테스트 데이터 관리

### 테스트 도메인 명명 규칙
- **형식**: `mcp-test-{기능명}-{날짜}`
- **예시**: `mcp-test-attributes-20240723`

### 테스트 URL 패턴
- **기본**: `https://example.com/test`
- **속성 테스트**: `https://example.com/test-attributes`
- **필터링 테스트**: `https://example.com/test-filter`

### 테스트 속성 패턴
- **카테고리**: `category`
- **태그**: `tags`
- **우선순위**: `priority`
- **상태**: `status`

### 범용 테스트 데이터 규칙
- **도메인**: `test-domain-{timestamp}`
- **노드**: `test-node-{timestamp}`
- **속성**: `test-attribute-{timestamp}`
- **URL**: `https://test.example.com/{timestamp}`

## 에러 시나리오 테스트

### 1. 잘못된 형식의 입력
```bash
# 잘못된 Composite ID
- composite_id: "invalid:format:id"

# 존재하지 않는 도메인
- domain_name: "non-existent-domain"

# 잘못된 속성 타입
- type: "invalid_type"
```

### 2. 필수 파라미터 누락
```bash
# 필수 파라미터 없이 도구 호출
# 예: domain_name, url, composite_id 등
```

### 3. 잘못된 데이터 타입
```bash
# 문자열이 필요한 곳에 숫자 입력
# 숫자가 필요한 곳에 문자열 입력
# 배열이 필요한 곳에 단일 값 입력
```

### 4. 경계 조건 테스트
```bash
# 빈 문자열
# 매우 긴 문자열
# 특수 문자 포함
# 최대값/최소값 테스트
```

## 성능 테스트 시나리오

### 1. 대량 데이터 테스트
```bash
# 대량 데이터 생성
# - 100개 이상의 노드 생성
# - 다양한 속성 조합으로 테스트
# - 페이지네이션 정상 작동 확인

# 페이지네이션 테스트
# - 다양한 page, size 조합
# - 경계값 테스트 (page=0, size=0 등)
```

### 2. 복잡한 쿼리 테스트
```bash
# 복잡한 필터링 조건
# - 다중 조건 필터링
# - 다양한 연산자 조합
# - 정렬 기능 테스트

# 검색 기능 테스트
# - 부분 문자열 검색
# - 대소문자 구분
# - 특수 문자 처리
```

## 테스트 결과 문서화

### 성공 케이스
- ✅ 도구명: 성공한 기능
- 📝 응답 예시: 실제 응답 데이터
- ⏱️ 응답 시간: 성능 측정 결과

### 실패 케이스
- ❌ 도구명: 실패한 기능
- 🐛 에러 메시지: 실제 에러 응답
- 🔧 해결 방법: 문제 해결 방법

## 자동화 테스트 고려사항

### 1. 테스트 스크립트 작성
```bash
#!/bin/bash
# mcp_test_runner.sh

echo "=== MCP 도구 테스트 시작 ==="

# Phase 1: 시스템 확인
echo "1. 서버 정보 확인"
# 서버 정보 조회 도구 호출

# Phase 2: 전체 도구 목록 확인
echo "2. 등록된 도구 목록 확인"
# 시스템에서 사용 가능한 모든 도구 확인

# Phase 3: 도구별 테스트
echo "3. 각 도구별 기본 기능 테스트"
# 각 도구에 대해 기본 파라미터로 테스트

# Phase 4: 에러 처리 테스트
echo "4. 에러 상황 테스트"
# 필수 파라미터 누락, 잘못된 형식 등 테스트

# Phase 5: 정리 작업
echo "5. 테스트 데이터 정리"
# 생성된 테스트 데이터 삭제
```

### 2. CI/CD 통합
- GitHub Actions에서 자동 테스트 실행
- 테스트 결과 리포트 생성
- 실패 시 알림 발송

## 유지보수 고려사항

### 1. 정기적 테스트
- **주기**: 새로운 기능 추가 시마다
- **범위**: 전체 도구 세트
- **결과**: 문서 업데이트

### 2. 버전 관리
- 테스트 시나리오 버전 관리
- 호환성 테스트 추가
- 하위 호환성 보장

### 3. 문서 업데이트
- 새로운 도구 추가 시 가이드 업데이트
- 에러 케이스 추가
- 성능 개선 사항 반영

---

## 결론

이 테스트 가이드를 따라 체계적으로 **모든 MCP 도구들을 테스트**하면:

1. **품질 보장**: 모든 기능이 정상 작동함을 확인
2. **안정성**: 예상치 못한 에러 상황 대비
3. **유지보수성**: 문제 발생 시 빠른 진단 가능
4. **문서화**: 팀원들이 쉽게 테스트 수행 가능
5. **완전성**: 모든 등록된 도구의 기능 검증

**중요**: 새로운 도구가 추가되거나 기존 도구가 수정될 때마다 반드시 전체 도구 세트를 테스트해야 합니다. 테스트는 개발 과정의 핵심 부분이며, 이 가이드를 통해 일관되고 신뢰할 수 있는 테스트를 수행할 수 있습니다. 