-- URL 데이터베이스 스키마

-- 도메인 폴더 테이블
CREATE TABLE domains (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 노드 테이블 (URL 등의 컨텐츠 저장)
CREATE TABLE nodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
    domain_id INTEGER NOT NULL,
    title TEXT,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE,
    UNIQUE(content, domain_id)
);

-- 속성 정의 테이블
CREATE TABLE attributes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('tag', 'ordered_tag', 'number', 'string', 'markdown', 'image')),
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE,
    UNIQUE(domain_id, name)
);

-- 노드 속성 값 테이블
CREATE TABLE node_attributes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    node_id INTEGER NOT NULL,
    attribute_id INTEGER NOT NULL,
    value TEXT NOT NULL,
    order_index INTEGER, -- 순서가 중요한 태그용
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE,
    FOREIGN KEY (attribute_id) REFERENCES attributes(id) ON DELETE CASCADE
);

-- 노드 간 연결관계 테이블
CREATE TABLE node_connections (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source_node_id INTEGER NOT NULL,
    target_node_id INTEGER NOT NULL,
    relationship_type TEXT NOT NULL, -- 'parent', 'child', 'related', 'next', 'prev' 등
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (source_node_id) REFERENCES nodes(id) ON DELETE CASCADE,
    FOREIGN KEY (target_node_id) REFERENCES nodes(id) ON DELETE CASCADE,
    UNIQUE(source_node_id, target_node_id, relationship_type)
);

-- 인덱스 생성
CREATE INDEX idx_nodes_domain ON nodes(domain_id);
CREATE INDEX idx_nodes_content ON nodes(content);
CREATE INDEX idx_attributes_domain ON attributes(domain_id);
CREATE INDEX idx_node_attributes_node ON node_attributes(node_id);
CREATE INDEX idx_node_attributes_attribute ON node_attributes(attribute_id);
CREATE INDEX idx_node_connections_source ON node_connections(source_node_id);
CREATE INDEX idx_node_connections_target ON node_connections(target_node_id);

-- 트리거: updated_at 자동 업데이트
CREATE TRIGGER domains_updated_at 
    AFTER UPDATE ON domains 
    FOR EACH ROW 
    BEGIN 
        UPDATE domains SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER nodes_updated_at 
    AFTER UPDATE ON nodes 
    FOR EACH ROW 
    BEGIN 
        UPDATE nodes SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

-- 외부 종속성 관리 테이블

-- 노드 구독 테이블
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

-- 노드 종속성 테이블
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

-- 노드 이벤트 로그
CREATE TABLE node_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    node_id INTEGER NOT NULL,
    event_type TEXT NOT NULL,             -- 'created', 'updated', 'deleted', 'attribute_changed'
    event_data TEXT,                      -- JSON: 이벤트 상세 데이터
    occurred_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    processed_at DATETIME,                -- 처리 완료 시간
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE
);

-- 외부 종속성 관련 인덱스
CREATE INDEX idx_subscriptions_service ON node_subscriptions(subscriber_service);
CREATE INDEX idx_subscriptions_node ON node_subscriptions(subscribed_node_id);
CREATE INDEX idx_subscriptions_active ON node_subscriptions(is_active);
CREATE INDEX idx_dependencies_dependent ON node_dependencies(dependent_node_id);
CREATE INDEX idx_dependencies_dependency ON node_dependencies(dependency_node_id);
CREATE INDEX idx_events_node ON node_events(node_id);
CREATE INDEX idx_events_type ON node_events(event_type);
CREATE INDEX idx_events_occurred ON node_events(occurred_at);
CREATE INDEX idx_events_unprocessed ON node_events(processed_at) WHERE processed_at IS NULL;

-- 트리거: node_subscriptions updated_at 자동 업데이트
CREATE TRIGGER node_subscriptions_updated_at 
    AFTER UPDATE ON node_subscriptions 
    FOR EACH ROW 
    BEGIN 
        UPDATE node_subscriptions SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;