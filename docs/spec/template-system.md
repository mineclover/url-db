# 템플릿 시스템 개요

## 소개
템플릿 시스템은 URL-DB에서 재사용 가능한 데이터 구조를 정의하고 관리하는 기능입니다. 노드(Node)와 유사한 구조를 가지지만, 실제 URL이나 컨텐츠 대신 템플릿 데이터를 저장합니다.

## 주요 특징

### 1. 도메인 기반 관리
- 템플릿은 특정 도메인에 속함
- 도메인의 속성 스키마를 공유
- 도메인별로 템플릿 네임스페이스 분리

### 2. 속성 시스템 통합
- 노드와 동일한 속성 시스템 사용
- 6가지 속성 타입 지원 (tag, ordered_tag, number, string, markdown, image)
- 도메인 스키마 검증 적용

### 3. 활성화 상태 관리
- `is_active` 플래그로 템플릿 활성화/비활성화
- 비활성화된 템플릿은 읽기 전용
- 활성화 상태에 따른 접근 제어

## 데이터베이스 구조

### templates 테이블
```sql
CREATE TABLE IF NOT EXISTS templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    domain_id INTEGER NOT NULL,
    template_data TEXT NOT NULL, -- JSON 형식의 템플릿 데이터
    title TEXT,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE,
    UNIQUE(name, domain_id)
);
```

### template_attributes 테이블
```sql
CREATE TABLE IF NOT EXISTS template_attributes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    template_id INTEGER NOT NULL,
    attribute_id INTEGER NOT NULL,
    value TEXT NOT NULL,
    order_index INTEGER, -- ordered_tag 타입용
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (template_id) REFERENCES templates(id) ON DELETE CASCADE,
    FOREIGN KEY (attribute_id) REFERENCES attributes(id) ON DELETE CASCADE
);
```

## 합성키 규칙

템플릿은 특별한 합성키 형식을 사용합니다:
```
{tool_name}:{domain_name}:template:{id}
```

예시:
- `url-db:products:template:1`
- `url-db:articles:template:42`

## 사용 예시

### 1. 템플릿 생성
```json
{
  "name": "product-template",
  "domain_name": "products",
  "template_data": {
    "layout": "grid",
    "fields": ["name", "price", "category"],
    "validation": {...}
  },
  "title": "기본 상품 템플릿",
  "description": "상품 페이지용 기본 템플릿"
}
```

### 2. 템플릿 속성 설정
```json
{
  "composite_id": "url-db:products:template:1",
  "attributes": [
    {"name": "category", "value": "layout"},
    {"name": "version", "value": "1.0"},
    {"name": "tags", "value": "responsive"},
    {"name": "tags", "value": "mobile-friendly"}
  ]
}
```

### 3. 템플릿 조회
```json
// 응답
{
  "composite_id": "url-db:products:template:1",
  "name": "product-template",
  "domain_name": "products",
  "template_data": {...},
  "title": "기본 상품 템플릿",
  "description": "상품 페이지용 기본 템플릿",
  "is_active": true,
  "attributes": [
    {"name": "category", "value": "layout", "type": "tag"},
    {"name": "version", "value": "1.0", "type": "string"}
  ]
}
```

## 검증 규칙

### 템플릿 검증
1. **이름**: 도메인 내에서 고유해야 함
2. **템플릿 데이터**: 유효한 JSON 형식이어야 함
3. **도메인**: 존재하는 도메인이어야 함

### 속성 검증
1. **도메인 일치**: 템플릿과 속성은 같은 도메인에 속해야 함
2. **타입 검증**: 속성 값은 정의된 타입과 일치해야 함
3. **활성화 상태**: 비활성화된 템플릿의 속성은 수정 불가

## 권한 및 제약사항

### 생성/수정 권한
- 활성화된 템플릿만 속성 추가/수정 가능
- 템플릿 이름과 도메인은 생성 후 변경 불가

### 삭제 제약
- 다른 엔티티에서 참조 중인 템플릿은 삭제 불가
- CASCADE DELETE로 템플릿 삭제 시 관련 속성도 함께 삭제

### 조회 권한
- 비활성화된 템플릿도 조회는 가능 (읽기 전용)
- 도메인 접근 권한에 따라 템플릿 접근 제어

## 성능 최적화

### 인덱스
- `idx_templates_domain`: 도메인별 템플릿 조회
- `idx_templates_name`: 이름으로 템플릿 검색
- `idx_templates_active`: 활성화 상태별 필터링
- `idx_template_attributes_template`: 템플릿의 속성 조회
- `idx_template_attributes_attribute`: 속성별 템플릿 검색

### 캐싱 전략
- 자주 사용되는 템플릿은 애플리케이션 레벨 캐싱
- 템플릿 데이터는 JSON으로 저장하여 파싱 오버헤드 최소화

## 향후 확장 계획

### 1. 템플릿 버전 관리
- 템플릿 버전 히스토리 추적
- 버전 간 마이그레이션 지원

### 2. 템플릿 상속
- 기본 템플릿에서 확장 가능한 하위 템플릿
- 속성 상속 및 오버라이드

### 3. 템플릿 검증 규칙
- 템플릿 데이터 구조 검증을 위한 JSON Schema
- 커스텀 검증 규칙 정의

### 4. 템플릿 인스턴스화
- 템플릿을 기반으로 노드 자동 생성
- 템플릿과 노드 간 연결 관계 추적