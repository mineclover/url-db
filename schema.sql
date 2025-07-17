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