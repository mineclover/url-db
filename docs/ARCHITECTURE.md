# Architecture Improvement Plan

## Current Issues
- Mixed domain/layer organization
- Duplicated handler structures  
- Scattered test files
- Inconsistent naming

## Target Structure (Clean Architecture)

```
internal/
├── core/                    # Business Logic (Domain Layer)
│   ├── domain/             # Domain entities & interfaces
│   │   ├── models/         # Core business models
│   │   ├── repositories/   # Repository interfaces
│   │   └── services/       # Service interfaces
│   └── usecases/          # Application Use Cases
│       ├── domain/        # Domain use cases
│       ├── node/          # Node use cases
│       └── attribute/     # Attribute use cases
│
├── infrastructure/         # External Concerns (Infrastructure Layer)
│   ├── database/          # Database implementations
│   ├── repositories/      # Repository implementations
│   └── config/           # Configuration
│
├── interfaces/            # Interface Adapters (Interface Layer)
│   ├── http/             # HTTP handlers & routes
│   ├── mcp/              # MCP protocol handlers
│   └── cli/              # CLI interfaces
│
└── shared/               # Shared Utilities
    ├── errors/           # Error definitions
    ├── utils/            # Common utilities
    └── constants/        # Application constants
```

## Migration Strategy
1. Create new structure alongside old
2. Move files gradually with tests
3. Update imports progressively
4. Remove old structure
5. Improve test coverage simultaneously

## Benefits
- Clear separation of concerns
- Easier testing (dependency injection)
- Better maintainability
- Consistent organization
- Higher test coverage potential