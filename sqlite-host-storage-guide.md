# SQLite를 컨테이너 밖(호스트)에 저장하기

## 방법별 비교

### 1. 🐳 **Docker Volume (기본)**
```bash
docker run -it --rm -v url-db-data:/data asfdassdssa/url-db:latest
```

**장점**:
- Docker가 관리하므로 안전
- 컨테이너 간 쉬운 데이터 공유
- 백업/복원이 표준화됨

**단점**:
- 호스트에서 직접 접근 어려움
- 파일 위치가 숨겨져 있음

### 2. 📁 **로컬 디렉토리 마운트 (권장)**
```bash
docker run -it --rm -v $(pwd)/database:/data asfdassdssa/url-db:latest
```

**장점**:
- ✅ 호스트에서 직접 파일 접근 가능
- ✅ 백업이 간단 (파일 복사)
- ✅ 다른 도구로 DB 직접 조작 가능
- ✅ 파일 위치가 명확함

**단점**:
- 권한 문제 가능성
- 경로 관리 필요

### 3. 🏠 **절대 경로 마운트**
```bash
docker run -it --rm -v /Users/username/url-db:/data asfdassdssa/url-db:latest
```

**장점**:
- 명확한 위치 지정
- 여러 프로젝트에서 공유 가능

**단점**:
- 하드코딩된 경로
- 이식성 떨어짐

## 실제 사용 시나리오

### Claude Desktop 설정 (로컬 폴더)

```json
{
  "mcpServers": {
    "url-db": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-v", "/Users/junwoobang/Documents/url-db:/data",
        "asfdassdssa/url-db:latest"
      ]
    }
  }
}
```

### 개발/테스트용 설정

```bash
# 프로젝트 폴더에 데이터베이스 생성
mkdir -p ./database
docker run -it --rm -v $(pwd)/database:/data asfdassdssa/url-db:latest

# 데이터베이스 직접 확인
sqlite3 ./database/url-db.sqlite "SELECT * FROM domains;"
```

### 프로덕션용 설정

```bash
# 전용 데이터 디렉토리 생성
sudo mkdir -p /opt/url-db/data
sudo chown $USER:$USER /opt/url-db/data

# 프로덕션 실행
docker run -d --name url-db-prod \
  -v /opt/url-db/data:/data \
  -p 8080:8080 \
  --restart unless-stopped \
  asfdassdssa/url-db:latest \
  -port=8080 -db-path=/data/url-db.sqlite
```

## 데이터베이스 관리

### 백업
```bash
# 단순 파일 복사
cp ./database/url-db.sqlite ./backups/url-db-$(date +%Y%m%d).sqlite

# SQLite 덤프
sqlite3 ./database/url-db.sqlite .dump > backup.sql
```

### 복원
```bash
# 파일 복사로 복원
cp ./backups/url-db-20250724.sqlite ./database/url-db.sqlite

# SQL 덤프로 복원
sqlite3 ./database/url-db.sqlite < backup.sql
```

### 데이터베이스 분석
```bash
# 테이블 목록
sqlite3 ./database/url-db.sqlite ".tables"

# 스키마 확인
sqlite3 ./database/url-db.sqlite ".schema"

# 데이터 조회
sqlite3 ./database/url-db.sqlite "SELECT * FROM domains LIMIT 10;"
```

## 권한 문제 해결

### macOS/Linux
```bash
# 현재 사용자로 소유권 설정
chown -R $USER:$USER ./database

# 권한 설정
chmod 755 ./database
chmod 644 ./database/url-db.sqlite
```

### Docker 사용자 매핑
```bash
# 현재 사용자 ID로 실행
docker run -it --rm \
  -v $(pwd)/database:/data \
  -u $(id -u):$(id -g) \
  asfdassdssa/url-db:latest
```

## 보안 고려사항

1. **파일 권한**: 데이터베이스 파일을 적절히 보호
2. **백업 암호화**: 민감한 데이터가 있다면 백업 암호화
3. **접근 제한**: 필요한 사용자만 디렉토리 접근 허용

## 권장사항

### 개발 환경
```bash
# 프로젝트 루트에 database 폴더 생성
mkdir -p ./database
docker run -it --rm -v $(pwd)/database:/data asfdassdssa/url-db:latest
```

### Claude Desktop
```json
{
  "mcpServers": {
    "url-db": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-v", "~/url-db-data:/data",
        "asfdassdssa/url-db:latest"
      ]
    }
  }
}
```

### 프로덕션 환경
```bash
# 전용 데이터 디렉토리 사용
docker run -d \
  -v /opt/url-db/data:/data \
  --restart unless-stopped \
  asfdassdssa/url-db:latest
```

이렇게 하면 SQLite 데이터베이스를 컨테이너 밖에서 완전히 관리할 수 있습니다!