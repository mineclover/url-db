# 템플릿 JSON 구조 명세

## 개요
템플릿의 `template_data` 필드는 JSON 형식으로 저장되며, 다양한 용도로 활용 가능한 구조화된 데이터를 포함합니다.

## 기본 JSON 구조

### 최상위 구조
```json
{
  "version": "1.0",
  "type": "layout|form|document|custom",
  "metadata": {
    "author": "string",
    "created": "ISO 8601 date",
    "modified": "ISO 8601 date",
    "tags": ["string"]
  },
  "schema": {},
  "content": {},
  "validation": {},
  "presentation": {}
}
```

## 템플릿 타입별 구조

### 1. Layout 템플릿
페이지 레이아웃이나 UI 구조를 정의하는 템플릿

```json
{
  "version": "1.0",
  "type": "layout",
  "metadata": {
    "name": "Two Column Layout",
    "description": "Responsive two-column layout with sidebar"
  },
  "content": {
    "structure": {
      "container": {
        "type": "grid",
        "columns": 2,
        "gap": "20px",
        "responsive": {
          "mobile": { "columns": 1 },
          "tablet": { "columns": 2 }
        }
      },
      "areas": [
        {
          "name": "sidebar",
          "width": "300px",
          "position": "left",
          "sticky": true
        },
        {
          "name": "main",
          "width": "auto",
          "position": "right"
        }
      ]
    },
    "components": {
      "sidebar": {
        "allowed": ["navigation", "widget", "advertisement"],
        "max_items": 5
      },
      "main": {
        "allowed": ["article", "gallery", "video", "form"],
        "required": ["article"]
      }
    }
  },
  "presentation": {
    "theme": "light",
    "css_classes": ["layout-two-column", "responsive"],
    "custom_css": "/* Optional custom CSS */"
  }
}
```

### 2. Form 템플릿
데이터 입력 폼을 정의하는 템플릿

```json
{
  "version": "1.0",
  "type": "form",
  "metadata": {
    "name": "Product Registration Form",
    "description": "Form for adding new products"
  },
  "schema": {
    "fields": [
      {
        "name": "title",
        "type": "text",
        "label": "Product Title",
        "required": true,
        "maxLength": 200,
        "placeholder": "Enter product name"
      },
      {
        "name": "price",
        "type": "number",
        "label": "Price",
        "required": true,
        "min": 0,
        "max": 999999.99,
        "step": 0.01,
        "currency": "USD"
      },
      {
        "name": "category",
        "type": "select",
        "label": "Category",
        "required": true,
        "options": [
          { "value": "electronics", "label": "Electronics" },
          { "value": "clothing", "label": "Clothing" },
          { "value": "food", "label": "Food & Beverages" }
        ],
        "multiple": false
      },
      {
        "name": "description",
        "type": "textarea",
        "label": "Description",
        "required": false,
        "maxLength": 5000,
        "rows": 10
      },
      {
        "name": "images",
        "type": "file",
        "label": "Product Images",
        "accept": "image/*",
        "multiple": true,
        "maxFiles": 10,
        "maxSize": "5MB"
      },
      {
        "name": "tags",
        "type": "tags",
        "label": "Tags",
        "maxTags": 20,
        "autocomplete": true
      }
    ],
    "sections": [
      {
        "title": "Basic Information",
        "fields": ["title", "price", "category"]
      },
      {
        "title": "Details",
        "fields": ["description", "images", "tags"]
      }
    ]
  },
  "validation": {
    "rules": [
      {
        "field": "title",
        "rules": ["required", "minLength:3", "maxLength:200"]
      },
      {
        "field": "price",
        "rules": ["required", "numeric", "min:0"]
      }
    ],
    "custom": [
      {
        "name": "unique_title",
        "message": "Product title must be unique",
        "check": "async_validation_endpoint"
      }
    ]
  },
  "presentation": {
    "layout": "vertical",
    "submit_button": {
      "text": "Add Product",
      "position": "bottom-right"
    },
    "show_progress": true,
    "auto_save": true
  }
}
```

### 3. Document 템플릿
구조화된 문서나 콘텐츠를 정의하는 템플릿

```json
{
  "version": "1.0",
  "type": "document",
  "metadata": {
    "name": "Blog Post Template",
    "description": "Standard blog post structure"
  },
  "schema": {
    "sections": [
      {
        "id": "header",
        "type": "header",
        "required": true,
        "fields": {
          "title": {
            "type": "text",
            "required": true,
            "maxLength": 150
          },
          "subtitle": {
            "type": "text",
            "required": false,
            "maxLength": 200
          },
          "author": {
            "type": "author",
            "required": true
          },
          "publish_date": {
            "type": "datetime",
            "required": true
          },
          "featured_image": {
            "type": "image",
            "required": false,
            "dimensions": {
              "min_width": 1200,
              "aspect_ratio": "16:9"
            }
          }
        }
      },
      {
        "id": "content",
        "type": "body",
        "required": true,
        "blocks": [
          {
            "type": "paragraph",
            "min": 1,
            "max": null
          },
          {
            "type": "heading",
            "levels": [2, 3, 4],
            "max": 10
          },
          {
            "type": "image",
            "max": 20,
            "caption": true
          },
          {
            "type": "code",
            "languages": ["javascript", "python", "go", "sql"],
            "syntax_highlight": true
          },
          {
            "type": "quote",
            "max": 5
          },
          {
            "type": "list",
            "ordered": true,
            "unordered": true
          }
        ]
      },
      {
        "id": "footer",
        "type": "footer",
        "required": false,
        "components": ["tags", "related_posts", "comments"]
      }
    ]
  },
  "content": {
    "default_blocks": [
      {
        "type": "paragraph",
        "content": "Start writing your blog post here..."
      }
    ],
    "templates": {
      "introduction": "In this post, we'll explore...",
      "conclusion": "To summarize..."
    }
  },
  "validation": {
    "min_words": 300,
    "max_words": 5000,
    "required_sections": ["header", "content"],
    "seo": {
      "title_length": { "min": 30, "max": 60 },
      "meta_description": { "min": 120, "max": 160 }
    }
  }
}
```

### 4. Custom 템플릿
특수 목적의 커스텀 템플릿

```json
{
  "version": "1.0",
  "type": "custom",
  "metadata": {
    "name": "API Endpoint Template",
    "description": "Template for API endpoint documentation"
  },
  "content": {
    "endpoint": {
      "method": "POST",
      "path": "/api/v1/resources",
      "authentication": "Bearer Token",
      "rate_limit": "100 requests/hour"
    },
    "request": {
      "headers": {
        "Content-Type": "application/json",
        "Authorization": "Bearer {token}"
      },
      "body": {
        "type": "object",
        "properties": {
          "name": {
            "type": "string",
            "required": true,
            "description": "Resource name"
          },
          "data": {
            "type": "object",
            "required": false,
            "description": "Additional data"
          }
        }
      },
      "examples": [
        {
          "title": "Create new resource",
          "body": {
            "name": "My Resource",
            "data": { "key": "value" }
          }
        }
      ]
    },
    "response": {
      "success": {
        "status": 201,
        "body": {
          "id": "string",
          "name": "string",
          "created_at": "datetime"
        }
      },
      "errors": [
        {
          "status": 400,
          "code": "VALIDATION_ERROR",
          "message": "Invalid request data"
        },
        {
          "status": 401,
          "code": "UNAUTHORIZED",
          "message": "Invalid or missing token"
        }
      ]
    }
  }
}
```

## 공통 컴포넌트

### Metadata 구조
```json
{
  "author": "user@example.com",
  "created": "2024-01-15T10:30:00Z",
  "modified": "2024-01-20T14:45:00Z",
  "version": "1.2",
  "tags": ["template", "reusable", "production"],
  "category": "layout",
  "status": "published",
  "permissions": {
    "view": ["all"],
    "edit": ["owner", "admin"],
    "delete": ["owner"]
  }
}
```

### Validation 규칙
```json
{
  "rules": {
    "field_name": {
      "type": "string|number|boolean|array|object",
      "required": true,
      "min": 1,
      "max": 100,
      "pattern": "^[A-Za-z0-9]+$",
      "enum": ["option1", "option2"],
      "custom": "validation_function_name"
    }
  },
  "messages": {
    "field_name.required": "This field is required",
    "field_name.min": "Minimum value is {min}"
  },
  "dependencies": {
    "field1": ["field2", "field3"]
  }
}
```

### Presentation 옵션
```json
{
  "theme": "light|dark|custom",
  "layout": "vertical|horizontal|grid",
  "responsive": true,
  "animations": {
    "enabled": true,
    "type": "fade|slide|zoom"
  },
  "css": {
    "classes": ["custom-class-1", "custom-class-2"],
    "inline": "color: #333; font-size: 16px;",
    "external": ["https://example.com/styles.css"]
  },
  "javascript": {
    "inline": "console.log('Template loaded');",
    "external": ["https://example.com/script.js"]
  }
}
```

## 사용 예시

### 템플릿 생성 요청
```json
{
  "name": "product-card",
  "domain_name": "ecommerce",
  "template_data": {
    "version": "1.0",
    "type": "layout",
    "content": {
      "structure": {
        "type": "card",
        "sections": ["image", "title", "price", "actions"]
      }
    }
  },
  "title": "Product Card Template",
  "description": "Reusable product card component"
}
```

### 템플릿 응답
```json
{
  "composite_id": "url-db:ecommerce:template:1",
  "name": "product-card",
  "domain_name": "ecommerce",
  "template_data": {
    "version": "1.0",
    "type": "layout",
    "content": {...}
  },
  "title": "Product Card Template",
  "description": "Reusable product card component",
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

## 검증 및 제약사항

1. **JSON 유효성**: template_data는 반드시 유효한 JSON이어야 함
2. **버전 관리**: version 필드는 필수이며 semantic versioning 권장
3. **타입 제한**: type은 정의된 값 중 하나여야 함
4. **크기 제한**: template_data의 최대 크기는 1MB
5. **필수 필드**: version과 type은 최소 필수 필드

## 확장성

### 플러그인 시스템
```json
{
  "plugins": [
    {
      "name": "markdown-renderer",
      "version": "2.0",
      "config": {
        "sanitize": true,
        "breaks": true
      }
    }
  ]
}
```

### 변수 및 표현식
```json
{
  "variables": {
    "site_name": "{{config.site.name}}",
    "current_date": "{{now | date:'YYYY-MM-DD'}}"
  },
  "expressions": {
    "show_if": "user.role === 'admin'",
    "calculate": "price * quantity * (1 - discount)"
  }
}
```