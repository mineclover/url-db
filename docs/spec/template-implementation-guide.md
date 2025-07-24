# 템플릿 시스템 구현 가이드 (jsonschema 기반)

## 개요
이 문서는 [santhosh-tekiri/jsonschema](https://github.com/santhosh-tekuri/jsonschema) 라이브러리를 사용하여 URL-DB 템플릿 시스템을 구현하는 방법을 설명합니다.

## 라이브러리 특징

### 주요 기능
- JSON Schema draft 4, 6, 7, 2019-09, 2020-12 지원
- 커스텀 정규식 엔진 지원
- 고급 format 검증 (UUID, IP, 날짜 등)
- Content assertions 및 encoding 검증
- 상세한 에러 분석
- 커스텀 vocabulary 지원
- 무한 검증 루프 감지

### 설치
```bash
go get github.com/santhosh-tekuri/jsonschema/v6
```

## 구현 아키텍처

### 1. 디렉토리 구조
```
internal/
├── domain/
│   └── entity/
│       └── template.go          # Template 엔티티
├── infrastructure/
│   └── validation/
│       ├── schema_loader.go     # 스키마 로더
│       ├── template_validator.go # 템플릿 검증기
│       └── schemas/             # JSON Schema 파일들
│           ├── base.json
│           ├── layout.json
│           ├── form.json
│           ├── document.json
│           └── custom.json
└── application/
    └── service/
        └── template_service.go   # 템플릿 서비스
```

### 2. 스키마 로더 구현

```go
// internal/infrastructure/validation/schema_loader.go
package validation

import (
    "embed"
    "fmt"
    "io/fs"
    "path/filepath"
    
    "github.com/santhosh-tekuri/jsonschema/v6"
)

//go:embed schemas/*.json
var schemaFS embed.FS

type SchemaLoader struct {
    compiler *jsonschema.Compiler
    schemas  map[string]*jsonschema.Schema
}

func NewSchemaLoader() (*SchemaLoader, error) {
    compiler := jsonschema.NewCompiler()
    
    // 커스텀 format 추가
    compiler.RegisterExtension("custom-formats", CustomFormats)
    
    loader := &SchemaLoader{
        compiler: compiler,
        schemas:  make(map[string]*jsonschema.Schema),
    }
    
    if err := loader.loadSchemas(); err != nil {
        return nil, fmt.Errorf("failed to load schemas: %w", err)
    }
    
    return loader, nil
}

func (sl *SchemaLoader) loadSchemas() error {
    return fs.WalkDir(schemaFS, "schemas", func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        
        if d.IsDir() || filepath.Ext(path) != ".json" {
            return nil
        }
        
        schemaData, err := schemaFS.ReadFile(path)
        if err != nil {
            return fmt.Errorf("failed to read schema file %s: %w", path, err)
        }
        
        schemaName := filepath.Base(path[:len(path)-5]) // .json 제거
        schemaURL := fmt.Sprintf("https://url-db.internal/schemas/%s.json", schemaName)
        
        if err := sl.compiler.AddResource(schemaURL, string(schemaData)); err != nil {
            return fmt.Errorf("failed to add schema resource %s: %w", schemaName, err)
        }
        
        schema, err := sl.compiler.Compile(schemaURL)
        if err != nil {
            return fmt.Errorf("failed to compile schema %s: %w", schemaName, err)
        }
        
        sl.schemas[schemaName] = schema
        return nil
    })
}

func (sl *SchemaLoader) GetSchema(templateType string) (*jsonschema.Schema, error) {
    schema, exists := sl.schemas[templateType]
    if !exists {
        return nil, fmt.Errorf("schema not found for template type: %s", templateType)
    }
    return schema, nil
}

// 커스텀 format 정의
var CustomFormats = map[string]jsonschema.Format{
    "semantic-version": {
        Validate: func(value string) error {
            // semantic version 검증 로직
            if !isValidSemanticVersion(value) {
                return fmt.Errorf("invalid semantic version: %s", value)
            }
            return nil
        },
    },
    "template-name": {
        Validate: func(value string) error {
            // 템플릿 이름 검증 로직
            if !isValidTemplateName(value) {
                return fmt.Errorf("invalid template name: %s", value)
            }
            return nil
        },
    },
}

func isValidSemanticVersion(version string) bool {
    // semantic version 정규식 검증
    // 예: "1.0.0", "2.1.3-beta"
    return true // 실제 구현 필요
}

func isValidTemplateName(name string) bool {
    // 템플릿 이름 검증 (영문자, 숫자, 하이픈, 언더스코어만 허용)
    return true // 실제 구현 필요
}
```

### 3. 템플릿 검증기 구현

```go
// internal/infrastructure/validation/template_validator.go
package validation

import (
    "encoding/json"
    "fmt"
    
    "github.com/santhosh-tekuri/jsonschema/v6"
)

type TemplateValidator struct {
    schemaLoader *SchemaLoader
}

func NewTemplateValidator() (*TemplateValidator, error) {
    loader, err := NewSchemaLoader()
    if err != nil {
        return nil, err
    }
    
    return &TemplateValidator{
        schemaLoader: loader,
    }, nil
}

type ValidationResult struct {
    Valid  bool                    `json:"valid"`
    Errors []ValidationError       `json:"errors,omitempty"`
}

type ValidationError struct {
    Path    string `json:"path"`
    Message string `json:"message"`
    Value   any    `json:"value,omitempty"`
}

func (tv *TemplateValidator) ValidateTemplateData(templateType string, data []byte) (*ValidationResult, error) {
    // JSON 파싱 검증
    var jsonData interface{}
    if err := json.Unmarshal(data, &jsonData); err != nil {
        return &ValidationResult{
            Valid: false,
            Errors: []ValidationError{
                {
                    Path:    "$",
                    Message: fmt.Sprintf("Invalid JSON: %s", err.Error()),
                },
            },
        }, nil
    }
    
    // 기본 템플릿 스키마 검증
    baseSchema, err := tv.schemaLoader.GetSchema("base")
    if err != nil {
        return nil, fmt.Errorf("failed to get base schema: %w", err)
    }
    
    if err := baseSchema.Validate(jsonData); err != nil {
        return tv.convertValidationError(err), nil
    }
    
    // 타입별 스키마 검증
    if templateType != "" {
        typeSchema, err := tv.schemaLoader.GetSchema(templateType)
        if err != nil {
            return nil, fmt.Errorf("failed to get schema for type %s: %w", templateType, err)
        }
        
        if err := typeSchema.Validate(jsonData); err != nil {
            return tv.convertValidationError(err), nil
        }
    }
    
    return &ValidationResult{Valid: true}, nil
}

func (tv *TemplateValidator) convertValidationError(err error) *ValidationResult {
    var errors []ValidationError
    
    if ve, ok := err.(*jsonschema.ValidationError); ok {
        errors = tv.extractValidationErrors(ve)
    } else {
        errors = []ValidationError{
            {
                Path:    "$",
                Message: err.Error(),
            },
        }
    }
    
    return &ValidationResult{
        Valid:  false,
        Errors: errors,
    }
}

func (tv *TemplateValidator) extractValidationErrors(ve *jsonschema.ValidationError) []ValidationError {
    var errors []ValidationError
    
    // 메인 에러 추가
    errors = append(errors, ValidationError{
        Path:    ve.InstanceLocation,
        Message: ve.Message,
        Value:   ve.InstanceValue,
    })
    
    // 하위 에러들 재귀적으로 추가
    for _, cause := range ve.Causes {
        if subVe, ok := cause.(*jsonschema.ValidationError); ok {
            errors = append(errors, tv.extractValidationErrors(subVe)...)
        }
    }
    
    return errors
}

// 특정 템플릿 타입에 맞는 기본 구조 생성
func (tv *TemplateValidator) GenerateTemplate(templateType string) (map[string]interface{}, error) {
    templates := map[string]map[string]interface{}{
        "layout": {
            "version": "1.0",
            "type":    "layout",
            "metadata": map[string]interface{}{
                "name":        "",
                "description": "",
                "tags":        []string{},
            },
            "content": map[string]interface{}{
                "structure": map[string]interface{}{
                    "type":  "grid",
                    "areas": []map[string]interface{}{},
                },
            },
            "presentation": map[string]interface{}{
                "theme":      "light",
                "responsive": true,
            },
        },
        "form": {
            "version": "1.0",
            "type":    "form",
            "metadata": map[string]interface{}{
                "name":        "",
                "description": "",
            },
            "schema": map[string]interface{}{
                "fields":   []map[string]interface{}{},
                "sections": []map[string]interface{}{},
            },
            "validation": map[string]interface{}{
                "rules": map[string]interface{}{},
            },
            "presentation": map[string]interface{}{
                "layout": "vertical",
            },
        },
        "document": {
            "version": "1.0",
            "type":    "document",
            "metadata": map[string]interface{}{
                "name": "",
            },
            "schema": map[string]interface{}{
                "sections": []map[string]interface{}{},
            },
        },
        "custom": {
            "version": "1.0",
            "type":    "custom",
            "metadata": map[string]interface{}{
                "name": "",
            },
            "content": map[string]interface{}{},
        },
    }
    
    template, exists := templates[templateType]
    if !exists {
        return nil, fmt.Errorf("unknown template type: %s", templateType)
    }
    
    return template, nil
}
```

### 4. JSON Schema 파일들

#### Base Schema
```json
// internal/infrastructure/validation/schemas/base.json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://url-db.internal/schemas/base.json",
  "title": "Base Template Schema",
  "type": "object",
  "required": ["version", "type"],
  "properties": {
    "version": {
      "type": "string",
      "format": "semantic-version",
      "description": "Template version using semantic versioning"
    },
    "type": {
      "type": "string",
      "enum": ["layout", "form", "document", "custom"],
      "description": "Template type"
    },
    "metadata": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string",
          "format": "template-name",
          "maxLength": 255
        },
        "description": {
          "type": "string",
          "maxLength": 1000
        },
        "author": {
          "type": "string",
          "format": "email"
        },
        "created": {
          "type": "string",
          "format": "date-time"
        },
        "modified": {
          "type": "string",
          "format": "date-time"
        },
        "tags": {
          "type": "array",
          "items": {
            "type": "string",
            "minLength": 1
          },
          "uniqueItems": true
        },
        "category": {
          "type": "string"
        },
        "status": {
          "type": "string",
          "enum": ["draft", "published", "archived"],
          "default": "draft"
        }
      }
    },
    "validation": {
      "$ref": "#/$defs/validation"
    },
    "presentation": {
      "$ref": "#/$defs/presentation"
    }
  },
  "$defs": {
    "validation": {
      "type": "object",
      "properties": {
        "rules": {
          "type": "object",
          "additionalProperties": {
            "type": "object",
            "properties": {
              "type": {
                "type": "string"
              },
              "required": {
                "type": "boolean"
              },
              "min": {
                "type": "number"
              },
              "max": {
                "type": "number"
              },
              "pattern": {
                "type": "string"
              },
              "enum": {
                "type": "array"
              }
            }
          }
        },
        "messages": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        }
      }
    },
    "presentation": {
      "type": "object",
      "properties": {
        "theme": {
          "type": "string",
          "enum": ["light", "dark", "auto", "custom"]
        },
        "layout": {
          "type": "string"
        },
        "responsive": {
          "type": "boolean",
          "default": true
        },
        "css": {
          "type": "object",
          "properties": {
            "classes": {
              "type": "array",
              "items": {
                "type": "string"
              }
            },
            "inline": {
              "type": "string"
            }
          }
        }
      }
    }
  }
}
```

#### Layout Schema
```json
// internal/infrastructure/validation/schemas/layout.json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://url-db.internal/schemas/layout.json",
  "title": "Layout Template Schema",
  "allOf": [
    {
      "$ref": "base.json"
    },
    {
      "if": {
        "properties": {
          "type": { "const": "layout" }
        }
      },
      "then": {
        "required": ["content"],
        "properties": {
          "content": {
            "type": "object",
            "required": ["structure"],
            "properties": {
              "structure": {
                "type": "object",
                "required": ["type"],
                "properties": {
                  "type": {
                    "type": "string",
                    "enum": ["grid", "flex", "card", "list"]
                  },
                  "container": {
                    "type": "object",
                    "properties": {
                      "maxWidth": {
                        "type": "string"
                      },
                      "padding": {
                        "type": "string"
                      },
                      "margin": {
                        "type": "string"
                      }
                    }
                  },
                  "areas": {
                    "type": "array",
                    "items": {
                      "type": "object",
                      "required": ["name"],
                      "properties": {
                        "name": {
                          "type": "string",
                          "minLength": 1
                        },
                        "width": {
                          "type": "string"
                        },
                        "height": {
                          "type": "string"
                        },
                        "position": {
                          "type": "string",
                          "enum": ["left", "right", "top", "bottom", "center"]
                        },
                        "gridArea": {
                          "type": "string"
                        }
                      }
                    }
                  }
                }
              },
              "components": {
                "type": "object",
                "additionalProperties": {
                  "type": "object"
                }
              }
            }
          }
        }
      }
    }
  ]
}
```

### 5. 템플릿 서비스 구현

```go
// internal/application/service/template_service.go
package service

import (
    "encoding/json"
    "fmt"
    
    "url-db/internal/domain/entity"
    "url-db/internal/domain/repository"
    "url-db/internal/infrastructure/validation"
)

type TemplateService struct {
    templateRepo repository.TemplateRepository
    validator    *validation.TemplateValidator
}

func NewTemplateService(templateRepo repository.TemplateRepository) (*TemplateService, error) {
    validator, err := validation.NewTemplateValidator()
    if err != nil {
        return nil, fmt.Errorf("failed to create template validator: %w", err)
    }
    
    return &TemplateService{
        templateRepo: templateRepo,
        validator:    validator,
    }, nil
}

func (ts *TemplateService) CreateTemplate(req *CreateTemplateRequest) (*entity.Template, error) {
    // 템플릿 데이터 검증
    result, err := ts.validator.ValidateTemplateData(req.Type, []byte(req.TemplateData))
    if err != nil {
        return nil, fmt.Errorf("validation error: %w", err)
    }
    
    if !result.Valid {
        return nil, &ValidationError{
            Message: "Template data validation failed",
            Errors:  result.Errors,
        }
    }
    
    // 템플릿 엔티티 생성
    template := &entity.Template{
        Name:         req.Name,
        DomainID:     req.DomainID,
        TemplateData: req.TemplateData,
        Title:        req.Title,
        Description:  req.Description,
        IsActive:     true,
    }
    
    // 저장
    if err := ts.templateRepo.Create(template); err != nil {
        return nil, fmt.Errorf("failed to create template: %w", err)
    }
    
    return template, nil
}

func (ts *TemplateService) UpdateTemplate(id int, req *UpdateTemplateRequest) (*entity.Template, error) {
    // 기존 템플릿 조회
    template, err := ts.templateRepo.GetByID(id)
    if err != nil {
        return nil, fmt.Errorf("template not found: %w", err)
    }
    
    // 새 템플릿 데이터 검증 (제공된 경우)
    if req.TemplateData != "" {
        var templateType string
        var templateData map[string]interface{}
        if err := json.Unmarshal([]byte(req.TemplateData), &templateData); err == nil {
            if t, ok := templateData["type"].(string); ok {
                templateType = t
            }
        }
        
        result, err := ts.validator.ValidateTemplateData(templateType, []byte(req.TemplateData))
        if err != nil {
            return nil, fmt.Errorf("validation error: %w", err)
        }
        
        if !result.Valid {
            return nil, &ValidationError{
                Message: "Template data validation failed",
                Errors:  result.Errors,
            }
        }
        
        template.TemplateData = req.TemplateData
    }
    
    // 다른 필드 업데이트
    if req.Title != "" {
        template.Title = req.Title
    }
    if req.Description != "" {
        template.Description = req.Description
    }
    if req.IsActive != nil {
        template.IsActive = *req.IsActive
    }
    
    // 저장
    if err := ts.templateRepo.Update(template); err != nil {
        return nil, fmt.Errorf("failed to update template: %w", err)
    }
    
    return template, nil
}

func (ts *TemplateService) GenerateTemplateScaffold(templateType string) (string, error) {
    template, err := ts.validator.GenerateTemplate(templateType)
    if err != nil {
        return "", fmt.Errorf("failed to generate template: %w", err)
    }
    
    data, err := json.MarshalIndent(template, "", "  ")
    if err != nil {
        return "", fmt.Errorf("failed to marshal template: %w", err)
    }
    
    return string(data), nil
}

type CreateTemplateRequest struct {
    Name         string `json:"name"`
    DomainID     int    `json:"domain_id"`
    Type         string `json:"type"`
    TemplateData string `json:"template_data"`
    Title        string `json:"title"`
    Description  string `json:"description"`
}

type UpdateTemplateRequest struct {
    TemplateData string `json:"template_data,omitempty"`
    Title        string `json:"title,omitempty"`
    Description  string `json:"description,omitempty"`
    IsActive     *bool  `json:"is_active,omitempty"`
}

type ValidationError struct {
    Message string                         `json:"message"`
    Errors  []validation.ValidationError   `json:"errors"`
}

func (e *ValidationError) Error() string {
    return e.Message
}
```

### 6. 사용 예제

#### 템플릿 생성 및 검증
```go
func ExampleTemplateCreation() {
    // 서비스 초기화
    templateService, err := service.NewTemplateService(templateRepo)
    if err != nil {
        log.Fatal(err)
    }
    
    // 레이아웃 템플릿 생성
    layoutData := `{
        "version": "1.0",
        "type": "layout",
        "metadata": {
            "name": "Two Column Layout",
            "description": "Responsive two-column layout"
        },
        "content": {
            "structure": {
                "type": "grid",
                "areas": [
                    {
                        "name": "sidebar",
                        "width": "300px",
                        "position": "left"
                    },
                    {
                        "name": "main",
                        "width": "auto",
                        "position": "right"
                    }
                ]
            }
        },
        "presentation": {
            "theme": "light",
            "responsive": true
        }
    }`
    
    req := &service.CreateTemplateRequest{
        Name:         "two-column-layout",
        DomainID:     1,
        Type:         "layout",
        TemplateData: layoutData,
        Title:        "Two Column Layout",
        Description:  "Standard two-column responsive layout",
    }
    
    template, err := templateService.CreateTemplate(req)
    if err != nil {
        if validationErr, ok := err.(*service.ValidationError); ok {
            for _, e := range validationErr.Errors {
                fmt.Printf("Validation error at %s: %s\n", e.Path, e.Message)
            }
        } else {
            log.Fatal(err)
        }
        return
    }
    
    fmt.Printf("Template created with ID: %d\n", template.ID)
}
```

#### CLI 도구를 사용한 스키마 검증
```bash
# jsonschema CLI 도구 설치
go install github.com/santhosh-tekuri/jsonschema/cmd/jv@latest

# 템플릿 데이터 검증
jv -s internal/infrastructure/validation/schemas/layout.json template-data.json

# 상세 오류 출력
jv -s internal/infrastructure/validation/schemas/layout.json -o detailed template-data.json
```

## 모범 사례

### 1. 스키마 버전 관리
- 스키마 파일에 버전 정보 포함
- 하위 호환성 유지를 위한 gradual migration
- 스키마 변경 시 deprecation 경고

### 2. 에러 처리
- 상세한 검증 오류 메시지 제공
- 사용자 친화적인 오류 표시
- 검증 실패 시 수정 가이드 제공

### 3. 성능 최적화
- 스키마 컴파일 결과 캐싱
- 대용량 템플릿 데이터 스트리밍 검증
- 검증 결과 캐싱

### 4. 확장성
- 커스텀 format 및 vocabulary 활용
- 플러그인 시스템과 연동
- 동적 스키마 로딩

이 가이드를 통해 jsonschema 라이브러리를 활용한 강력하고 유연한 템플릿 시스템을 구축할 수 있습니다.