# 템플릿 속성 값 정의 시스템

## 템플릿의 주 목적

템플릿은 **속성이 허용하는 값의 범위와 제약사항을 정의**하는 시스템입니다:

- **값 템플릿 정의**: 각 속성에 대해 허용되는 값의 형식과 범위를 미리 정의
- **데이터 일관성 보장**: 도메인 내 모든 노드가 동일한 속성 값 규칙을 따르도록 강제
- **입력 가이드라인**: 사용자가 노드 속성을 설정할 때 어떤 값이 허용되는지 명확히 제시
- **재사용 가능한 스키마**: 한 번 정의된 템플릿을 여러 노드에서 참조하여 사용

### 예시 시나리오

```json
// 1. 템플릿 정의: "제품" 도메인의 속성 값 템플릿
{
  "name": "product-template",
  "domain": "products", 
  "attributes": {
    "category": ["전자제품", "의류", "도서", "생활용품"],    // 허용되는 카테고리 값들
    "price_range": ["0-10000", "10000-50000", "50000+"],   // 허용되는 가격 범위
    "status": ["판매중", "품절", "단종"]                    // 허용되는 상태 값들
  }
}

// 2. 노드 생성 시 템플릿 기반 검증
// ✅ 성공: 템플릿에 정의된 값 사용
create_node_with_attributes(
  url: "https://shop.com/laptop",
  attributes: {
    "category": "전자제품",    // 템플릿에 정의된 값 ✓
    "price_range": "50000+",  // 템플릿에 정의된 값 ✓
    "status": "판매중"        // 템플릿에 정의된 값 ✓
  }
)

// ❌ 실패: 템플릿에 정의되지 않은 값 사용
create_node_with_attributes(
  url: "https://shop.com/phone", 
  attributes: {
    "category": "스마트폰",    // 템플릿에 없는 값 ❌
    "status": "예약판매"      // 템플릿에 없는 값 ❌  
  }
)
```

## 에러 코드 매핑

### 템플릿 기반 값 검증 에러
- `template_value_not_allowed`: 템플릿에 정의되지 않은 값 사용 시도
- `template_required_but_missing`: 템플릿이 필수인 속성에 템플릿 없이 값 설정 시도  
- `template_value_format_mismatch`: 템플릿 정의 형식과 실제 값 형식 불일치

### 속성 값 생성/수정  
- `validation_error`: 필수 필드 누락, 값 형식 오류
- `conflict`: 태그 값 중복
- `business_rule_violation`: 템플릿과 속성 도메인 불일치
- `constraint_violation`: 순서 중복, 단일 값 속성 중복
- `template_inactive`: 비활성화된 템플릿에 속성 추가 시도

### 속성 값 조회
- `not_found`: 템플릿 속성 존재하지 않음
- `template_inactive`: 비활성화된 템플릿의 속성 조회 시도

### 속성 값 삭제
- `not_found`: 템플릿 속성 존재하지 않음
- `template_inactive`: 비활성화된 템플릿의 속성 삭제 시도

### 속성 값 일괄 설정
- `validation_error`: 속성 배열 형식 오류
- `business_rule_violation`: 템플릿과 속성 도메인 불일치
- `template_inactive`: 비활성화된 템플릿에 속성 설정 시도

## 템플릿 기반 검증 규칙

### 기본 검증
- `attribute_id`: 필수, 존재하는 속성 ID
- `value`: 필수, 속성 타입에 따른 형식 검증
- `order_index`: ordered_tag 타입에서만 필수
- 템플릿과 속성은 같은 도메인에 속해야 함
- 비활성화된 템플릿의 속성은 읽기 전용

### 템플릿 기반 값 제약
1. **허용 값 목록 검증**: 템플릿에 정의된 값만 허용
   ```json
   // 템플릿: ["small", "medium", "large"]
   // ✅ 허용: "medium"
   // ❌ 거부: "extra-large"
   ```

2. **값 형식 패턴 검증**: 정규식 패턴 기반 검증
   ```json
   // 템플릿: {"pattern": "^[A-Z]{2}-\\d{6}$"}
   // ✅ 허용: "KR-123456"
   // ❌ 거부: "kr-123456"
   ```

3. **값 범위 검증**: 숫자/날짜 범위 제한
   ```json
   // 템플릿: {"min": 0, "max": 100}
   // ✅ 허용: "75"
   // ❌ 거부: "150"
   ```

4. **조건부 값 검증**: 다른 속성 값에 따른 조건부 허용
   ```json
   // 템플릿: {"if": {"category": "전자제품"}, "then": ["1년", "2년", "3년"]}
   // category가 "전자제품"일 때만 보증기간 값 허용
   ```