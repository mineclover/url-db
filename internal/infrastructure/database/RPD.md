# 데이터베이스 모듈 RPD

## 참조 문서
- [schema.sql](../../schema.sql) - 데이터베이스 스키마
- [docs/spec/error-codes.md](../../docs/spec/error-codes.md) - 에러 코드 정의

## 요구사항 분석

### 기능 요구사항
1. **데이터베이스 초기화**: SQLite 데이터베이스 연결 및 설정
2. **스키마 관리**: 테이블 생성, 인덱스 생성, 트리거 설정
3. **마이그레이션**: 스키마 버전 관리 및 업그레이드
4. **연결 관리**: 연결 풀링, 재연결, 타임아웃 처리
5. **트랜잭션 관리**: 트랜잭션 시작, 커밋, 롤백
6. **헬스 체크**: 데이터베이스 상태 확인

### 비기능 요구사항
- 연결 풀링으로 성능 최적화
- 트랜잭션 안전성 보장
- 에러 처리 및 로깅
- 테스트 환경 지원

## 데이터베이스 설계

### 테이블 구조
1. **domains**: 도메인 관리
2. **nodes**: 노드(URL) 관리
3. **attributes**: 속성 정의
4. **node_attributes**: 노드 속성 값
5. **node_connections**: 노드 간 연결

### 제약 조건
- UNIQUE 제약: 도메인 이름, 노드 content+domain_id
- FOREIGN KEY: 참조 무결성
- CHECK 제약: 속성 타입 검증

### 인덱스 전략
- 기본 키 인덱스: 자동 생성
- 외래 키 인덱스: 조인 성능 향상
- 검색 인덱스: 컨텐츠 검색 최적화

## 아키텍처 설계

### 계층 구조
```
Application -> Database -> SQLite
```

### 주요 컴포넌트
1. **Database**: 데이터베이스 연결 관리
2. **Migration**: 스키마 마이그레이션
3. **Transaction**: 트랜잭션 관리
4. **Health**: 헬스 체크

## 구현 계획

### Phase 1: Core Database
- [ ] Database 구조체 정의
- [ ] SQLite 연결 설정
- [ ] 기본 설정 및 옵션
- [ ] 단위 테스트 작성

### Phase 2: Schema Management
- [ ] 스키마 생성 로직
- [ ] 인덱스 생성 로직
- [ ] 트리거 설정 로직
- [ ] 스키마 테스트 작성

### Phase 3: Migration System
- [ ] 마이그레이션 프레임워크
- [ ] 버전 관리 시스템
- [ ] 롤백 기능
- [ ] 마이그레이션 테스트

### Phase 4: Advanced Features
- [ ] 연결 풀링 구현
- [ ] 트랜잭션 헬퍼
- [ ] 헬스 체크 구현
- [ ] 성능 모니터링

## 데이터베이스 설정

### SQLite 설정
```go
type Config struct {
    URL             string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
    WALMode         bool
    ForeignKeys     bool
    JournalMode     string
    Synchronous     string
}
```

### 기본 설정
- WAL 모드: 동시성 향상
- Foreign Keys: 참조 무결성
- Journal Mode: 데이터 안전성
- Synchronous: 성능/안전성 균형

## 마이그레이션 시스템

### 마이그레이션 파일 구조
```
migrations/
├── 001_initial_schema.up.sql
├── 001_initial_schema.down.sql
├── 002_add_indexes.up.sql
├── 002_add_indexes.down.sql
└── ...
```

### 마이그레이션 관리
- 버전 테이블: 현재 스키마 버전 추적
- 체크섬 검증: 마이그레이션 파일 무결성
- 롤백 지원: 이전 버전으로 복원

## 트랜잭션 관리

### 트랜잭션 헬퍼
```go
func (db *Database) WithTransaction(fn func(*sql.Tx) error) error {
    tx, err := db.Begin()
    if err != nil {
        return err
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
```

### 트랜잭션 격리 수준
- READ_UNCOMMITTED: 최고 성능
- READ_COMMITTED: 기본 설정
- SERIALIZABLE: 최고 안전성

## 에러 처리

### 데이터베이스 에러 코드
- `DATABASE_CONNECTION_FAILED`: 연결 실패
- `DATABASE_MIGRATION_FAILED`: 마이그레이션 실패
- `DATABASE_TRANSACTION_FAILED`: 트랜잭션 실패
- `DATABASE_CONSTRAINT_VIOLATION`: 제약 조건 위반
- `DATABASE_TIMEOUT`: 타임아웃

### SQLite 에러 매핑
- SQLITE_CONSTRAINT_UNIQUE → UNIQUE_CONSTRAINT_VIOLATION
- SQLITE_CONSTRAINT_FOREIGN_KEY → FOREIGN_KEY_CONSTRAINT_VIOLATION
- SQLITE_BUSY → DATABASE_BUSY
- SQLITE_LOCKED → DATABASE_LOCKED

## 테스트 전략

### 단위 테스트
- 데이터베이스 연결 테스트
- 스키마 생성 테스트
- 마이그레이션 테스트
- 트랜잭션 테스트

### 통합 테스트
- 실제 데이터베이스 연동 테스트
- 동시성 테스트
- 성능 테스트

### 테스트 환경
- 인메모리 SQLite: 빠른 단위 테스트
- 임시 파일 SQLite: 통합 테스트
- 실제 파일 SQLite: 성능 테스트

## 파일 구조
```
internal/database/
├── RPD.md
├── database.go         # Database 구조체
├── database_test.go    # Database 테스트
├── config.go           # 설정 관리
├── migration.go        # 마이그레이션 관리
├── migration_test.go   # 마이그레이션 테스트
├── transaction.go      # 트랜잭션 헬퍼
├── transaction_test.go # 트랜잭션 테스트
├── health.go           # 헬스 체크
├── health_test.go      # 헬스 체크 테스트
├── errors.go           # 데이터베이스 에러
├── migrations/         # 마이그레이션 파일
│   ├── 001_initial_schema.up.sql
│   ├── 001_initial_schema.down.sql
│   └── ...
└── testdata/           # 테스트 데이터
    ├── test_schema.sql
    └── ...
```

## 성능 최적화

### 연결 풀링
- 최대 연결 수: CPU 코어 * 2
- 유휴 연결 수: 최대 연결 수 / 2
- 연결 수명: 1시간

### 쿼리 최적화
- 프리페어드 스테이트먼트 사용
- 인덱스 힌트 제공
- 쿼리 플랜 분석

### 메모리 관리
- 페이지 크기 최적화
- 캐시 크기 조정
- 메모리 매핑 활용

## 보안 고려사항

### SQL 인젝션 방지
- 프리페어드 스테이트먼트 강제
- 사용자 입력 검증
- 쿼리 빌더 사용

### 접근 제어
- 파일 권한 설정
- 네트워크 접근 제한
- 암호화 지원

## 모니터링

### 메트릭 수집
- 연결 수 모니터링
- 쿼리 성능 측정
- 에러 발생률 추적

### 로깅
- 쿼리 로깅 (개발 환경)
- 에러 로깅
- 성능 로깅

## 의존성
- `database/sql`: 표준 데이터베이스 인터페이스
- `github.com/mattn/go-sqlite3`: SQLite 드라이버
- `github.com/golang-migrate/migrate/v4`: 마이그레이션 도구
- `github.com/stretchr/testify`: 테스트 유틸리티

## 설정 예제

### 개발 환경
```go
config := &Config{
    URL:             "file:./dev.db",
    MaxOpenConns:    10,
    MaxIdleConns:    5,
    ConnMaxLifetime: time.Hour,
    WALMode:         true,
    ForeignKeys:     true,
    JournalMode:     "WAL",
    Synchronous:     "NORMAL",
}
```

### 프로덕션 환경
```go
config := &Config{
    URL:             "file:./prod.db",
    MaxOpenConns:    100,
    MaxIdleConns:    50,
    ConnMaxLifetime: time.Hour,
    WALMode:         true,
    ForeignKeys:     true,
    JournalMode:     "WAL",
    Synchronous:     "FULL",
}
```