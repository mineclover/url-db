# 속성 관리 API 엔드포인트

## 개요
도메인 내 속성 정의를 생성, 조회, 수정, 삭제하는 REST API

## 엔드포인트 목록

### 1. 속성 생성
- **POST** `/api/domains/{domain_id}/attributes`
- **요청 본문**:
```json
{
  "name": "category",
  "type": "tag",
  "description": "Category tag for URLs"
}
```
- **응답 (201)**:
```json
{
  "id": 1,
  "domain_id": 1,
  "name": "category",
  "type": "tag",
  "description": "Category tag for URLs",
  "created_at": "2023-01-01T00:00:00Z"
}
```

### 2. 도메인별 속성 목록 조회
- **GET** `/api/domains/{domain_id}/attributes`
- **응답 (200)**:
```json
{
  "attributes": [
    {
      "id": 1,
      "domain_id": 1,
      "name": "category",
      "type": "tag",
      "description": "Category tag for URLs",
      "created_at": "2023-01-01T00:00:00Z"
    },
    {
      "id": 2,
      "domain_id": 1,
      "name": "priority",
      "type": "ordered_tag",
      "description": "Priority level",
      "created_at": "2023-01-01T00:00:00Z"
    }
  ]
}
```

### 3. 속성 조회
- **GET** `/api/attributes/{id}`
- **응답 (200)**:
```json
{
  "id": 1,
  "domain_id": 1,
  "name": "category",
  "type": "tag",
  "description": "Category tag for URLs",
  "created_at": "2023-01-01T00:00:00Z"
}
```

### 4. 속성 수정
- **PUT** `/api/attributes/{id}`
- **요청 본문**:
```json
{
  "name": "updated-category",
  "description": "Updated category description"
}
```
- **응답 (200)**:
```json
{
  "id": 1,
  "domain_id": 1,
  "name": "updated-category",
  "type": "tag",
  "description": "Updated category description",
  "created_at": "2023-01-01T00:00:00Z"
}
```

### 5. 속성 삭제
- **DELETE** `/api/attributes/{id}`
- **응답 (204)**: 본문 없음

## 속성 타입

### 지원하는 속성 타입
- `tag`: 일반 태그 ([상세](../spec/attribute-types/tag.md))
- `ordered_tag`: 순서가 있는 태그 ([상세](../spec/attribute-types/ordered_tag.md))
- `number`: 숫자 값 ([상세](../spec/attribute-types/number.md))
- `string`: 문자열 값 ([상세](../spec/attribute-types/string.md))
- `markdown`: 마크다운 텍스트 ([상세](../spec/attribute-types/markdown.md))
- `image`: 이미지 데이터 ([상세](../spec/attribute-types/image.md))

## 에러 응답

> 에러 응답 형식: [`../spec/error-codes.md`](../spec/error-codes.md)  
> 속성 관련 에러: [`../spec/attribute-errors.md`](../spec/attribute-errors.md)

## 검증 규칙
- `name`: 필수, 1-100자, 영숫자와 언더스코어만 허용
- `type`: 필수, 지원하는 타입 중 하나
- `description`: 선택, 최대 500자
- 속성 타입은 생성 후 변경 불가