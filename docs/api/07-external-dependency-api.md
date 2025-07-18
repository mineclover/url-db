# 외부 종속성 관리 API

## 개요
외부 종속성 관리 시스템을 위한 구독, 종속성, 이벤트 관리 API입니다.

## 구독 관리 API

### 구독 등록
노드의 이벤트를 구독합니다.

```http
POST /api/nodes/{nodeId}/subscriptions
```

**Request Body**
```json
{
  "subscriber_service": "analytics-service",
  "subscriber_endpoint": "https://analytics.example.com/webhook",
  "event_types": ["updated", "deleted"],
  "filter_conditions": {
    "attribute_filters": [
      {
        "attribute_name": "status",
        "operator": "equals",
        "value": "published"
      }
    ]
  }
}
```

**Response**
```json
{
  "id": 1,
  "subscriber_service": "analytics-service",
  "subscriber_endpoint": "https://analytics.example.com/webhook",
  "subscribed_node_id": 123,
  "event_types": ["updated", "deleted"],
  "filter_conditions": {
    "attribute_filters": [
      {
        "attribute_name": "status",
        "operator": "equals",
        "value": "published"
      }
    ]
  },
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### 구독 조회
```http
GET /api/subscriptions/{id}
```

### 구독 수정
```http
PUT /api/subscriptions/{id}
```

**Request Body**
```json
{
  "subscriber_endpoint": "https://analytics.example.com/webhook/v2",
  "event_types": ["updated", "deleted", "attribute_changed"],
  "is_active": false
}
```

### 구독 삭제
```http
DELETE /api/subscriptions/{id}
```

### 노드 구독 목록 조회
```http
GET /api/nodes/{nodeId}/subscriptions
```

### 서비스 구독 목록 조회
```http
GET /api/subscriptions?service=analytics-service
```

### 전체 구독 목록 조회
```http
GET /api/subscriptions?page=1&page_size=20
```

## 종속성 관리 API

### 종속성 등록
노드 간의 종속성을 등록합니다.

```http
POST /api/nodes/{nodeId}/dependencies
```

**Request Body**
```json
{
  "dependency_node_id": 456,
  "dependency_type": "soft",
  "cascade_delete": false,
  "cascade_update": true,
  "metadata": {
    "relationship": "parent-child",
    "description": "Article depends on category"
  }
}
```

**Response**
```json
{
  "id": 1,
  "dependent_node_id": 123,
  "dependency_node_id": 456,
  "dependency_type": "soft",
  "cascade_delete": false,
  "cascade_update": true,
  "metadata": {
    "relationship": "parent-child",
    "description": "Article depends on category"
  },
  "created_at": "2024-01-01T00:00:00Z"
}
```

### 종속성 조회
```http
GET /api/dependencies/{id}
```

### 종속성 삭제
```http
DELETE /api/dependencies/{id}
```

### 노드 종속성 목록 조회
노드가 의존하는 다른 노드들을 조회합니다.

```http
GET /api/nodes/{nodeId}/dependencies
```

### 노드 종속자 목록 조회
노드에 의존하는 다른 노드들을 조회합니다.

```http
GET /api/nodes/{nodeId}/dependents
```

## 이벤트 관리 API

### 노드 이벤트 조회
```http
GET /api/nodes/{nodeId}/events?limit=50
```

**Response**
```json
[
  {
    "id": 1,
    "node_id": 123,
    "event_type": "updated",
    "event_data": {
      "event_id": "uuid-123",
      "node_id": 123,
      "event_type": "updated",
      "timestamp": "2024-01-01T00:00:00Z",
      "changes": {
        "before": { "title": "Old Title" },
        "after": { "title": "New Title" }
      },
      "metadata": {
        "user": "system",
        "reason": "bulk_update"
      }
    },
    "occurred_at": "2024-01-01T00:00:00Z",
    "processed_at": "2024-01-01T00:00:01Z"
  }
]
```

### 미처리 이벤트 조회
```http
GET /api/events/pending?limit=100
```

### 이벤트 처리 완료 표시
```http
POST /api/events/{eventId}/process
```

### 이벤트 타입별 조회
```http
GET /api/events?type=updated&start=2024-01-01T00:00:00Z&end=2024-01-02T00:00:00Z
```

### 이벤트 통계 조회
```http
GET /api/events/stats
```

**Response**
```json
{
  "total_events": 1000,
  "pending_events": 5,
  "events_by_type": {
    "created": 300,
    "updated": 500,
    "deleted": 100,
    "attribute_changed": 100
  }
}
```

### 이벤트 정리
지정된 일수보다 오래된 처리 완료된 이벤트를 삭제합니다.

```http
POST /api/events/cleanup?days=30
```

**Response**
```json
{
  "deleted_events": 250,
  "message": "Events cleaned up successfully"
}
```

## 종속성 타입

### Hard Dependency
- 상위 노드 삭제 시 하위 노드도 자동 삭제
- 중요한 변경사항 즉시 반영

### Soft Dependency
- 상위 노드 삭제 시 하위 노드는 유지
- 변경사항 알림만 전송

### Reference Dependency
- 단순 참조 관계만 유지
- 변경사항의 직접적인 영향 없음

## 이벤트 타입

- `created`: 노드 생성
- `updated`: 노드 수정
- `deleted`: 노드 삭제
- `attribute_changed`: 속성 변경
- `connection_changed`: 연결 관계 변경
- `dependency_updated`: 종속성 업데이트

## 필터 조건

### 속성 필터
```json
{
  "attribute_filters": [
    {
      "attribute_name": "status",
      "operator": "equals",
      "value": "published"
    }
  ]
}
```

### 변경 필터
```json
{
  "change_filters": {
    "fields": ["title", "content"],
    "ignore_fields": ["updated_at"]
  }
}
```

## 에러 응답

### 400 Bad Request
- `validation_error`: 요청 데이터 검증 실패
- `invalid_request`: 잘못된 요청 구조

### 404 Not Found
- `not_found`: 리소스 존재하지 않음

### 409 Conflict
- `conflict`: 순환 종속성 감지

### 500 Internal Server Error
- `internal_error`: 서버 내부 오류