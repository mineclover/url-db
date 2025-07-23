package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// Schema file path relative to project root
const schemaFilePath = "schema.sql"

type Database struct {
	db     *sql.DB
	config *Config
}

func New(config *Config) (*Database, error) {
	if config == nil {
		config = DefaultConfig()
	}

	db, err := sql.Open("sqlite3", config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := configureDatabase(db, config); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to configure database: %w", err)
	}

	database := &Database{
		db:     db,
		config: config,
	}

	if err := database.createSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return database, nil
}

func configureDatabase(db *sql.DB, config *Config) error {
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	pragmas := []string{
		fmt.Sprintf("PRAGMA journal_mode = %s", config.JournalMode),
		fmt.Sprintf("PRAGMA synchronous = %s", config.Synchronous),
	}

	if config.ForeignKeys {
		pragmas = append(pragmas, "PRAGMA foreign_keys = ON")
	}

	if config.WALMode {
		pragmas = append(pragmas, "PRAGMA journal_mode = WAL")
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return fmt.Errorf("failed to execute pragma %s: %w", pragma, err)
		}
	}

	return nil
}

func (d *Database) createSchema() error {
	// Load schema from external file
	schema, err := d.loadSchemaFromFile()
	if err != nil {
		// Fallback to inline schema if file not found
		fmt.Printf("Warning: Could not load schema.sql file, using inline schema: %v\n", err)
		schema = d.getInlineSchema()
	}

	if _, err := d.db.Exec(schema); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}

// loadSchemaFromFile loads schema from external schema.sql file
func (d *Database) loadSchemaFromFile() (string, error) {
	// Find project root by looking for go.mod
	projectRoot, err := findProjectRoot()
	if err != nil {
		return "", fmt.Errorf("could not find project root: %w", err)
	}

	schemaPath := filepath.Join(projectRoot, schemaFilePath)
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return "", fmt.Errorf("could not read schema file %s: %w", schemaPath, err)
	}

	return string(schemaBytes), nil
}

// findProjectRoot finds the project root by looking for go.mod file
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Walk up the directory tree looking for go.mod
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached root directory
		}
		dir = parent
	}

	return "", fmt.Errorf("go.mod not found")
}

// getInlineSchema returns the inline schema as fallback
func (d *Database) getInlineSchema() string {
	return `
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
		FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE,
		UNIQUE(domain_id, name)
	);

	-- 노드 속성 값 테이블
	CREATE TABLE IF NOT EXISTS node_attributes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		node_id INTEGER NOT NULL,
		attribute_id INTEGER NOT NULL,
		value TEXT NOT NULL,
		order_index INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE,
		FOREIGN KEY (attribute_id) REFERENCES attributes(id) ON DELETE CASCADE
	);

	-- 노드 간 연결관계 테이블
	CREATE TABLE IF NOT EXISTS node_connections (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		source_node_id INTEGER NOT NULL,
		target_node_id INTEGER NOT NULL,
		relationship_type TEXT NOT NULL,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (source_node_id) REFERENCES nodes(id) ON DELETE CASCADE,
		FOREIGN KEY (target_node_id) REFERENCES nodes(id) ON DELETE CASCADE,
		UNIQUE(source_node_id, target_node_id, relationship_type)
	);

	-- 노드 구독 테이블 (이벤트 구독 관리)
	CREATE TABLE IF NOT EXISTS node_subscriptions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		subscriber_service TEXT NOT NULL,
		subscriber_endpoint TEXT,
		subscribed_node_id INTEGER NOT NULL,
		event_types TEXT NOT NULL, -- JSON array of event types
		filter_conditions TEXT,    -- JSON object for filter conditions
		is_active BOOLEAN DEFAULT 1,
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

	-- 노드 의존성 테이블 (기존 호환성용)
	CREATE TABLE IF NOT EXISTS node_dependencies (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		dependent_node_id INTEGER NOT NULL,
		dependency_node_id INTEGER NOT NULL,
		dependency_type TEXT NOT NULL,
		cascade_delete BOOLEAN DEFAULT FALSE,
		cascade_update BOOLEAN DEFAULT FALSE,
		metadata TEXT, -- JSON metadata
		description TEXT,
		is_required BOOLEAN DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (dependent_node_id) REFERENCES nodes(id) ON DELETE CASCADE,
		FOREIGN KEY (dependency_node_id) REFERENCES nodes(id) ON DELETE CASCADE,
		UNIQUE(dependent_node_id, dependency_node_id, dependency_type)
	);

	-- 고도화된 노드 의존성 테이블 V2
	CREATE TABLE IF NOT EXISTS node_dependencies_v2 (
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
		FOREIGN KEY (dependency_id) REFERENCES node_dependencies_v2(id)
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
	-- 기존 의존성 인덱스
	CREATE INDEX IF NOT EXISTS idx_node_dependencies_dependent ON node_dependencies(dependent_node_id);
	CREATE INDEX IF NOT EXISTS idx_node_dependencies_dependency ON node_dependencies(dependency_node_id);
	
	-- 고도화된 의존성 인덱스
	CREATE INDEX IF NOT EXISTS idx_deps_v2_dependent ON node_dependencies_v2(dependent_node_id);
	CREATE INDEX IF NOT EXISTS idx_deps_v2_dependency ON node_dependencies_v2(dependency_node_id);
	CREATE INDEX IF NOT EXISTS idx_deps_v2_type ON node_dependencies_v2(dependency_type_id);
	CREATE INDEX IF NOT EXISTS idx_deps_v2_active ON node_dependencies_v2(is_active);
	CREATE INDEX IF NOT EXISTS idx_deps_v2_valid_from ON node_dependencies_v2(valid_from);
	CREATE INDEX IF NOT EXISTS idx_deps_v2_valid_until ON node_dependencies_v2(valid_until);
	CREATE INDEX IF NOT EXISTS idx_deps_v2_strength ON node_dependencies_v2(strength);
	CREATE INDEX IF NOT EXISTS idx_deps_v2_priority ON node_dependencies_v2(priority);
	
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

	-- 이벤트 테이블 인덱스
	CREATE INDEX IF NOT EXISTS idx_events_node ON node_events(node_id);
	CREATE INDEX IF NOT EXISTS idx_events_type ON node_events(event_type);
	CREATE INDEX IF NOT EXISTS idx_events_occurred ON node_events(occurred_at);
	CREATE INDEX IF NOT EXISTS idx_events_unprocessed ON node_events(processed_at) WHERE processed_at IS NULL;

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
	`
}

func (d *Database) DB() *sql.DB {
	return d.db
}

func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

func (d *Database) Ping() error {
	return d.db.Ping()
}

func (d *Database) WithTransaction(fn func(*sql.Tx) error) error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// InitDB initializes the database with the given URL
func InitDB(url string) (*Database, error) {
	config := DefaultConfig()
	config.URL = url
	return New(config)
}
