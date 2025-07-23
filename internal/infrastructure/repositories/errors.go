package repositories

import (
	"database/sql"
	"errors"
	"strings"
)

// 리포지토리 에러 정의
var (
	ErrDomainNotFound         = errors.New("domain not found")
	ErrNodeNotFound           = errors.New("node not found")
	ErrAttributeNotFound      = errors.New("attribute not found")
	ErrNodeAttributeNotFound  = errors.New("node attribute not found")
	ErrNodeConnectionNotFound = errors.New("node connection not found")
	ErrDuplicateEntry         = errors.New("duplicate entry")
	ErrForeignKeyConstraint   = errors.New("foreign key constraint violation")
	ErrInvalidInput           = errors.New("invalid input")
	ErrConnectionTimeout      = errors.New("database connection timeout")
	ErrTransactionFailed      = errors.New("transaction failed")
)

// SQLiteError 는 SQLite 관련 에러를 처리하는 구조체입니다.
type SQLiteError struct {
	Code    int
	Message string
	Query   string
}

func (e SQLiteError) Error() string {
	return e.Message
}

// MapSQLiteError 는 SQLite 에러를 도메인 에러로 매핑합니다.
func MapSQLiteError(err error) error {
	if err == nil {
		return nil
	}

	// sql.ErrNoRows 처리
	if err == sql.ErrNoRows {
		return ErrDomainNotFound // 기본값, 각 리포지토리에서 적절히 변경
	}

	// SQLite 에러 처리 (문자열 기반)
	errStr := err.Error()
	if strings.Contains(errStr, "UNIQUE constraint failed") {
		return ErrDuplicateEntry
	}
	if strings.Contains(errStr, "FOREIGN KEY constraint failed") {
		return ErrForeignKeyConstraint
	}
	if strings.Contains(errStr, "CHECK constraint failed") ||
		strings.Contains(errStr, "NOT NULL constraint failed") {
		return ErrInvalidInput
	}
	if strings.Contains(errStr, "database is locked") ||
		strings.Contains(errStr, "database is busy") {
		return ErrConnectionTimeout
	}

	return err
}

// IsNotFoundError 는 해당 에러가 "찾을 수 없음" 에러인지 확인합니다.
func IsNotFoundError(err error) bool {
	return err == ErrDomainNotFound ||
		err == ErrNodeNotFound ||
		err == ErrAttributeNotFound ||
		err == ErrNodeAttributeNotFound ||
		err == ErrNodeConnectionNotFound
}

// IsDuplicateError 는 해당 에러가 중복 에러인지 확인합니다.
func IsDuplicateError(err error) bool {
	return err == ErrDuplicateEntry
}

// IsForeignKeyError 는 해당 에러가 외래키 제약 조건 에러인지 확인합니다.
func IsForeignKeyError(err error) bool {
	return err == ErrForeignKeyConstraint
}

// IsInvalidInputError 는 해당 에러가 잘못된 입력 에러인지 확인합니다.
func IsInvalidInputError(err error) bool {
	return err == ErrInvalidInput
}

// IsConnectionError 는 해당 에러가 연결 관련 에러인지 확인합니다.
func IsConnectionError(err error) bool {
	return err == ErrConnectionTimeout
}

// IsTransactionError 는 해당 에러가 트랜잭션 관련 에러인지 확인합니다.
func IsTransactionError(err error) bool {
	return err == ErrTransactionFailed
}
