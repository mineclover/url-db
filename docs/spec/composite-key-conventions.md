# 합성키 컨벤션

## 개요
외부 시스템(MCP 클라이언트 등)과 데이터를 주고받을 때, 내부 ID를 숨기고 의미 있는 식별자를 제공하기 위해 합성키를 사용합니다.

## 합성키 구조
합성키는 다음과 같은 형태로 구성됩니다:
```
{tool_name}:{domain_name}:{id}
```

### 구성 요소
- **tool_name**: 도구/서비스 이름 (예: `url-db`, `bookmark-manager`)
- **domain_name**: 도메인 이름 (예: `tech-articles`, `recipes`) 
- **id**: 내부 노드 ID (예: `123`, `456`)

## 형식 규칙

### 1. 구분자
- 콜론(`:`)을 사용하여 각 구성 요소를 구분
- 예: `url-db:tech-articles:123`

### 2. 문자 제한
- 각 구성 요소는 영문자, 숫자, 하이픈(`-`), 언더스코어(`_`)만 허용
- 공백, 특수문자는 하이픈으로 변환
- 예: `"My Domain"` → `"my-domain"`

### 3. 대소문자
- 모든 구성 요소는 소문자로 변환
- 예: `"TechArticles"` → `"tech-articles"`

### 4. 길이 제한
- tool_name: 최대 50자
- domain_name: 최대 50자  
- id: 최대 20자

## 예시

### 기본 예시
```
url-db:tech-articles:123
url-db:recipes:456
url-db:personal-bookmarks:789
```

### 변환 예시
```
Original: "URL Database", "Tech Articles", 123
Composite: url-database:tech-articles:123

Original: "BookmarkManager", "My Personal Links", 456  
Composite: bookmarkmanager:my-personal-links:456
```

## 파싱 규칙

### 1. 분해 (Decomposition)
합성키를 개별 구성 요소로 분해:
```
"url-db:tech-articles:123" → ["url-db", "tech-articles", "123"]
```

### 2. 검증
- 정확히 3개의 구성 요소가 있는지 확인
- 각 구성 요소가 형식 규칙을 만족하는지 검증
- ID가 유효한 정수인지 확인

### 3. 내부 ID 변환
- 마지막 구성 요소를 정수로 변환하여 내부 ID로 사용
- 예: `"123"` → `123`

## 사용 사례

### 1. 노드 (Node) 합성키
```json
{
  "composite_id": "url-db:tech-articles:123",
  "url": "https://example.com/article",
  "title": "Example Article"
}
```

### 2. 템플릿 (Template) 합성키
```json
{
  "composite_id": "url-db:tech-articles:template:456",
  "name": "article-template",
  "template_data": "{...}"
}
```

템플릿 합성키는 추가 구분자를 포함합니다:
```
{tool_name}:{domain_name}:template:{id}
```

### 3. 외부 API 호출
```
GET /api/mcp/nodes/url-db:tech-articles:123
GET /api/mcp/templates/url-db:tech-articles:template:456
```

### 4. 배치 처리
```json
{
  "composite_ids": [
    "url-db:tech-articles:123",
    "url-db:recipes:456",
    "url-db:personal-bookmarks:789"
  ]
}
```

## 오류 처리

### 1. 잘못된 형식
```json
{
  "error": "INVALID_COMPOSITE_KEY",
  "message": "합성키 형식이 올바르지 않습니다",
  "expected_format": "tool_name:domain_name:id"
}
```

### 2. 존재하지 않는 리소스
```json
{
  "error": "RESOURCE_NOT_FOUND", 
  "message": "지정된 합성키에 해당하는 리소스를 찾을 수 없습니다",
  "composite_id": "url-db:tech-articles:999"
}
```

### 3. 권한 부족
```json
{
  "error": "ACCESS_DENIED",
  "message": "해당 도메인에 대한 접근 권한이 없습니다",
  "domain": "tech-articles"
}
```

## 구현 가이드

### 1. 합성키 생성
```javascript
function createCompositeKey(toolName, domainName, id) {
    const cleanTool = toolName.toLowerCase().replace(/[^a-z0-9\-_]/g, '-');
    const cleanDomain = domainName.toLowerCase().replace(/[^a-z0-9\-_]/g, '-');
    return `${cleanTool}:${cleanDomain}:${id}`;
}
```

### 2. 합성키 파싱
```javascript
function parseCompositeKey(compositeKey) {
    const parts = compositeKey.split(':');
    
    // 노드 합성키 (3개 부분)
    if (parts.length === 3) {
        const [toolName, domainName, idStr] = parts;
        const id = parseInt(idStr);
        
        if (isNaN(id)) {
            throw new Error('Invalid ID in composite key');
        }
        
        return { type: 'node', toolName, domainName, id };
    }
    
    // 템플릿 합성키 (4개 부분)
    if (parts.length === 4 && parts[2] === 'template') {
        const [toolName, domainName, _, idStr] = parts;
        const id = parseInt(idStr);
        
        if (isNaN(id)) {
            throw new Error('Invalid ID in composite key');
        }
        
        return { type: 'template', toolName, domainName, id };
    }
    
    throw new Error('Invalid composite key format');
}
```

## 마이그레이션 계획

### 1. 단계별 도입
- Phase 1: 내부 구현 및 테스트
- Phase 2: MCP 서버 API에 적용
- Phase 3: 기존 API에 추가 지원
- Phase 4: 기존 API 완전 대체

### 2. 호환성 유지
- 기존 내부 ID 기반 API는 당분간 유지
- 새로운 합성키 기반 API와 병행 지원
- 클라이언트가 점진적으로 마이그레이션할 수 있도록 지원