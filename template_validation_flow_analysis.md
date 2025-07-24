# 템플릿 검증 로직 흐름 분석

## 개요
MCP를 통해 태그 속성을 생성한 후 템플릿을 생성할 때, 템플릿 생성 과정에서 검증 로직이 어떻게 실행되는지 상세히 분석합니다.

## 검증 로직 실행 흐름

### 1. 전체 아키텍처 흐름

```
MCP Client Request
    ↓
MCPToolHandler.handleCreateTemplate()
    ↓
TemplateService.CreateTemplate()
    ↓
TemplateService.ValidateTemplateData()  ← 핵심 검증 지점
    ↓
TemplateValidator.ValidateTemplate()
    ↓
JSON 파싱 → 구조 검증 → 필수 필드 검증 → 타입 검증
    ↓
ValidationResult 반환
    ↓
성공/실패에 따른 템플릿 생성 또는 에러 반환
```

### 2. 의존성 주입 구조

#### 2.1 ApplicationFactory에서 TemplateService 생성
```go
// internal/interface/setup/factory.go:108-112
templateService, err := service.NewTemplateService(templateRepo, domainRepo)
if err != nil {
    panic("Failed to create template service: " + err.Error())
}
```

#### 2.2 TemplateService 생성 시 검증기 주입
```go
// internal/domain/service/template_service.go:56-66
func NewTemplateService(templateRepo repository.TemplateRepository, domainRepo repository.DomainRepository) (TemplateService, error) {
    validator, err := validation.NewTemplateValidator()
    if err != nil {
        return nil, fmt.Errorf("failed to create template validator: %w", err)
    }

    return &templateService{
        templateRepo: templateRepo,
        domainRepo:   domainRepo,
        validator:    validator,  ← 검증기가 주입됨
    }, nil
}
```

### 3. 검증 로직 실행 상세 분석

#### 3.1 MCP 핸들러에서 서비스 호출
```go
// internal/interface/mcp/tools.go:1330-1335
req := &service.CreateTemplateRequest{
    Name:         name,
    DomainName:   domainName,
    TemplateData: templateData,  ← 검증 대상 데이터
    Title:        title,
    Description:  description,
}

template, err := h.dependencies.TemplateService.CreateTemplate(ctx, req)
```

#### 3.2 TemplateService.CreateTemplate에서 검증 호출
```go
// internal/domain/service/template_service.go:95-105
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
    
    // 3. 이후 로직 (도메인 확인, 중복 확인, 엔티티 생성, 저장)
}
```

#### 3.3 ValidateTemplateData에서 검증기 호출
```go
// internal/domain/service/template_service.go:370-372
func (s *templateService) ValidateTemplateData(templateData string) (*validation.ValidationResult, error) {
    return s.validator.ValidateTemplate(templateData)  ← 실제 검증기 호출
}
```

#### 3.4 TemplateValidator에서 실제 검증 수행
```go
// internal/infrastructure/validation/template_validator.go:22-85
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
    var errors []ValidationError

    // version 필드 검증
    if version, exists := dataMap["version"]; !exists {
        errors = append(errors, ValidationError{
            Path:    "$.version",
            Message: "Version field is required",
        })
    } else if versionStr, ok := version.(string); !ok {
        errors = append(errors, ValidationError{
            Path:    "$.version",
            Message: "Version must be a string",
        })
    } else if !isValidSemanticVersion(versionStr) {
        errors = append(errors, ValidationError{
            Path:    "$.version",
            Message: "Invalid semantic version format",
            Value:   versionStr,
        })
    }

    // type 필드 검증
    if templateType, exists := dataMap["type"]; !exists {
        errors = append(errors, ValidationError{
            Path:    "$.type",
            Message: "Type field is required",
        })
    } else if typeStr, ok := templateType.(string); !ok {
        errors = append(errors, ValidationError{
            Path:    "$.type",
            Message: "Type must be a string",
        })
    } else if !isValidTemplateType(typeStr) {
        errors = append(errors, ValidationError{
            Path:    "$.type",
            Message: "Invalid template type",
            Value:   typeStr,
        })
    }

    // 4. 검증 결과 반환
    if len(errors) > 0 {
        return &ValidationResult{
            Valid:  false,
            Errors: errors,
        }, nil
    }

    return &ValidationResult{Valid: true}, nil
}
```

### 4. 검증 규칙 상세

#### 4.1 JSON 파싱 검증
- 템플릿 데이터가 유효한 JSON 형식인지 확인
- 파싱 실패 시 즉시 에러 반환

#### 4.2 구조 검증
- JSON 데이터가 객체(맵) 형태인지 확인
- 배열이나 기본 타입이 아닌 객체여야 함

#### 4.3 필수 필드 검증

##### 4.3.1 version 필드
- 필수 존재 여부 확인
- 문자열 타입 확인
- Semantic Version 형식 검증 (예: "1.0", "2.1.3")

##### 4.3.2 type 필드
- 필수 존재 여부 확인
- 문자열 타입 확인
- 유효한 템플릿 타입 확인 (layout, form, document, custom)

#### 4.4 Semantic Version 검증
```go
// internal/infrastructure/validation/template_validator.go:182-222
func isValidSemanticVersion(version string) bool {
    if version == "" {
        return false
    }

    // Basic semantic version validation: major.minor.patch[-prerelease]
    parts := strings.Split(version, ".")
    if len(parts) < 2 || len(parts) > 3 {
        return false
    }

    for i, part := range parts {
        if part == "" {
            return false
        }

        // For the last part, allow prerelease suffix
        if i == len(parts)-1 && strings.Contains(part, "-") {
            mainPart := strings.Split(part, "-")[0]
            if mainPart == "" {
                return false
            }
            // Check if main part is numeric
            for _, r := range mainPart {
                if r < '0' || r > '9' {
                    return false
                }
            }
        } else {
            // Check if part is numeric
            for _, r := range part {
                if r < '0' || r > '9' {
                    return false
                }
            }
        }
    }

    return true
}
```

#### 4.5 Template Type 검증
```go
// internal/infrastructure/validation/template_validator.go:223-231
func isValidTemplateType(templateType string) bool {
    validTypes := []string{"layout", "form", "document", "custom"}
    for _, validType := range validTypes {
        if templateType == validType {
            return true
        }
    }
    return false
}
```

### 5. 에러 처리 및 응답

#### 5.1 검증 실패 시 에러 구조
```go
type ValidationError struct {
    Message string                       `json:"message"`
    Errors  []validation.ValidationError `json:"errors"`
}

type ValidationError struct {
    Path    string      `json:"path"`
    Message string      `json:"message"`
    Value   interface{} `json:"value,omitempty"`
}
```

#### 5.2 MCP 응답 형식
```go
// 성공 시
return map[string]interface{}{
    "content": []map[string]interface{}{
        {
            "type": "text",
            "text": fmt.Sprintf("Template created successfully!\n\nComposite ID: url-db:%s:template:%d\nName: %s\nType: %s\nVersion: %s\nTitle: %s\nDescription: %s\nStatus: %s\nCreated: %s",
                domainName, template.ID(), template.Name(), templateType, templateVersion,
                template.Title(), template.Description(), getTemplateStatus(template.IsActive()),
                template.CreatedAt().Format("2006-01-02 15:04:05")),
        },
    },
    "isError": false,
}, nil

// 실패 시
return nil, fmt.Errorf("failed to create template: %w", err)
```

### 6. 검증 로직 실행 확인 방법

#### 6.1 로그 확인
- 템플릿 생성 시 검증 로직이 실행되는지 로그 확인
- 검증 실패 시 상세한 에러 메시지 확인

#### 6.2 테스트 시나리오
1. **유효한 템플릿 생성**: 검증 통과 후 템플릿 생성
2. **잘못된 JSON**: JSON 파싱 오류 발생
3. **필수 필드 누락**: version 또는 type 필드 누락 시 오류
4. **잘못된 타입**: 지원하지 않는 템플릿 타입 사용 시 오류

### 7. 결론

MCP를 통해 태그 속성을 생성한 후 템플릿을 생성할 때, **템플릿 생성 과정에서 반드시 검증 로직이 실행됩니다**:

1. **CreateTemplate 호출** → **ValidateTemplateData 호출** → **TemplateValidator.ValidateTemplate 실행**
2. 검증 로직은 JSON 파싱, 구조 검증, 필수 필드 검증, 타입 검증을 순차적으로 수행
3. 검증 실패 시 템플릿 생성이 중단되고 구조화된 에러 메시지 반환
4. 검증 성공 시에만 템플릿이 실제로 생성됨

이를 통해 데이터 무결성과 템플릿 품질을 보장할 수 있습니다. 