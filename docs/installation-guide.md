# URL-DB 설치 가이드

## 개요
URL-DB는 Go로 작성된 URL 관리 시스템으로, 속성 기반 태깅과 외부 종속성 관리를 지원합니다.

## 시스템 요구사항

### 최소 요구사항
- Go 1.21 이상
- SQLite 3.x
- 메모리: 512MB 이상
- 디스크: 1GB 이상

### 권장 요구사항
- Go 1.21 이상
- SQLite 3.x
- 메모리: 2GB 이상
- 디스크: 10GB 이상

## 설치 방법

### 1. 소스 코드 다운로드
```bash
git clone https://github.com/yourusername/url-db.git
cd url-db
```

### 2. 의존성 설치
```bash
go mod download
```

### 3. 빌드

#### Windows
```bash
# PowerShell 또는 Command Prompt에서
.\build.bat
```

#### Unix/Linux/macOS
```bash
chmod +x build.sh
./build.sh
```

### 4. 데이터베이스 초기화
빌드 시 자동으로 초기화되지만, 수동으로 초기화하려면:

```bash
# 데이터베이스 파일 생성
sqlite3 url-db.db < schema.sql
```

## 실행 방법

### 1. HTTP 서버 모드 (기본)
```bash
# Windows
bin\url-db.exe

# Unix/Linux/macOS
./bin/url-db
```

### 2. MCP stdio 모드
```bash
# Windows
bin\url-db.exe -mcp-mode=stdio

# Unix/Linux/macOS
./bin/url-db -mcp-mode=stdio
```

### 3. 포트 설정
```bash
# 환경변수로 포트 설정
export PORT=9090
./bin/url-db

# 또는 직접 설정
PORT=9090 ./bin/url-db
```

## 환경변수 설정

### 필수 환경변수
```bash
# 데이터베이스 파일 경로
export DATABASE_URL="file:url-db.db"

# 서버 포트 (기본값: 8080)
export PORT=8080

# 합성키용 도구 이름 (기본값: url-db)
export TOOL_NAME=url-db
```

### 선택적 환경변수
```bash
# 로그 레벨 설정
export LOG_LEVEL=info

# CORS 설정
export CORS_ALLOW_ORIGINS="*"
export CORS_ALLOW_METHODS="GET,POST,PUT,DELETE,OPTIONS"
export CORS_ALLOW_HEADERS="Content-Type,Authorization"
```

## 설정 파일 예시

### `.env` 파일
```bash
# 데이터베이스 설정
DATABASE_URL=file:url-db.db

# 서버 설정
PORT=8080
TOOL_NAME=url-db

# 로그 설정
LOG_LEVEL=info

# CORS 설정
CORS_ALLOW_ORIGINS=*
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOW_HEADERS=Content-Type,Authorization
```

## Docker 설치 (선택사항)

### Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o url-db cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite
WORKDIR /root/

COPY --from=builder /app/url-db .
COPY --from=builder /app/schema.sql .

# 데이터베이스 초기화
RUN sqlite3 url-db.db < schema.sql

EXPOSE 8080
CMD ["./url-db"]
```

### docker-compose.yml
```yaml
version: '3.8'

services:
  url-db:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./data:/data
    environment:
      - DATABASE_URL=file:/data/url-db.db
      - PORT=8080
    restart: unless-stopped
```

## 빌드 스크립트

### build.sh (Unix/Linux/macOS)
```bash
#!/bin/bash
set -e

echo "Building URL-DB..."

# 빌드 디렉토리 생성
mkdir -p bin

# Go 빌드
echo "Building Go application..."
go build -o bin/url-db cmd/server/main.go

# 데이터베이스 초기화
echo "Initializing database..."
if [ ! -f url-db.db ]; then
    sqlite3 url-db.db < schema.sql
    echo "Database initialized."
else
    echo "Database already exists."
fi

# 실행 권한 설정
chmod +x bin/url-db

echo "Build completed successfully!"
echo "Run: ./bin/url-db"
```

### build.bat (Windows)
```batch
@echo off
echo Building URL-DB...

:: 빌드 디렉토리 생성
if not exist bin mkdir bin

:: Go 빌드
echo Building Go application...
go build -o bin\url-db.exe cmd\server\main.go

:: 데이터베이스 초기화
echo Initializing database...
if not exist url-db.db (
    sqlite3 url-db.db < schema.sql
    echo Database initialized.
) else (
    echo Database already exists.
)

echo Build completed successfully!
echo Run: bin\url-db.exe
pause
```

## 서비스 등록

### systemd (Linux)
```ini
# /etc/systemd/system/url-db.service
[Unit]
Description=URL Database Service
After=network.target

[Service]
Type=simple
User=urldb
Group=urldb
WorkingDirectory=/opt/url-db
ExecStart=/opt/url-db/bin/url-db
Restart=always
RestartSec=10
Environment=DATABASE_URL=file:/opt/url-db/data/url-db.db
Environment=PORT=8080

[Install]
WantedBy=multi-user.target
```

### 서비스 시작
```bash
sudo systemctl daemon-reload
sudo systemctl enable url-db
sudo systemctl start url-db
sudo systemctl status url-db
```

## 업그레이드 가이드

### 1. 백업
```bash
# 데이터베이스 백업
cp url-db.db url-db.db.backup.$(date +%Y%m%d_%H%M%S)
```

### 2. 새 버전 설치
```bash
# 새 버전 다운로드
git pull origin main

# 빌드
./build.sh  # 또는 build.bat
```

### 3. 스키마 업데이트 (필요시)
```bash
# 스키마 변경사항 적용
sqlite3 url-db.db < migration.sql
```

## 트러블슈팅

### 일반적인 문제

#### 1. 포트 이미 사용 중
```bash
# 포트 사용 중인 프로세스 확인
lsof -i :8080  # macOS/Linux
netstat -ano | findstr :8080  # Windows

# 다른 포트 사용
PORT=9090 ./bin/url-db
```

#### 2. 데이터베이스 권한 문제
```bash
# 데이터베이스 파일 권한 확인
ls -la url-db.db

# 권한 수정
chmod 644 url-db.db
```

#### 3. 메모리 부족
```bash
# 메모리 사용량 확인
free -h  # Linux
top -l 1 | grep PhysMem  # macOS

# 스왑 메모리 설정 (Linux)
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

### 로그 확인

#### 애플리케이션 로그
```bash
# 실시간 로그 확인
tail -f /var/log/url-db.log

# systemd 로그 확인
journalctl -u url-db -f
```

#### 데이터베이스 로그
```bash
# SQLite 로그 활성화
sqlite3 url-db.db
.log stdout
.headers on
.mode column
```

## 성능 최적화

### 1. 데이터베이스 최적화
```sql
-- 인덱스 재구성
REINDEX;

-- 통계 업데이트
ANALYZE;

-- 불필요한 공간 제거
VACUUM;
```

### 2. 메모리 최적화
```bash
# Go 가비지 컬렉션 튜닝
export GOGC=100
export GOMEMLIMIT=2GiB
```

### 3. 파일 시스템 최적화
```bash
# SSD 최적화 (Linux)
echo deadline > /sys/block/sda/queue/scheduler
```

## 모니터링

### 1. 헬스체크
```bash
# HTTP 헬스체크
curl http://localhost:8080/health

# 응답 예시
{
  "status": "healthy",
  "service": "url-db"
}
```

### 2. 메트릭 수집
```bash
# 시스템 메트릭
ps aux | grep url-db
top -p $(pgrep url-db)

# 데이터베이스 메트릭
sqlite3 url-db.db "SELECT name, COUNT(*) FROM sqlite_master WHERE type='table' GROUP BY name;"
```

## 보안 권장사항

### 1. 방화벽 설정
```bash
# Ubuntu/Debian
sudo ufw allow 8080/tcp

# CentOS/RHEL
sudo firewall-cmd --permanent --add-port=8080/tcp
sudo firewall-cmd --reload
```

### 2. 사용자 권한
```bash
# 전용 사용자 생성
sudo useradd -r -s /bin/false urldb
sudo chown -R urldb:urldb /opt/url-db
```

### 3. 데이터베이스 백업
```bash
# 자동 백업 스크립트
#!/bin/bash
BACKUP_DIR="/backup/url-db"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR
cp url-db.db $BACKUP_DIR/url-db_$DATE.db
find $BACKUP_DIR -name "*.db" -mtime +7 -delete
```

## 지원 및 문의

- GitHub Issues: https://github.com/yourusername/url-db/issues
- 문서: https://docs.url-db.com
- 이메일: support@url-db.com