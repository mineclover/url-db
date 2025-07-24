# 템플릿 시스템 사용 예제

## 1. 제품 카탈로그 템플릿

### 시나리오
전자상거래 사이트에서 제품 페이지 레이아웃을 표준화하기 위한 템플릿

### 템플릿 생성
```json
POST /api/templates
{
  "name": "product-detail-template",
  "domain_name": "ecommerce",
  "title": "제품 상세 페이지 템플릿",
  "description": "모든 제품 페이지에서 사용할 표준 레이아웃",
  "template_data": {
    "version": "1.0",
    "type": "layout",
    "metadata": {
      "author": "design-team@company.com",
      "created": "2024-01-15T10:00:00Z",
      "tags": ["product", "detail", "responsive"],
      "category": "ecommerce"
    },
    "content": {
      "structure": {
        "type": "grid",
        "container": {
          "maxWidth": "1200px",
          "padding": "20px",
          "margin": "0 auto"
        },
        "areas": [
          {
            "name": "gallery",
            "gridArea": "1 / 1 / 3 / 2",
            "width": "60%"
          },
          {
            "name": "info",
            "gridArea": "1 / 2 / 2 / 3",
            "width": "40%"
          },
          {
            "name": "actions",
            "gridArea": "2 / 2 / 3 / 3",
            "width": "40%"
          },
          {
            "name": "details",
            "gridArea": "3 / 1 / 4 / 3",
            "width": "100%"
          }
        ]
      },
      "components": {
        "gallery": {
          "type": "image-carousel",
          "config": {
            "showThumbnails": true,
            "enableZoom": true,
            "maxImages": 10
          }
        },
        "info": {
          "fields": ["title", "price", "rating", "short_description"],
          "layout": "vertical"
        },
        "actions": {
          "buttons": ["add_to_cart", "buy_now", "wishlist"],
          "showStock": true
        },
        "details": {
          "tabs": ["description", "specifications", "reviews", "shipping"]
        }
      }
    },
    "presentation": {
      "theme": "light",
      "responsive": true,
      "breakpoints": {
        "mobile": "768px",
        "tablet": "1024px"
      }
    }
  }
}
```

### 템플릿에 속성 추가
```json
POST /api/template-attributes
{
  "composite_id": "url-db:ecommerce:template:1",
  "attributes": [
    { "name": "category", "value": "layout" },
    { "name": "usage", "value": "product-pages" },
    { "name": "performance", "value": "optimized" },
    { "name": "seo_score", "value": "95" }
  ]
}
```

## 2. 블로그 포스트 템플릿

### 시나리오
블로그 플랫폼에서 일관된 포스트 구조를 유지하기 위한 템플릿

### 템플릿 생성
```json
POST /api/templates
{
  "name": "blog-post-standard",
  "domain_name": "content",
  "title": "표준 블로그 포스트 템플릿",
  "description": "SEO 최적화된 블로그 포스트 구조",
  "template_data": {
    "version": "2.1",
    "type": "document",
    "metadata": {
      "author": "content-team@blog.com",
      "tags": ["blog", "article", "seo-optimized"]
    },
    "schema": {
      "sections": [
        {
          "id": "meta",
          "type": "header",
          "required": true,
          "fields": {
            "title": {
              "type": "text",
              "required": true,
              "minLength": 30,
              "maxLength": 60,
              "placeholder": "SEO 친화적인 제목 (30-60자)"
            },
            "slug": {
              "type": "text",
              "required": true,
              "pattern": "^[a-z0-9-]+$",
              "generated": true
            },
            "excerpt": {
              "type": "textarea",
              "required": true,
              "minLength": 120,
              "maxLength": 160,
              "placeholder": "메타 설명 (120-160자)"
            },
            "featured_image": {
              "type": "image",
              "required": true,
              "dimensions": {
                "width": 1200,
                "height": 630
              },
              "alt_text_required": true
            },
            "author": {
              "type": "author_select",
              "required": true
            },
            "category": {
              "type": "select",
              "required": true,
              "options": ["Technology", "Design", "Business", "Lifestyle"]
            },
            "tags": {
              "type": "tags",
              "required": true,
              "min": 3,
              "max": 10
            }
          }
        },
        {
          "id": "content",
          "type": "body",
          "required": true,
          "blocks": [
            {
              "type": "introduction",
              "required": true,
              "min": 1,
              "max": 1,
              "minWords": 50,
              "maxWords": 150
            },
            {
              "type": "heading",
              "levels": [2, 3],
              "min": 3,
              "max": 10
            },
            {
              "type": "paragraph",
              "min": 5,
              "minWords": 50
            },
            {
              "type": "image",
              "max": 10,
              "caption": true,
              "alt_required": true
            },
            {
              "type": "code",
              "syntax_highlight": true,
              "copy_button": true
            },
            {
              "type": "quote",
              "citation_required": true
            },
            {
              "type": "list",
              "min": 1
            }
          ],
          "toc": {
            "enabled": true,
            "min_headings": 3,
            "position": "sidebar"
          }
        },
        {
          "id": "conclusion",
          "type": "footer",
          "required": true,
          "components": {
            "summary": {
              "required": true,
              "minWords": 50,
              "maxWords": 200
            },
            "cta": {
              "required": true,
              "types": ["newsletter", "related_posts", "social_share"]
            },
            "author_bio": {
              "required": true,
              "show_social": true
            }
          }
        }
      ]
    },
    "validation": {
      "total_words": {
        "min": 800,
        "max": 3000
      },
      "readability": {
        "target_score": 60,
        "algorithm": "flesch-kincaid"
      },
      "seo": {
        "keyword_density": {
          "min": 0.5,
          "max": 2.5
        },
        "internal_links": {
          "min": 2
        },
        "external_links": {
          "max": 5,
          "nofollow": true
        }
      }
    },
    "presentation": {
      "layout": "single-column",
      "typography": {
        "font_family": "Georgia, serif",
        "line_height": 1.8,
        "paragraph_spacing": "1.5em"
      },
      "code_theme": "monokai",
      "show_reading_time": true,
      "enable_comments": true
    }
  }
}
```

## 3. API 문서 템플릿

### 시나리오
API 엔드포인트 문서화를 위한 표준 템플릿

### 템플릿 생성
```json
POST /api/templates
{
  "name": "api-endpoint-doc",
  "domain_name": "documentation",
  "title": "API 엔드포인트 문서 템플릿",
  "description": "RESTful API 엔드포인트 문서화 표준",
  "template_data": {
    "version": "1.0",
    "type": "custom",
    "metadata": {
      "author": "api-team@company.com",
      "tags": ["api", "documentation", "rest"]
    },
    "content": {
      "endpoint": {
        "base_url": "https://api.example.com/v1",
        "path": "/resources/{id}",
        "method": "PUT",
        "summary": "Update a resource",
        "description": "Updates an existing resource with the provided data"
      },
      "authentication": {
        "type": "Bearer Token",
        "required": true,
        "location": "header",
        "example": "Bearer eyJhbGciOiJIUzI1NiIs..."
      },
      "parameters": {
        "path": [
          {
            "name": "id",
            "type": "string",
            "required": true,
            "description": "The unique identifier of the resource",
            "example": "res_123456"
          }
        ],
        "query": [
          {
            "name": "include",
            "type": "string",
            "required": false,
            "description": "Comma-separated list of related resources to include",
            "example": "author,tags",
            "enum": ["author", "tags", "comments", "metadata"]
          }
        ],
        "headers": {
          "Content-Type": {
            "required": true,
            "value": "application/json"
          },
          "X-Request-ID": {
            "required": false,
            "description": "Unique request identifier for tracking"
          }
        }
      },
      "request": {
        "body": {
          "type": "application/json",
          "schema": {
            "type": "object",
            "required": ["name"],
            "properties": {
              "name": {
                "type": "string",
                "minLength": 1,
                "maxLength": 255,
                "description": "The name of the resource"
              },
              "description": {
                "type": "string",
                "maxLength": 1000,
                "description": "Optional description"
              },
              "metadata": {
                "type": "object",
                "description": "Additional metadata",
                "additionalProperties": true
              },
              "tags": {
                "type": "array",
                "items": {
                  "type": "string"
                },
                "maxItems": 20
              }
            }
          },
          "examples": {
            "basic": {
              "name": "Updated Resource Name",
              "description": "This is an updated description"
            },
            "complete": {
              "name": "Complete Update Example",
              "description": "Full update with all fields",
              "metadata": {
                "category": "important",
                "priority": "high"
              },
              "tags": ["updated", "example", "api"]
            }
          }
        }
      },
      "responses": {
        "200": {
          "description": "Resource updated successfully",
          "body": {
            "type": "application/json",
            "schema": {
              "type": "object",
              "properties": {
                "id": { "type": "string" },
                "name": { "type": "string" },
                "description": { "type": "string" },
                "metadata": { "type": "object" },
                "tags": { 
                  "type": "array",
                  "items": { "type": "string" }
                },
                "updated_at": { 
                  "type": "string",
                  "format": "date-time"
                }
              }
            },
            "example": {
              "id": "res_123456",
              "name": "Updated Resource Name",
              "description": "This is an updated description",
              "metadata": {
                "category": "important",
                "priority": "high"
              },
              "tags": ["updated", "example", "api"],
              "updated_at": "2024-01-15T14:30:00Z"
            }
          }
        },
        "400": {
          "description": "Bad Request",
          "body": {
            "type": "application/json",
            "schema": {
              "$ref": "#/components/schemas/Error"
            },
            "examples": {
              "validation_error": {
                "error": {
                  "code": "VALIDATION_ERROR",
                  "message": "Validation failed",
                  "details": [
                    {
                      "field": "name",
                      "message": "Name is required"
                    }
                  ]
                }
              }
            }
          }
        },
        "404": {
          "description": "Resource not found",
          "body": {
            "type": "application/json",
            "example": {
              "error": {
                "code": "NOT_FOUND",
                "message": "Resource with id 'res_123456' not found"
              }
            }
          }
        }
      },
      "code_samples": [
        {
          "language": "curl",
          "code": "curl -X PUT https://api.example.com/v1/resources/res_123456 \\\n  -H 'Authorization: Bearer YOUR_TOKEN' \\\n  -H 'Content-Type: application/json' \\\n  -d '{\n    \"name\": \"Updated Resource Name\",\n    \"description\": \"This is an updated description\"\n  }'"
        },
        {
          "language": "javascript",
          "code": "const response = await fetch('https://api.example.com/v1/resources/res_123456', {\n  method: 'PUT',\n  headers: {\n    'Authorization': 'Bearer YOUR_TOKEN',\n    'Content-Type': 'application/json'\n  },\n  body: JSON.stringify({\n    name: 'Updated Resource Name',\n    description: 'This is an updated description'\n  })\n});\n\nconst data = await response.json();"
        },
        {
          "language": "python",
          "code": "import requests\n\nresponse = requests.put(\n    'https://api.example.com/v1/resources/res_123456',\n    headers={\n        'Authorization': 'Bearer YOUR_TOKEN',\n        'Content-Type': 'application/json'\n    },\n    json={\n        'name': 'Updated Resource Name',\n        'description': 'This is an updated description'\n    }\n)\n\ndata = response.json()"
        }
      ]
    }
  }
}
```

## 4. 사용자 등록 폼 템플릿

### 시나리오
다양한 서비스에서 재사용 가능한 사용자 등록 폼 템플릿

### 템플릿 생성
```json
POST /api/templates
{
  "name": "user-registration-form",
  "domain_name": "auth",
  "title": "사용자 등록 폼 템플릿",
  "description": "GDPR 준수 사용자 등록 폼",
  "template_data": {
    "version": "3.0",
    "type": "form",
    "metadata": {
      "compliance": ["GDPR", "CCPA"],
      "last_audit": "2024-01-10T00:00:00Z"
    },
    "schema": {
      "fields": [
        {
          "name": "email",
          "type": "text",
          "label": "이메일 주소",
          "required": true,
          "placeholder": "your@email.com",
          "validation": {
            "type": "email",
            "unique": true,
            "async_check": "/api/check-email"
          },
          "autocomplete": "email"
        },
        {
          "name": "password",
          "type": "password",
          "label": "비밀번호",
          "required": true,
          "placeholder": "최소 8자, 대소문자 및 숫자 포함",
          "validation": {
            "minLength": 8,
            "pattern": "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d).+$",
            "strength_meter": true
          },
          "autocomplete": "new-password"
        },
        {
          "name": "password_confirm",
          "type": "password",
          "label": "비밀번호 확인",
          "required": true,
          "placeholder": "비밀번호를 다시 입력하세요",
          "validation": {
            "match": "password"
          },
          "autocomplete": "new-password"
        },
        {
          "name": "full_name",
          "type": "text",
          "label": "이름",
          "required": true,
          "placeholder": "홍길동",
          "maxLength": 100,
          "autocomplete": "name"
        },
        {
          "name": "phone",
          "type": "tel",
          "label": "전화번호",
          "required": false,
          "placeholder": "+82-10-1234-5678",
          "validation": {
            "pattern": "^\\+?[0-9\\s\\-\\(\\)]+$"
          },
          "autocomplete": "tel"
        },
        {
          "name": "birth_date",
          "type": "date",
          "label": "생년월일",
          "required": true,
          "validation": {
            "min_age": 13,
            "max_age": 120
          }
        },
        {
          "name": "country",
          "type": "select",
          "label": "국가",
          "required": true,
          "options": "dynamic:/api/countries",
          "default": "KR",
          "searchable": true
        },
        {
          "name": "marketing_consent",
          "type": "checkbox",
          "label": "마케팅 정보 수신 동의",
          "required": false,
          "default": false,
          "help_text": "제품 업데이트 및 프로모션 정보를 이메일로 받아보실 수 있습니다."
        },
        {
          "name": "terms_consent",
          "type": "checkbox",
          "label": "이용약관 동의",
          "required": true,
          "validation": {
            "must_be_true": true
          },
          "link": {
            "text": "이용약관",
            "url": "/terms"
          }
        },
        {
          "name": "privacy_consent",
          "type": "checkbox",
          "label": "개인정보 처리방침 동의",
          "required": true,
          "validation": {
            "must_be_true": true
          },
          "link": {
            "text": "개인정보 처리방침",
            "url": "/privacy"
          }
        }
      ],
      "sections": [
        {
          "title": "계정 정보",
          "fields": ["email", "password", "password_confirm"],
          "icon": "account"
        },
        {
          "title": "개인 정보",
          "fields": ["full_name", "phone", "birth_date", "country"],
          "icon": "person"
        },
        {
          "title": "약관 동의",
          "fields": ["marketing_consent", "terms_consent", "privacy_consent"],
          "icon": "document"
        }
      ]
    },
    "validation": {
      "rules": {
        "email": {
          "required": "이메일 주소는 필수입니다",
          "email": "올바른 이메일 형식이 아닙니다",
          "unique": "이미 사용 중인 이메일입니다"
        },
        "password": {
          "required": "비밀번호는 필수입니다",
          "minLength": "비밀번호는 최소 8자 이상이어야 합니다",
          "pattern": "비밀번호는 대문자, 소문자, 숫자를 포함해야 합니다"
        },
        "password_confirm": {
          "required": "비밀번호 확인은 필수입니다",
          "match": "비밀번호가 일치하지 않습니다"
        }
      },
      "async_validation": {
        "debounce": 500,
        "show_spinner": true
      }
    },
    "presentation": {
      "layout": "wizard",
      "steps": ["계정 정보", "개인 정보", "약관 동의"],
      "progress_bar": true,
      "save_draft": true,
      "theme": {
        "primary_color": "#007bff",
        "error_color": "#dc3545",
        "success_color": "#28a745"
      },
      "submit_button": {
        "text": "회원가입 완료",
        "loading_text": "처리 중...",
        "success_redirect": "/welcome"
      }
    },
    "plugins": [
      {
        "name": "recaptcha",
        "version": "3.0",
        "config": {
          "site_key": "YOUR_RECAPTCHA_SITE_KEY",
          "threshold": 0.5
        }
      },
      {
        "name": "analytics",
        "version": "1.0",
        "config": {
          "track_steps": true,
          "track_errors": true
        }
      }
    ]
  }
}
```

## 5. 템플릿 활용 예제

### 템플릿 복제 및 커스터마이징
```json
POST /api/templates/clone
{
  "source_composite_id": "url-db:ecommerce:template:1",
  "new_name": "product-detail-custom",
  "modifications": {
    "template_data.content.structure.areas[0].width": "50%",
    "template_data.content.structure.areas[1].width": "50%",
    "template_data.presentation.theme": "dark"
  }
}
```

### 템플릿 기반 노드 생성
```json
POST /api/nodes/from-template
{
  "template_composite_id": "url-db:content:template:2",
  "node_data": {
    "url": "https://blog.example.com/my-first-post",
    "title": "템플릿을 활용한 첫 번째 블로그 포스트",
    "content": {
      "meta": {
        "title": "초보자를 위한 Go 프로그래밍 입문",
        "excerpt": "Go 언어의 기초부터 실무 활용까지, 단계별로 알아보는 종합 가이드",
        "author": "john.doe@blog.com",
        "category": "Technology",
        "tags": ["golang", "programming", "tutorial", "beginner"]
      }
    }
  }
}
```

### 템플릿 버전 비교
```json
GET /api/templates/compare?template1=url-db:auth:template:3&version1=2.0&version2=3.0

Response:
{
  "changes": [
    {
      "path": "schema.fields",
      "type": "added",
      "field": "phone",
      "description": "전화번호 필드 추가"
    },
    {
      "path": "metadata.compliance",
      "type": "added",
      "value": "CCPA",
      "description": "CCPA 규정 준수 추가"
    },
    {
      "path": "presentation.layout",
      "type": "modified",
      "old": "single-page",
      "new": "wizard",
      "description": "단일 페이지에서 마법사 형식으로 변경"
    }
  ]
}
```

## 모범 사례

### 1. 버전 관리
- Semantic Versioning 사용 (major.minor.patch)
- 하위 호환성 유지
- 변경 사항 문서화

### 2. 메타데이터 활용
- 작성자, 생성일, 수정일 기록
- 태그를 통한 분류 및 검색
- 권한 설정으로 접근 제어

### 3. 검증 규칙
- 클라이언트/서버 양쪽 검증
- 명확한 오류 메시지 제공
- 비동기 검증 최적화

### 4. 성능 최적화
- 자주 사용되는 템플릿 캐싱
- 큰 템플릿의 경우 압축 저장
- 불필요한 데이터 제거

### 5. 확장성
- 플러그인 시스템 활용
- 변수 및 표현식으로 동적 처리
- 템플릿 상속 구조 설계