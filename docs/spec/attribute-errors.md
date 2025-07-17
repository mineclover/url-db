# 속성 관리 에러 코드

## 에러 코드 매핑

### 속성 생성/수정
- `validation_error`: 필수 필드 누락, 잘못된 타입
- `conflict`: 도메인 내 속성 이름 중복

### 속성 조회
- `not_found`: 속성 존재하지 않음

### 속성 삭제
- `not_found`: 속성 존재하지 않음
- `state_conflict`: 속성 값이 있는 속성 삭제 시도

## 검증 규칙
- `name`: 필수, 1-100자, 영숫자와 언더스코어만 허용
- `type`: 필수, tag|ordered_tag|number|string|markdown|image
- `description`: 선택, 최대 500자
- 속성 타입은 생성 후 변경 불가