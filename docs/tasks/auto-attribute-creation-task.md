# 속성 자동 생성 기능 구현 계획 ✅ COMPLETED

## Status: ✅ COMPLETED (2025-07-23)

## 개요

현재 MCP 서비스에서 존재하지 않는 속성을 사용할 때 `VALIDATION_ERROR: attribute 'status' not found` 에러가 발생합니다. 이를 개선하기 위해 속성 자동 생성 옵션을 구현합니다.

## 현재 문제점

### 에러 시나리오
```bash
# 현재 발생하는 에러
Error setting node attributes: VALIDATION_ERROR: attribute 'status' not found
```

### 사용자 경험 문제
1. 사용자가 속성 이름을 정확히 알아야 함
2. 속성 생성과 할당을 두 단계로 나누어야 함
3. 실패 시 수동으로 속성을 먼저 생성해야 함

## 구현 목표

### 1. 자동 생성 옵션 추가
- `set_node_attributes` 호출 시 `auto_create_attributes` 파라미터 추가
- 기본값: `false` (기존 동작 유지)
- `true`로 설정 시 존재하지 않는 속성을 자동 생성

### 2. 속성 타입 추론
- 자동 생성 시 적절한 속성 타입 결정
- 기본값: `tag` 타입
- 향후 확장 가능한 구조

### 3. 사용자 친화적 에러 메시지
- 자동 생성이 비활성화된 경우 명확한 안내 메시지
- 속성 생성 방법 안내

## 구현 계획 ✅ COMPLETED

### Phase 1: MCP 인터페이스 확장 ✅ COMPLETED

#### 1.1 `set_node_attributes` 파라미터 추가 ✅ COMPLETED
```json
{
  "composite_id": "url-db:domain:node-id",
  "attributes": [
    {
      "name": "status",
      "value": "Production"
    }
  ],
  "auto_create_attributes": true
}
```

#### 1.2 새로운 MCP 도구 추가
```json
{
  "name": "set_node_attributes_with_auto_create",
  "description": "Set node attributes with automatic attribute creation",
  "inputSchema": {
    "type": "object",
    "properties": {
      "composite_id": {"type": "string"},
      "attributes": {"type": "array"},
      "auto_create_attributes": {"type": "boolean", "default": true}
    }
  }
}
```

### Phase 2: 백엔드 로직 구현

#### 2.1 속성 자동 생성 서비스
```go
// internal/services/auto_attribute_service.go
type AutoAttributeService struct {
    domainRepo    domains.DomainRepository
    attributeRepo attributes.AttributeRepository
    logger        *log.Logger
}

func (s *AutoAttributeService) SetNodeAttributesWithAutoCreate(
    ctx context.Context,
    compositeID string,
    attributes []models.NodeAttributeRequest,
    autoCreate bool,
) (*models.NodeAttributesResponse, error) {
    // 1. Composite ID 파싱
    // 2. 도메인 확인
    // 3. 각 속성에 대해:
    //    - 속성 존재 여부 확인
    //    - 존재하지 않으면 자동 생성 (autoCreate=true인 경우)
    //    - 노드에 속성 할당
    // 4. 결과 반환
}
```

#### 2.2 속성 타입 추론 로직
```go
func inferAttributeType(value string) models.AttributeType {
    // 숫자 패턴 확인
    if matched, _ := regexp.MatchString(`^\d+(\.\d+)?$`, value); matched {
        return models.AttributeTypeNumber
    }
    
    // URL 패턴 확인
    if strings.HasPrefix(value, "http") {
        return models.AttributeTypeString
    }
    
    // 기본값: tag
    return models.AttributeTypeTag
}
```

### Phase 3: 도메인 속성 관리 개선

#### 3.1 속성 생성 시 기본값 설정
```go
func (s *AutoAttributeService) createAttributeWithDefaults(
    ctx context.Context,
    domainID int,
    name string,
    value string,
) (*models.Attribute, error) {
    attrType := inferAttributeType(value)
    description := fmt.Sprintf("Auto-created attribute: %s", name)
    
    return &models.Attribute{
        DomainID:    domainID,
        Name:        name,
        Type:        attrType,
        Description: description,
        CreatedAt:   time.Now(),
    }, nil
}
```

### Phase 4: 사용자 경험 개선

#### 4.1 향상된 에러 메시지
```go
func (s *AutoAttributeService) getHelpfulErrorMessage(
    missingAttributes []string,
    domainName string,
) string {
    return fmt.Sprintf(
        "The following attributes do not exist in domain '%s': %s. "+
        "You can either:\n"+
        "1. Create them manually using create_domain_attribute\n"+
        "2. Use set_node_attributes_with_auto_create to create them automatically\n"+
        "3. Use set_node_attributes with auto_create_attributes=true",
        domainName,
        strings.Join(missingAttributes, ", "),
    )
}
```

#### 4.2 자동 생성 로그
```go
func (s *AutoAttributeService) logAutoCreation(
    domainName string,
    createdAttributes []string,
) {
    s.logger.Printf(
        "Auto-created %d attributes in domain '%s': %s",
        len(createdAttributes),
        domainName,
        strings.Join(createdAttributes, ", "),
    )
}
```

## 구현 단계

### Step 1: 기본 구조 구현 (1-2일)
- [ ] `AutoAttributeService` 구조체 생성
- [ ] 기본 속성 타입 추론 로직 구현
- [ ] 단위 테스트 작성

### Step 2: MCP 인터페이스 확장 (1일)
- [ ] `set_node_attributes_with_auto_create` 도구 추가
- [ ] 기존 `set_node_attributes`에 `auto_create_attributes` 파라미터 추가
- [ ] MCP 핸들러 업데이트

### Step 3: 백엔드 로직 구현 (2-3일)
- [ ] 속성 자동 생성 로직 구현
- [ ] 도메인 확인 및 에러 처리
- [ ] 트랜잭션 관리

### Step 4: 사용자 경험 개선 (1일)
- [ ] 향상된 에러 메시지 구현
- [ ] 자동 생성 로그 추가
- [ ] 문서 업데이트

### Step 5: 테스트 및 검증 (1일)
- [ ] 통합 테스트 작성
- [ ] 다양한 시나리오 테스트
- [ ] 성능 테스트

## 테스트 시나리오

### 1. 자동 생성 성공 케이스
```bash
# 존재하지 않는 속성으로 노드 속성 설정
set_node_attributes_with_auto_create({
  "composite_id": "url-db:test:1",
  "attributes": [
    {"name": "new_status", "value": "Active"},
    {"name": "priority", "value": "5"}
  ]
})
```

### 2. 기존 동작 유지
```bash
# auto_create_attributes=false로 기존 동작 테스트
set_node_attributes({
  "composite_id": "url-db:test:1",
  "attributes": [{"name": "nonexistent", "value": "test"}],
  "auto_create_attributes": false
})
# → VALIDATION_ERROR 발생
```

### 3. 혼합 케이스
```bash
# 일부는 존재, 일부는 없는 속성
set_node_attributes_with_auto_create({
  "composite_id": "url-db:test:1",
  "attributes": [
    {"name": "existing_attr", "value": "value1"},
    {"name": "new_attr", "value": "value2"}
  ]
})
# → existing_attr은 기존 사용, new_attr은 자동 생성
```

## 성능 고려사항

### 1. 데이터베이스 쿼리 최적화
- 속성 존재 여부를 배치로 확인
- 트랜잭션 내에서 속성 생성 및 할당

### 2. 메모리 사용량
- 대량의 속성 처리 시 메모리 효율성 고려
- 스트리밍 처리 가능성 검토

### 3. 동시성 처리
- 동일한 속성명으로 동시 생성 시 경쟁 상태 방지
- 데이터베이스 락 또는 UUID 사용

## 향후 확장 계획

### 1. 스마트 속성 타입 추론
- 머신러닝 기반 속성 타입 예측
- 사용자 패턴 학습

### 2. 속성 템플릿
- 자주 사용되는 속성 조합을 템플릿으로 제공
- 도메인별 기본 속성 세트

### 3. 속성 검증 규칙
- 도메인별 속성 검증 규칙 설정
- 자동 생성 시 검증 규칙 적용

## 성공 기준

### 1. 기능적 요구사항
- [ ] 존재하지 않는 속성을 자동으로 생성
- [ ] 적절한 속성 타입 추론
- [ ] 기존 동작과의 호환성 유지

### 2. 성능 요구사항
- [ ] 자동 생성 시 응답 시간 < 500ms
- [ ] 동시 요청 처리 가능
- [ ] 메모리 사용량 증가 < 10%

### 3. 사용자 경험 요구사항
- [ ] 명확한 에러 메시지 제공
- [ ] 자동 생성 로그 기록
- [ ] 문서화 완료

## 결론

이 구현 계획을 통해 사용자는 더욱 직관적이고 효율적으로 노드 속성을 관리할 수 있게 됩니다. 자동 생성 기능은 초기 설정의 복잡성을 줄이면서도, 필요에 따라 수동 제어도 가능하도록 설계됩니다. 