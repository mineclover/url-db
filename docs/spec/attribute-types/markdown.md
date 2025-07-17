# markdown 속성 타입

## 개요
마크다운 형식의 텍스트. 다중 라인 허용.

## 검증 규칙
- 최대 길이: 10000자
- 마크다운 구문 허용
- 개행 문자 허용

## 제약사항
- 하나의 URL에 하나의 값만 허용
- `order_index` 사용하지 않음

## 예시
```json
{
  "value": "# Notes\n\nThis is a **great** resource.",
  "order_index": null
}
```

## 에러 코드
- `validation_error`: 길이 초과
- `constraint_violation`: 이미 값이 존재하는 경우