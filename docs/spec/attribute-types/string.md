# string 속성 타입

## 개요
일반 문자열 값. 단일 라인 텍스트.

## 검증 규칙
- 최대 길이: 1000자
- 최소 길이: 0자
- 개행 문자 허용하지 않음

## 제약사항
- 하나의 URL에 하나의 값만 허용
- `order_index` 사용하지 않음

## 예시
```json
{
  "value": "My favorite article",
  "order_index": null
}
```

## 에러 코드
- `validation_error`: 길이 초과
- `constraint_violation`: 이미 값이 존재하는 경우