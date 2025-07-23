# 노드 속성 값 관리 에러 코드

## 에러 코드 매핑

### 속성 값 생성/수정
- `validation_error`: 필수 필드 누락, 값 형식 오류
- `conflict`: 태그 값 중복
- `business_rule_violation`: 노드와 속성 도메인 불일치
- `constraint_violation`: 순서 중복, 단일 값 속성 중복

### 속성 값 조회
- `not_found`: 노드 속성 존재하지 않음

### 속성 값 삭제
- `not_found`: 노드 속성 존재하지 않음

## 검증 규칙
- `attribute_id`: 필수, 존재하는 속성 ID
- `value`: 필수, 속성 타입에 따른 형식 검증
- `order_index`: ordered_tag 타입에서만 필수
- 노드와 속성은 같은 도메인에 속해야 함