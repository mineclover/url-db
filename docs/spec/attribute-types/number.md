# number 속성 타입

## 개요
숫자 값. 정수 또는 실수 허용.

## 검증 규칙
- 유효한 숫자 형식
- 최소값/최대값 제한 없음 (기본값)
- 정수 제한 없음 (기본값)

## 제약사항
- 하나의 URL에 하나의 값만 허용
- `order_index` 사용하지 않음

## 예시
```json
{
  "value": "4.5",
  "order_index": null
}
```

## 에러 코드
- `validation_error`: 숫자 형식 오류
- `constraint_violation`: 이미 값이 존재하는 경우