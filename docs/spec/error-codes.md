# 에러 코드 정의

## 표준 에러 응답 형식
```json
{
  "error": "error_code",
  "message": "Human readable message",
  "details": [
    {
      "field": "field_name",
      "message": "field error message"
    }
  ]
}
```

## HTTP 상태 코드별 에러 코드

### 400 Bad Request
- `validation_error`: 요청 데이터 검증 실패
- `invalid_format`: 데이터 형식 오류
- `invalid_request`: 잘못된 요청 구조

### 401 Unauthorized
- `authentication_required`: 인증 필요
- `invalid_credentials`: 잘못된 인증 정보

### 403 Forbidden
- `permission_denied`: 권한 부족

### 404 Not Found
- `not_found`: 리소스 존재하지 않음
- `endpoint_not_found`: 엔드포인트 존재하지 않음

### 409 Conflict
- `conflict`: 리소스 충돌 (중복 등)
- `state_conflict`: 리소스 상태 충돌

### 422 Unprocessable Entity
- `business_rule_violation`: 비즈니스 규칙 위반
- `constraint_violation`: 제약 조건 위반

### 429 Too Many Requests
- `rate_limit_exceeded`: 요청 한도 초과

### 500 Internal Server Error
- `internal_error`: 서버 내부 오류
- `database_error`: 데이터베이스 오류

### 503 Service Unavailable
- `service_unavailable`: 서비스 일시 중단