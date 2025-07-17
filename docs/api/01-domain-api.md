# 도메인 관리 API 엔드포인트

## 개요
도메인 폴더 생성, 조회, 수정, 삭제 기능을 제공하는 REST API

## 엔드포인트 목록

### 1. 도메인 생성
- **POST** `/api/domains`
- **요청 본문**:
```json
{
  "name": "programming",
  "description": "Programming related URLs"
}
```
- **응답 (201)**:
```json
{
  "id": 1,
  "name": "programming",
  "description": "Programming related URLs",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### 2. 도메인 목록 조회
- **GET** `/api/domains`
- **쿼리 파라미터**:
  - `page` (optional): 페이지 번호 (기본값: 1)
  - `size` (optional): 페이지 크기 (기본값: 20, 최대: 100)
- **응답 (200)**:
```json
{
  "domains": [
    {
      "id": 1,
      "name": "programming",
      "description": "Programming related URLs",
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  ],
  "total_count": 1,
  "page": 1,
  "size": 20,
  "total_pages": 1
}
```

### 3. 도메인 조회
- **GET** `/api/domains/{id}`
- **응답 (200)**:
```json
{
  "id": 1,
  "name": "programming",
  "description": "Programming related URLs",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### 4. 도메인 수정
- **PUT** `/api/domains/{id}`
- **요청 본문**:
```json
{
  "name": "updated-programming",
  "description": "Updated description"
}
```
- **응답 (200)**:
```json
{
  "id": 1,
  "name": "updated-programming",
  "description": "Updated description",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T01:00:00Z"
}
```

### 5. 도메인 삭제
- **DELETE** `/api/domains/{id}`
- **응답 (204)**: 본문 없음

## 에러 응답

> 에러 응답 형식: [`../spec/error-codes.md`](../spec/error-codes.md)  
> 도메인 관련 에러: [`../spec/domain-errors.md`](../spec/domain-errors.md)

## 검증 규칙
- `name`: 필수, 1-255자, 영숫자와 하이픈만 허용
- `description`: 선택, 최대 1000자