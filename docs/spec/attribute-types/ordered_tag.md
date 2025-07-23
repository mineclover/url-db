# ordered_tag 속성 타입

## 개요
순서가 있는 태그. `order_index` 필수.

## 검증 규칙
- 최대 길이: 50자
- 금지 문자: `,`, `;`, `|`, `\n`, `\t`
- 대소문자 구분하지 않음 (소문자로 변환)
- `order_index`: 0 이상의 정수

## 제약사항
- `order_index` 필수
- 동일 URL에 같은 속성으로 중복 순서 불가

## 예시
```json
{
  "value": "high",
  "order_index": 1
}
```

## 에러 코드
- `validation_error`: 빈 값, 길이 초과, 금지 문자 포함, `order_index` 누락
- `constraint_violation`: 순서 중복