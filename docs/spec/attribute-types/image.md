# image 속성 타입

## 개요
이미지 데이터. Base64 인코딩 또는 URL 형식 지원.

## 검증 규칙
- Base64 형식: `data:image/{type};base64,{data}`
- URL 형식: `http://` 또는 `https://`로 시작
- 지원 MIME 타입: `image/jpeg`, `image/png`, `image/gif`, `image/webp`
- 최대 크기: 10MB (Base64의 경우)

## 제약사항
- 하나의 URL에 하나의 값만 허용
- `order_index` 사용하지 않음

## 예시
```json
{
  "value": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEAYABgAAD...",
  "order_index": null
}
```

## 에러 코드
- `validation_error`: 잘못된 형식, 지원하지 않는 MIME 타입, 크기 초과
- `constraint_violation`: 이미 값이 존재하는 경우