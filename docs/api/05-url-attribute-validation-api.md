# 노드 속성 확인 API 엔드포인트

## 개요
노드에 설정 가능한 속성을 확인하고 속성 값을 검증하는 REST API

## 엔드포인트 목록

### 1. 노드에 설정 가능한 속성 목록 조회
- **GET** `/api/urls/{url_id}/available-attributes`
- **응답 (200)**:
```json
{
  "attributes": [
    {
      "id": 1,
      "name": "category",
      "type": "tag",
      "description": "Category tag for URLs",
      "constraints": {
        "max_length": 50,
        "invalid_chars": [",", ";", "|", "\n", "\t"],
        "case_sensitive": false
      },
      "examples": ["programming", "tutorial", "reference"]
    },
    {
      "id": 2,
      "name": "priority",
      "type": "ordered_tag",
      "description": "Priority level",
      "constraints": {
        "max_length": 50,
        "invalid_chars": [",", ";", "|", "\n", "\t"],
        "case_sensitive": false
      },
      "examples": ["high", "medium", "low"]
    },
    {
      "id": 3,
      "name": "rating",
      "type": "number",
      "description": "Rating score",
      "constraints": {
        "min_value": 0,
        "max_value": 10,
        "integer": false
      },
      "examples": ["4.5", "8", "9.2"]
    }
  ]
}
```

### 2. 속성 값 검증
- **POST** `/api/attributes/validate`
- **요청 본문**:
```json
{
  "attribute_id": 1,
  "value": "programming",
  "order_index": 1
}
```
- **응답 (200)**:
```json
{
  "valid": true,
  "warnings": [
    "Tag will be converted to lowercase for consistency"
  ]
}
```

### 3. 속성 제약사항 조회
- **GET** `/api/attributes/{attribute_id}/constraints`
- **응답 (200)**:
```json
{
  "type": "tag",
  "constraints": {
    "max_length": 50,
    "invalid_chars": [",", ";", "|", "\n", "\t"],
    "case_sensitive": false
  }
}
```

### 4. 속성 예시 조회
- **GET** `/api/attributes/{attribute_id}/examples`
- **응답 (200)**:
```json
{
  "examples": [
    "programming",
    "tutorial",
    "reference",
    "documentation"
  ]
}
```

## 속성 타입별 제약사항

### 속성 타입별 상세 스펙
- `tag`: [상세](../spec/attribute-types/tag.md)
- `ordered_tag`: [상세](../spec/attribute-types/ordered_tag.md)
- `number`: [상세](../spec/attribute-types/number.md)
- `string`: [상세](../spec/attribute-types/string.md)
- `markdown`: [상세](../spec/attribute-types/markdown.md)
- `image`: [상세](../spec/attribute-types/image.md)

## 검증 응답

### 성공 응답
```json
{
  "valid": true,
  "warnings": [
    "Tag will be converted to lowercase for consistency"
  ]
}
```

### 실패 응답
```json
{
  "valid": false,
  "error_message": "Tag contains invalid characters",
  "suggestions": [
    "Remove special characters like commas, semicolons, or pipes",
    "Use spaces or hyphens instead of special characters"
  ]
}
```

## 에러 응답

> 에러 응답 형식: [`../spec/error-codes.md`](../spec/error-codes.md)  
> 속성 관련 에러: [`../spec/attribute-errors.md`](../spec/attribute-errors.md)

## 속성 타입별 예시

### tag 타입 예시
```json
{
  "examples": [
    "programming",
    "tutorial",
    "reference",
    "documentation"
  ]
}
```

### number 타입 예시
```json
{
  "examples": [
    "42",
    "3.14",
    "100"
  ]
}
```

### markdown 타입 예시
```json
{
  "examples": [
    "# Notes\n\nThis is a **great** resource.",
    "## Summary\n\n- Point 1\n- Point 2"
  ]
}
```

### image 타입 예시
```json
{
  "examples": [
    "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEAYABgAAD...",
    "https://example.com/image.jpg"
  ]
}
```

## 검증 규칙
- `attribute_id`: 필수, 존재하는 속성 ID
- `value`: 필수, 검증할 값
- `order_index`: ordered_tag 타입에서만 사용