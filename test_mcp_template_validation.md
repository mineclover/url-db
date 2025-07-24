# MCP 태그 속성 생성 및 템플릿 검증 테스트

## 테스트 목적
MCP를 통해 태그 속성을 생성한 후 템플릿을 생성하고, 템플릿 생성 과정에서 검증 로직이 실행되는지 확인합니다.

## 검증 로직 흐름 분석

### 1. 템플릿 생성 시 검증 로직 실행 흐름

```
MCP handleCreateTemplate
    ↓
TemplateService.CreateTemplate()
    ↓
TemplateService.ValidateTemplateData()  ← 검증 로직 실행 지점
    ↓
TemplateValidator.ValidateTemplate()
    ↓
JSON 파싱 검증 → 구조 검증 → 필수 필드 검증 → 타입 검증
```

### 2. 검증 로직 상세 분석

#### 2.1 TemplateService.CreateTemplate() 메서드
```go
func (s *templateService) CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*entity.Template, error) {
    // 1. 템플릿 이름 검증
    if err := s.ValidateTemplateName(req.Name); err != nil {
        return nil, fmt.Errorf("invalid template name: %w", err)
    }

    // 2. 템플릿 데이터 검증 ← 핵심 검증 로직
    result, err := s.ValidateTemplateData(req.TemplateData)
    if err != nil {
        return nil, fmt.Errorf("template validation error: %w", err)
    }

    if !result.Valid {
        return nil, &ValidationError{
            Message: "Template data validation failed",
            Errors:  result.Errors,
        }
    }
    
    // 3. 도메인 존재 확인
    // 4. 템플릿 이름 중복 확인
    // 5. 템플릿 엔티티 생성
    // 6. 저장소에 저장
}
```

#### 2.2 TemplateValidator.ValidateTemplate() 메서드
```go
func (tv *TemplateValidator) ValidateTemplate(templateData string) (*ValidationResult, error) {
    // 1. JSON 파싱 검증
    var data interface{}
    if err := json.Unmarshal([]byte(templateData), &data); err != nil {
        return &ValidationResult{
            Valid: false,
            Errors: []ValidationError{{
                Path:    "$",
                Message: fmt.Sprintf("Invalid JSON: %s", err.Error()),
            }},
        }, nil
    }

    // 2. 기본 구조 검증 (객체인지 확인)
    dataMap, ok := data.(map[string]interface{})
    if !ok {
        return &ValidationResult{
            Valid: false,
            Errors: []ValidationError{{
                Path:    "$",
                Message: "Template data must be an object",
            }},
        }, nil
    }

    // 3. 필수 필드 검증
    // - version 필드 검증 (semantic version 형식)
    // - type 필드 검증 (layout, form, document, custom 중 하나)
    
    // 4. 검증 결과 반환
    return &ValidationResult{Valid: true}, nil
}
```

## 테스트 시나리오

### 시나리오 1: 성공적인 템플릿 생성 (검증 통과)

#### 1.1 도메인 생성
```json
{
  "name": "create_domain",
  "arguments": {
    "name": "test-template-domain",
    "description": "템플릿 검증 테스트용 도메인"
  }
}
```

#### 1.2 태그 속성 생성
```json
{
  "name": "create_domain_attribute",
  "arguments": {
    "domain_name": "test-template-domain",
    "name": "category",
    "type": "tag",
    "description": "템플릿 카테고리 분류"
  }
}
```

#### 1.3 유효한 템플릿 생성 (검증 통과)
```json
{
  "name": "create_template",
  "arguments": {
    "domain_name": "test-template-domain",
    "name": "test-layout-template",
    "template_data": "{\"version\":\"1.0\",\"type\":\"layout\",\"metadata\":{\"name\":\"Test Layout\",\"description\":\"테스트용 레이아웃 템플릿\"}}",
    "title": "테스트 레이아웃 템플릿",
    "description": "템플릿 검증 테스트용 레이아웃"
  }
}
```

**예상 결과:**
- ✅ 템플릿 생성 성공
- ✅ 검증 로직이 실행되어 유효한 JSON과 필수 필드 확인
- ✅ Composite ID 반환: `url-db:test-template-domain:template:1`

### 시나리오 2: 실패하는 템플릿 생성 (검증 실패)

#### 2.1 잘못된 JSON 형식
```json
{
  "name": "create_template",
  "arguments": {
    "domain_name": "test-template-domain",
    "name": "invalid-json-template",
    "template_data": "{\"version\":\"1.0\",\"type\":\"layout\",\"metadata\":{\"name\":\"Test\",\"description\":\"Test\"}",
    "title": "잘못된 JSON 템플릿",
    "description": "JSON 형식 오류 테스트"
  }
}
```

**예상 결과:**
- ❌ 템플릿 생성 실패
- ❌ 검증 로직이 실행되어 JSON 파싱 오류 감지
- ❌ 에러 메시지: "Invalid JSON: unexpected end of JSON input"

#### 2.2 필수 필드 누락
```json
{
  "name": "create_template",
  "arguments": {
    "domain_name": "test-template-domain",
    "name": "missing-fields-template",
    "template_data": "{\"metadata\":{\"name\":\"Test\"}}",
    "title": "필수 필드 누락 템플릿",
    "description": "필수 필드 누락 테스트"
  }
}
```

**예상 결과:**
- ❌ 템플릿 생성 실패
- ❌ 검증 로직이 실행되어 필수 필드 누락 감지
- ❌ 에러 메시지: "Version field is required", "Type field is required"

#### 2.3 잘못된 타입
```json
{
  "name": "create_template",
  "arguments": {
    "domain_name": "test-template-domain",
    "name": "invalid-type-template",
    "template_data": "{\"version\":\"1.0\",\"type\":\"invalid_type\",\"metadata\":{\"name\":\"Test\"}}",
    "title": "잘못된 타입 템플릿",
    "description": "잘못된 타입 테스트"
  }
}
```

**예상 결과:**
- ❌ 템플릿 생성 실패
- ❌ 검증 로직이 실행되어 잘못된 타입 감지
- ❌ 에러 메시지: "Invalid template type"

### 시나리오 3: 템플릿 검증 도구 테스트

#### 3.1 유효한 템플릿 검증
```json
{
  "name": "validate_template",
  "arguments": {
    "template_data": "{\"version\":\"1.0\",\"type\":\"form\",\"metadata\":{\"name\":\"Test Form\",\"description\":\"테스트 폼\"}}"
  }
}
```

**예상 결과:**
- ✅ 검증 성공
- ✅ 응답: "✅ Template validation successful!\n\nType: form\nVersion: 1.0\n\nThe template data is valid and can be used to create a new template."

#### 3.2 잘못된 템플릿 검증
```json
{
  "name": "validate_template",
  "arguments": {
    "template_data": "{\"version\":\"1.0\",\"type\":\"invalid_type\"}"
  }
}
```

**예상 결과:**
- ❌ 검증 실패
- ❌ 응답: "❌ Template validation failed!\n\nErrors:\n1. Path: $.type - Invalid template type (value: invalid_type)"

## 검증 로직 실행 확인 포인트

### 1. CreateTemplate 메서드에서 검증 호출 확인
```go
// internal/domain/service/template_service.go:95-105
result, err := s.ValidateTemplateData(req.TemplateData)
if err != nil {
    return nil, fmt.Errorf("template validation error: %w", err)
}

if !result.Valid {
    return nil, &ValidationError{
        Message: "Template data validation failed",
        Errors:  result.Errors,
    }
}
```

### 2. ValidateTemplateData 메서드에서 검증기 호출 확인
```go
// internal/domain/service/template_service.go:370-372
func (s *templateService) ValidateTemplateData(templateData string) (*validation.ValidationResult, error) {
    return s.validator.ValidateTemplate(templateData)
}
```

### 3. TemplateValidator에서 실제 검증 로직 실행 확인
```go
// internal/infrastructure/validation/template_validator.go:22-85
func (tv *TemplateValidator) ValidateTemplate(templateData string) (*ValidationResult, error) {
    // JSON 파싱 검증
    // 구조 검증
    // 필수 필드 검증
    // 타입 검증
}
```

## 결론

MCP를 통해 태그 속성을 생성한 후 템플릿을 생성할 때, **템플릿 생성 과정에서 반드시 검증 로직이 실행됩니다**:

1. **CreateTemplate 호출** → **ValidateTemplateData 호출** → **TemplateValidator.ValidateTemplate 실행**
2. 검증 로직은 JSON 파싱, 구조 검증, 필수 필드 검증, 타입 검증을 수행
3. 검증 실패 시 템플릿 생성이 중단되고 적절한 에러 메시지 반환
4. 검증 성공 시에만 템플릿이 실제로 생성됨

이를 통해 데이터 무결성과 템플릿 품질을 보장할 수 있습니다. 