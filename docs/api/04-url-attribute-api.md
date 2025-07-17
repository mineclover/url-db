# 노드 속성 값 관리 API 엔드포인트

## 개요
노드에 속성 값을 설정하고 관리하는 REST API

## 엔드포인트 목록

### 1. 노드 속성 값 추가
- **POST** `/api/urls/{url_id}/attributes`
- **요청 본문**:
```json
{
  "attribute_id": 1,
  "value": "programming",
  "order_index": 1
}
```
- **응답 (201)**:
```json
{
  "id": 1,
  "node_id": 1,
  "attribute_id": 1,
  "value": "programming",
  "order_index": 1,
  "created_at": "2023-01-01T00:00:00Z"
}
```

### 2. 노드별 속성 값 목록 조회
- **GET** `/api/urls/{url_id}/attributes`
- **응답 (200)**:
```json
{
  "attributes": [
    {
      "id": 1,
      "node_id": 1,
      "attribute_id": 1,
      "attribute_name": "category",
      "attribute_type": "tag",
      "value": "programming",
      "order_index": null,
      "created_at": "2023-01-01T00:00:00Z"
    },
    {
      "id": 2,
      "node_id": 1,
      "attribute_id": 2,
      "attribute_name": "priority",
      "attribute_type": "ordered_tag",
      "value": "high",
      "order_index": 1,
      "created_at": "2023-01-01T00:00:00Z"
    }
  ]
}
```

### 3. 속성 값 조회
- **GET** `/api/url-attributes/{id}`
- **응답 (200)**:
```json
{
  "id": 1,
  "node_id": 1,
  "attribute_id": 1,
  "attribute_name": "category",
  "attribute_type": "tag",
  "value": "programming",
  "order_index": null,
  "created_at": "2023-01-01T00:00:00Z"
}
```

### 4. 속성 값 수정
- **PUT** `/api/url-attributes/{id}`
- **요청 본문**:
```json
{
  "value": "updated-programming",
  "order_index": 2
}
```
- **응답 (200)**:
```json
{
  "id": 1,
  "node_id": 1,
  "attribute_id": 1,
  "value": "updated-programming",
  "order_index": 2,
  "created_at": "2023-01-01T00:00:00Z"
}
```

### 5. 속성 값 삭제
- **DELETE** `/api/url-attributes/{id}`
- **응답 (204)**: 본문 없음

### 6. 노드의 특정 속성 값 삭제
- **DELETE** `/api/urls/{url_id}/attributes/{attribute_id}`
- **응답 (204)**: 본문 없음

## 속성 타입별 제약사항

### 속성 타입별 상세 스펙
- `tag`: [상세](../spec/attribute-types/tag.md)
- `ordered_tag`: [상세](../spec/attribute-types/ordered_tag.md)
- `number`: [상세](../spec/attribute-types/number.md)
- `string`: [상세](../spec/attribute-types/string.md)
- `markdown`: [상세](../spec/attribute-types/markdown.md)
- `image`: [상세](../spec/attribute-types/image.md)

## 에러 응답

> 에러 응답 형식: [`../spec/error-codes.md`](../spec/error-codes.md)  
> 노드 속성 관련 에러: [`../spec/node-attribute-errors.md`](../spec/node-attribute-errors.md)

## 검증 규칙
- `attribute_id`: 필수, 존재하는 속성 ID
- `value`: 필수, 속성 타입에 따른 형식 검증
- `order_index`: ordered_tag 타입에서만 필수
- 노드와 속성은 같은 도메인에 속해야 함