package repositories

import (
	"database/sql"
)

// BaseRepository 는 모든 리포지토리의 공통 베이스 구조체입니다.
type BaseRepository struct {
	db *sql.DB
}

// NewBaseRepository 는 새로운 베이스 리포지토리를 생성합니다.
func NewBaseRepository(db *sql.DB) *BaseRepository {
	return &BaseRepository{
		db: db,
	}
}

// WithTransaction 은 트랜잭션을 사용하여 함수를 실행합니다.
func (r *BaseRepository) WithTransaction(fn func(tx *sql.Tx) error) error {
	tx, err := r.db.Begin()
	if err != nil {
		return MapSQLiteError(err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			// 롤백 실패 시 원래 에러와 함께 로그 (실제 구현에서는 로거 사용)
			return MapSQLiteError(err)
		}
		return MapSQLiteError(err)
	}

	if err := tx.Commit(); err != nil {
		return MapSQLiteError(err)
	}

	return nil
}

// GetDB 는 데이터베이스 연결을 반환합니다.
func (r *BaseRepository) GetDB() *sql.DB {
	return r.db
}

// Ping 은 데이터베이스 연결을 확인합니다.
func (r *BaseRepository) Ping() error {
	return r.db.Ping()
}

// Close 는 데이터베이스 연결을 닫습니다.
func (r *BaseRepository) Close() error {
	return r.db.Close()
}

// ExecuteInTransaction 은 트랜잭션 내에서 쿼리를 실행합니다.
func (r *BaseRepository) ExecuteInTransaction(tx *sql.Tx, query string, args ...interface{}) (sql.Result, error) {
	result, err := tx.Exec(query, args...)
	if err != nil {
		return nil, MapSQLiteError(err)
	}
	return result, nil
}

// QueryInTransaction 은 트랜잭션 내에서 쿼리를 실행하고 결과를 반환합니다.
func (r *BaseRepository) QueryInTransaction(tx *sql.Tx, query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := tx.Query(query, args...)
	if err != nil {
		return nil, MapSQLiteError(err)
	}
	return rows, nil
}

// QueryRowInTransaction 은 트랜잭션 내에서 단일 행 쿼리를 실행합니다.
func (r *BaseRepository) QueryRowInTransaction(tx *sql.Tx, query string, args ...interface{}) *sql.Row {
	return tx.QueryRow(query, args...)
}

// PrepareStatement 는 프리페어드 스테이트먼트를 생성합니다.
func (r *BaseRepository) PrepareStatement(query string) (*sql.Stmt, error) {
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, MapSQLiteError(err)
	}
	return stmt, nil
}

// PrepareStatementInTransaction 은 트랜잭션 내에서 프리페어드 스테이트먼트를 생성합니다.
func (r *BaseRepository) PrepareStatementInTransaction(tx *sql.Tx, query string) (*sql.Stmt, error) {
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, MapSQLiteError(err)
	}
	return stmt, nil
}

// Execute 는 쿼리를 실행합니다.
func (r *BaseRepository) Execute(query string, args ...interface{}) (sql.Result, error) {
	result, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, MapSQLiteError(err)
	}
	return result, nil
}

// Query 는 쿼리를 실행하고 결과를 반환합니다.
func (r *BaseRepository) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, MapSQLiteError(err)
	}
	return rows, nil
}

// QueryRow 는 단일 행 쿼리를 실행합니다.
func (r *BaseRepository) QueryRow(query string, args ...interface{}) *sql.Row {
	return r.db.QueryRow(query, args...)
}
