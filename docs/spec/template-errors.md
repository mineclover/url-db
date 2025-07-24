# 템플릿 관리 에러 코드

## 에러 코드 매핑

### 템플릿 생성/수정
- `validation_error`: 필수 필드 누락 또는 유효하지 않은 값
- `conflict`: 도메인 내 템플릿 이름 중복
- `invalid_json`: template_data가 유효한 JSON 형식이 아님

### 템플릿 조회
- `not_found`: 템플릿 존재하지 않음
- `inactive`: 비활성화된 템플릿 접근 시도

### 템플릿 삭제
- `not_found`: 템플릿 존재하지 않음
- `in_use`: 다른 엔티티에서 참조 중인 템플릿

### 템플릿 활성화/비활성화
- `not_found`: 템플릿 존재하지 않음
- `already_active`: 이미 활성화된 템플릿
- `already_inactive`: 이미 비활성화된 템플릿

## 검증 규칙
- `name`: 필수, 최대 255자, 도메인 내 고유
- `template_data`: 필수, 유효한 JSON 형식
- `title`: 선택, 최대 255자
- `description`: 선택, 최대 1000자
- `is_active`: 기본값 true
- 템플릿 이름과 domain_id는 생성 후 수정 불가