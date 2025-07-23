# 도메인 관리 에러 코드

## 에러 코드 매핑

### 도메인 생성/수정
- `validation_error`: 필수 필드 누락, 형식 오류
- `conflict`: 도메인 이름 중복

### 도메인 조회
- `not_found`: 도메인 존재하지 않음

### 도메인 삭제
- `not_found`: 도메인 존재하지 않음
- `state_conflict`: URL이 있는 도메인 삭제 시도

## 검증 규칙
- `name`: 필수, 1-255자, 영숫자와 하이픈만 허용
- `description`: 선택, 최대 1000자