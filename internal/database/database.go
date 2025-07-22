package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

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
	schema := `
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

	-- 노드 의존성 테이블 (노드 간 의존성 관리)
	CREATE TABLE IF NOT EXISTS node_dependencies (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		dependent_node_id INTEGER NOT NULL,
		dependency_node_id INTEGER NOT NULL,
		dependency_type TEXT NOT NULL,
		description TEXT,
		is_required BOOLEAN DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (dependent_node_id) REFERENCES nodes(id) ON DELETE CASCADE,
		FOREIGN KEY (dependency_node_id) REFERENCES nodes(id) ON DELETE CASCADE,
		UNIQUE(dependent_node_id, dependency_node_id, dependency_type)
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
	CREATE INDEX IF NOT EXISTS idx_node_dependencies_dependent ON node_dependencies(dependent_node_id);
	CREATE INDEX IF NOT EXISTS idx_node_dependencies_dependency ON node_dependencies(dependency_node_id);

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
	`

	if _, err := d.db.Exec(schema); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
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
