-- URL-DB Complete Database Schema
-- Single source of truth for all database structures
-- Auto-generated from internal/database/database.go

-- 도메인 폴더 테이블
CREATE TABLE IF NOT EXISTS domains (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE,
	description TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 노드 테이블 (URL 등의 컨텐츠 저장)
CREATE TABLE IF NOT EXISTS nodes (
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
CREATE TABLE IF NOT EXISTS attributes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	domain_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	type TEXT NOT NULL CHECK (type IN ('tag', 'ordered_tag', 'number', 'string', 'markdown', 'image')),
	description TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE,
	UNIQUE(domain_id, name)
);

-- 노드 속성 값 테이블
CREATE TABLE IF NOT EXISTS node_attributes (
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
CREATE TABLE IF NOT EXISTS node_connections (
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

-- 템플릿 테이블 (노드와 유사하지만 템플릿용 데이터 저장)
CREATE TABLE IF NOT EXISTS templates (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	domain_id INTEGER NOT NULL,
	template_data TEXT NOT NULL, -- JSON 또는 다른 형식의 템플릿 데이터
	title TEXT,
	description TEXT,
	is_active BOOLEAN DEFAULT TRUE,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE,
	UNIQUE(name, domain_id)
);

-- 템플릿 속성 값 테이블
CREATE TABLE IF NOT EXISTS template_attributes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	template_id INTEGER NOT NULL,
	attribute_id INTEGER NOT NULL,
	value TEXT NOT NULL,
	order_index INTEGER, -- 순서가 중요한 태그용
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (template_id) REFERENCES templates(id) ON DELETE CASCADE,
	FOREIGN KEY (attribute_id) REFERENCES attributes(id) ON DELETE CASCADE
);

-- 노드 구독 테이블 (외부 서비스 알림)
CREATE TABLE IF NOT EXISTS node_subscriptions (
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

-- 의존성 타입 레지스트리
CREATE TABLE IF NOT EXISTS dependency_types (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	type_name TEXT NOT NULL UNIQUE,
	category TEXT NOT NULL, -- 'structural', 'behavioral', 'data'
	cascade_delete BOOLEAN DEFAULT FALSE,
	cascade_update BOOLEAN DEFAULT FALSE,
	validation_required BOOLEAN DEFAULT TRUE,
	metadata_schema TEXT, -- JSON schema for type-specific metadata
	description TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 노드 의존성 테이블 (엔터프라이즈급 의존성 관리)
CREATE TABLE IF NOT EXISTS node_dependencies (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	dependent_node_id INTEGER NOT NULL,
	dependency_node_id INTEGER NOT NULL,
	dependency_type_id INTEGER NOT NULL,
	strength INTEGER DEFAULT 50, -- 0-100, dependency strength
	priority INTEGER DEFAULT 50, -- 0-100, resolution priority
	metadata TEXT, -- JSON: type-specific metadata
	version_constraint TEXT, -- Semantic versioning constraint
	is_required BOOLEAN DEFAULT TRUE,
	is_active BOOLEAN DEFAULT TRUE,
	valid_from DATETIME DEFAULT CURRENT_TIMESTAMP,
	valid_until DATETIME,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	created_by TEXT,
	FOREIGN KEY (dependent_node_id) REFERENCES nodes(id) ON DELETE CASCADE,
	FOREIGN KEY (dependency_node_id) REFERENCES nodes(id) ON DELETE CASCADE,
	FOREIGN KEY (dependency_type_id) REFERENCES dependency_types(id),
	UNIQUE(dependent_node_id, dependency_node_id, dependency_type_id, valid_from)
);

-- 의존성 히스토리 추적
CREATE TABLE IF NOT EXISTS dependency_history (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	dependency_id INTEGER NOT NULL,
	action TEXT NOT NULL, -- 'created', 'updated', 'deleted', 'activated', 'deactivated'
	previous_state TEXT, -- JSON: previous dependency state
	new_state TEXT, -- JSON: new dependency state
	change_reason TEXT,
	changed_by TEXT,
	changed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (dependency_id) REFERENCES node_dependencies(id)
);

-- 의존성 그래프 캐시
CREATE TABLE IF NOT EXISTS dependency_graph_cache (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	node_id INTEGER NOT NULL,
	graph_data TEXT NOT NULL, -- JSON: pre-computed dependency graph
	depth INTEGER DEFAULT 0, -- Max depth in dependency tree
	total_dependencies INTEGER DEFAULT 0,
	total_dependents INTEGER DEFAULT 0,
	has_circular BOOLEAN DEFAULT FALSE,
	computed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	expires_at DATETIME,
	FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE
);

-- 의존성 검증 규칙
CREATE TABLE IF NOT EXISTS dependency_rules (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	domain_id INTEGER,
	rule_name TEXT NOT NULL,
	rule_type TEXT NOT NULL, -- 'circular_prevention', 'max_depth', 'type_compatibility'
	rule_config TEXT NOT NULL, -- JSON: rule configuration
	is_active BOOLEAN DEFAULT TRUE,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE
);

-- 의존성 영향 분석 결과
CREATE TABLE IF NOT EXISTS dependency_impact_analysis (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	source_node_id INTEGER NOT NULL,
	impact_type TEXT NOT NULL, -- 'delete', 'update', 'version_change'
	affected_nodes TEXT NOT NULL, -- JSON: array of affected node IDs with impact details
	impact_score INTEGER, -- 0-100, overall impact severity
	analysis_metadata TEXT, -- JSON: detailed analysis results
	analyzed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (source_node_id) REFERENCES nodes(id) ON DELETE CASCADE
);

-- 노드 이벤트 로그 테이블
CREATE TABLE IF NOT EXISTS node_events (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	node_id INTEGER NOT NULL,
	event_type TEXT NOT NULL,             -- 'created', 'updated', 'deleted', 'attribute_changed'
	event_data TEXT,                      -- JSON: 이벤트 상세 데이터
	occurred_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	processed_at DATETIME,                -- 처리 완료 시간
	FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE
);

-- 인덱스 생성
CREATE INDEX IF NOT EXISTS idx_nodes_domain ON nodes(domain_id);
CREATE INDEX IF NOT EXISTS idx_nodes_content ON nodes(content);
CREATE INDEX IF NOT EXISTS idx_attributes_domain ON attributes(domain_id);
CREATE INDEX IF NOT EXISTS idx_node_attributes_node ON node_attributes(node_id);
CREATE INDEX IF NOT EXISTS idx_node_attributes_attribute ON node_attributes(attribute_id);
CREATE INDEX IF NOT EXISTS idx_node_connections_source ON node_connections(source_node_id);
CREATE INDEX IF NOT EXISTS idx_node_connections_target ON node_connections(target_node_id);
CREATE INDEX IF NOT EXISTS idx_node_subscriptions_node ON node_subscriptions(subscribed_node_id);
CREATE INDEX IF NOT EXISTS idx_node_subscriptions_service ON node_subscriptions(subscriber_service);

-- 템플릿 인덱스
CREATE INDEX IF NOT EXISTS idx_templates_domain ON templates(domain_id);
CREATE INDEX IF NOT EXISTS idx_templates_name ON templates(name);
CREATE INDEX IF NOT EXISTS idx_templates_active ON templates(is_active);
CREATE INDEX IF NOT EXISTS idx_template_attributes_template ON template_attributes(template_id);
CREATE INDEX IF NOT EXISTS idx_template_attributes_attribute ON template_attributes(attribute_id);

-- 의존성 인덱스 (성능 최적화)
CREATE INDEX IF NOT EXISTS idx_node_dependencies_dependent ON node_dependencies(dependent_node_id);
CREATE INDEX IF NOT EXISTS idx_node_dependencies_dependency ON node_dependencies(dependency_node_id);
CREATE INDEX IF NOT EXISTS idx_node_dependencies_type ON node_dependencies(dependency_type_id);
CREATE INDEX IF NOT EXISTS idx_node_dependencies_active ON node_dependencies(is_active);
CREATE INDEX IF NOT EXISTS idx_node_dependencies_valid_from ON node_dependencies(valid_from);
CREATE INDEX IF NOT EXISTS idx_node_dependencies_valid_until ON node_dependencies(valid_until);
CREATE INDEX IF NOT EXISTS idx_node_dependencies_strength ON node_dependencies(strength);
CREATE INDEX IF NOT EXISTS idx_node_dependencies_priority ON node_dependencies(priority);

-- 의존성 히스토리 인덱스
CREATE INDEX IF NOT EXISTS idx_deps_history_dep ON dependency_history(dependency_id);
CREATE INDEX IF NOT EXISTS idx_deps_history_action ON dependency_history(action);
CREATE INDEX IF NOT EXISTS idx_deps_history_changed_at ON dependency_history(changed_at);

-- 의존성 그래프 캐시 인덱스
CREATE INDEX IF NOT EXISTS idx_deps_cache_node ON dependency_graph_cache(node_id);
CREATE INDEX IF NOT EXISTS idx_deps_cache_expires ON dependency_graph_cache(expires_at);
CREATE INDEX IF NOT EXISTS idx_deps_cache_computed ON dependency_graph_cache(computed_at);

-- 의존성 영향 분석 인덱스
CREATE INDEX IF NOT EXISTS idx_deps_impact_source ON dependency_impact_analysis(source_node_id);
CREATE INDEX IF NOT EXISTS idx_deps_impact_type ON dependency_impact_analysis(impact_type);
CREATE INDEX IF NOT EXISTS idx_deps_impact_analyzed ON dependency_impact_analysis(analyzed_at);

-- 의존성 타입 인덱스
CREATE INDEX IF NOT EXISTS idx_deps_types_name ON dependency_types(type_name);
CREATE INDEX IF NOT EXISTS idx_deps_types_category ON dependency_types(category);

-- 의존성 규칙 인덱스
CREATE INDEX IF NOT EXISTS idx_deps_rules_domain ON dependency_rules(domain_id);
CREATE INDEX IF NOT EXISTS idx_deps_rules_type ON dependency_rules(rule_type);
CREATE INDEX IF NOT EXISTS idx_deps_rules_active ON dependency_rules(is_active);

-- 이벤트 테이블 인덱스
CREATE INDEX IF NOT EXISTS idx_events_node ON node_events(node_id);
CREATE INDEX IF NOT EXISTS idx_events_type ON node_events(event_type);
CREATE INDEX IF NOT EXISTS idx_events_occurred ON node_events(occurred_at);
CREATE INDEX IF NOT EXISTS idx_events_unprocessed ON node_events(processed_at) WHERE processed_at IS NULL;

-- 트리거: updated_at 자동 업데이트
CREATE TRIGGER IF NOT EXISTS domains_updated_at 
	AFTER UPDATE ON domains 
	FOR EACH ROW 
	BEGIN 
		UPDATE domains SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;

CREATE TRIGGER IF NOT EXISTS nodes_updated_at 
	AFTER UPDATE ON nodes 
	FOR EACH ROW 
	BEGIN 
		UPDATE nodes SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;

CREATE TRIGGER IF NOT EXISTS node_subscriptions_updated_at 
	AFTER UPDATE ON node_subscriptions 
	FOR EACH ROW 
	BEGIN 
		UPDATE node_subscriptions SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;

CREATE TRIGGER IF NOT EXISTS templates_updated_at 
	AFTER UPDATE ON templates 
	FOR EACH ROW 
	BEGIN 
		UPDATE templates SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;

-- 기본 의존성 타입 데이터 초기화
INSERT OR IGNORE INTO dependency_types (type_name, category, cascade_delete, cascade_update, validation_required, description) VALUES
	('hard', 'structural', true, true, true, 'Strong coupling dependency with cascading operations'),
	('soft', 'structural', false, false, true, 'Loose coupling dependency without cascading'),
	('reference', 'structural', false, false, false, 'Informational reference link only'),
	('runtime', 'behavioral', false, true, true, 'Required at runtime execution'),
	('compile', 'behavioral', false, false, true, 'Required at build/compile time'),
	('optional', 'behavioral', false, false, false, 'Optional enhancement dependency'),
	('sync', 'data', false, true, true, 'Synchronous data dependency'),
	('async', 'data', false, false, false, 'Asynchronous data dependency');