# URL API

## 개요

URL API는 URL 노드의 생성, 조회, 수정, 삭제를 위한 엔드포인트를 제공합니다.

## 엔드포인트

### 노드 속성 설정

#### PUT /api/mcp/nodes/{composite_id}/attributes

노드의 속성을 설정합니다.

**요청 본문:**
```json
{
  "attributes": [
    {
      "name": "status",
      "value": "active",
      "order_index": 1
    }
  ],
  "auto_create_attributes": false
}
```

**파라미터:**
- `composite_id` (path): 노드의 복합 ID (예: `url-db:domain:123`)
- `attributes` (body): 설정할 속성 배열
  - `name`: 속성 이름
  - `value`: 속성 값
  - `order_index`: 순서 인덱스 (선택사항)
- `auto_create_attributes` (body): 존재하지 않는 속성을 자동 생성할지 여부 (기본값: false)

**응답:**
```json
{
  "composite_id": "url-db:domain:123",
  "attributes": [
    {
      "name": "status",
      "type": "tag",
      "value": "active"
    }
  ]
}
```

**자동 속성 생성 기능:**

`auto_create_attributes`가 `true`로 설정되면, 존재하지 않는 속성이 자동으로 생성됩니다:

- **숫자 값**: `number` 타입으로 추론
- **URL 값**: `string` 타입으로 추론  
- **기타 값**: `tag` 타입으로 추론

**예시:**
```json
{
  "attributes": [
    {"name": "priority", "value": "5"},
    {"name": "url", "value": "https://example.com"},
    {"name": "category", "value": "tutorial"}
  ],
  "auto_create_attributes": true
}
```

위 요청은 다음 속성들을 자동 생성합니다:
- `priority`: number 타입 (숫자 값)
- `url`: string 타입 (URL 값)
- `category`: tag 타입 (기본값)

**에러 처리:**

자동 생성이 비활성화된 경우 존재하지 않는 속성에 대해 명확한 에러 메시지를 제공합니다:

```json
{
  "error": "validation_error",
  "message": "The following attributes do not exist in domain 'test': status, priority. You can either:\n1. Create them manually using create_domain_attribute\n2. Use set_node_attributes_with_auto_create to create them automatically\n3. Use set_node_attributes with auto_create_attributes=true"
}
```

**환경변수 설정:**

서버 시작 시 `AUTO_CREATE_ATTRIBUTES` 환경변수로 기본값을 설정할 수 있습니다:

```bash
# 자동 속성 생성 활성화 (기본값)
export AUTO_CREATE_ATTRIBUTES=true
./bin/url-db

# 자동 속성 생성 비활성화
export AUTO_CREATE_ATTRIBUTES=false
./bin/url-db

# 또는 직접 설정
AUTO_CREATE_ATTRIBUTES=true ./bin/url-db
```

**지원되는 값:**
- `true`, `1`, `yes`, `on`: 자동 생성 활성화
- `false`, `0`, `no`, `off`: 자동 생성 비활성화
- 설정하지 않으면 기본값 `true` 사용