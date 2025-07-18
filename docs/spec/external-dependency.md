# 외부 종속성 관리 시스템

## 개요
URL-DB 시스템에서 노드 간의 외부 종속성을 관리하기 위한 구독 시스템입니다.
하나의 서비스가 자신의 상태 변경을 외부에 알리고, 다른 서비스들이 이를 구독하여 
종속성을 관리할 수 있도록 설계되었습니다.

## 주요 개념

### 1. 이벤트 게시자 (Publisher)
- 노드의 상태 변경 시 이벤트를 발생시키는 주체
- 생성, 수정, 삭제 등의 이벤트 타입 정의

### 2. 구독자 (Subscriber)
- 특정 노드의 이벤트를 구독하는 외부 서비스
- 구독 조건과 콜백 정보를 등록

### 3. 종속성 추적
- 노드 간의 종속 관계를 명시적으로 관리
- 계층적 종속성 및 순환 종속성 감지

## 데이터베이스 스키마

### node_subscriptions (노드 구독 테이블)
```sql
CREATE TABLE node_subscriptions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    subscriber_service TEXT NOT NULL,     -- 구독자 서비스 식별자
    subscriber_endpoint TEXT,             -- 콜백 엔드포인트 (옵션)
    subscribed_node_id INTEGER NOT NULL,  -- 구독 대상 노드
    event_types TEXT NOT NULL,            -- JSON 배열: ["created", "updated", "deleted"]
    filter_conditions TEXT,               -- JSON: 구독 필터 조건
    is_active BOOLEAN DEFAULT TRUE,       -- 구독 활성화 상태
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (subscribed_node_id) REFERENCES nodes(id) ON DELETE CASCADE
);

-- 인덱스
CREATE INDEX idx_subscriptions_service ON node_subscriptions(subscriber_service);
CREATE INDEX idx_subscriptions_node ON node_subscriptions(subscribed_node_id);
CREATE INDEX idx_subscriptions_active ON node_subscriptions(is_active);
```

### node_dependencies (노드 종속성 테이블)
```sql
CREATE TABLE node_dependencies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    dependent_node_id INTEGER NOT NULL,   -- 종속 노드 (하위)
    dependency_node_id INTEGER NOT NULL,  -- 종속 대상 노드 (상위)
    dependency_type TEXT NOT NULL,        -- 'hard', 'soft', 'reference'
    cascade_delete BOOLEAN DEFAULT FALSE, -- 상위 삭제 시 하위도 삭제
    cascade_update BOOLEAN DEFAULT FALSE, -- 상위 변경 시 하위 알림
    metadata TEXT,                        -- JSON: 추가 메타데이터
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (dependent_node_id) REFERENCES nodes(id) ON DELETE CASCADE,
    FOREIGN KEY (dependency_node_id) REFERENCES nodes(id) ON DELETE CASCADE,
    UNIQUE(dependent_node_id, dependency_node_id)
);

-- 인덱스
CREATE INDEX idx_dependencies_dependent ON node_dependencies(dependent_node_id);
CREATE INDEX idx_dependencies_dependency ON node_dependencies(dependency_node_id);
```

### node_events (노드 이벤트 로그)
```sql
CREATE TABLE node_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    node_id INTEGER NOT NULL,
    event_type TEXT NOT NULL,             -- 'created', 'updated', 'deleted', 'attribute_changed'
    event_data TEXT,                      -- JSON: 이벤트 상세 데이터
    occurred_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    processed_at DATETIME,                -- 처리 완료 시간
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE
);

-- 인덱스
CREATE INDEX idx_events_node ON node_events(node_id);
CREATE INDEX idx_events_type ON node_events(event_type);
CREATE INDEX idx_events_occurred ON node_events(occurred_at);
CREATE INDEX idx_events_unprocessed ON node_events(processed_at) WHERE processed_at IS NULL;
```

## 종속성 타입

### 1. Hard Dependency (강한 종속성)
- 상위 노드가 삭제되면 하위 노드도 반드시 삭제
- 상위 노드의 중요한 변경사항은 하위 노드에 즉시 반영

### 2. Soft Dependency (약한 종속성)
- 상위 노드가 삭제되어도 하위 노드는 유지
- 상위 노드 변경 시 하위 노드에 알림만 전송

### 3. Reference Dependency (참조 종속성)
- 단순 참조 관계만 유지
- 상위 노드 변경이 하위 노드에 직접적인 영향 없음

## 이벤트 타입

### 기본 이벤트
- `created`: 노드 생성
- `updated`: 노드 수정
- `deleted`: 노드 삭제
- `attribute_changed`: 속성 변경
- `connection_changed`: 연결 관계 변경

### 이벤트 데이터 형식
```json
{
  "event_id": "uuid",
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
}
```

## 구독 필터 조건

구독자는 특정 조건에 맞는 이벤트만 수신할 수 있습니다:

```json
{
  "attribute_filters": [
    {
      "attribute_name": "status",
      "operator": "equals",
      "value": "published"
    }
  ],
  "change_filters": {
    "fields": ["title", "content"],
    "ignore_fields": ["updated_at"]
  }
}
```

## API 엔드포인트

### 구독 관리
- `POST /api/nodes/{nodeId}/subscriptions` - 구독 등록
- `GET /api/subscriptions` - 구독 목록 조회
- `PUT /api/subscriptions/{id}` - 구독 수정
- `DELETE /api/subscriptions/{id}` - 구독 해제

### 종속성 관리
- `POST /api/nodes/{nodeId}/dependencies` - 종속성 등록
- `GET /api/nodes/{nodeId}/dependencies` - 종속성 조회
- `DELETE /api/dependencies/{id}` - 종속성 제거

### 이벤트 조회
- `GET /api/nodes/{nodeId}/events` - 노드 이벤트 조회
- `GET /api/events/pending` - 미처리 이벤트 조회

## 사용 예시

### 구독 등록
```bash
POST /api/nodes/123/subscriptions
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

### 종속성 등록
```bash
POST /api/nodes/456/dependencies
{
  "dependency_node_id": 123,
  "dependency_type": "soft",
  "cascade_update": true,
  "metadata": {
    "relationship": "parent-child",
    "description": "Article depends on category"
  }
}
```

## 보안 고려사항

1. **인증 및 권한**
   - 구독 등록 시 서비스 인증 필요
   - 종속성 생성 권한 검증

2. **이벤트 전달**
   - HTTPS를 통한 안전한 웹훅 전달
   - 재시도 및 실패 처리 메커니즘

3. **순환 종속성 방지**
   - 종속성 등록 시 순환 참조 검사
   - 깊이 제한 설정