# Scripts Directory

This directory contains automation scripts for the URL-DB MCP server project.

## Core Scripts (Keep Long-term)

### `generate-tool-constants.py` ⭐
**Purpose**: Generates consistent tool constants for Go and Python from YAML specification  
**Status**: Core infrastructure - Keep permanently  
**Usage**: `python3 generate-tool-constants.py`  
**Outputs**: 
- `generated/tool_constants.go` - Go constants
- `generated/tool_constants.py` - Python constants  
- `generated/mcp_tools_schema.json` - JSON schema

**Quality**: 9/10 - Excellent architecture, essential for consistency

### `test.sh` ⭐
**Purpose**: Comprehensive testing infrastructure for the project  
**Status**: Essential testing - Keep permanently  
**Usage**: `./test.sh [options]`  
**Features**: Unit tests, integration tests, coverage reports, linting

**Quality**: 9/10 - Well-structured, critical for CI/CD

### `test-cursor-mcp.sh` ⭐  
**Purpose**: Cursor IDE integration testing and configuration helper  
**Status**: User onboarding tool - Keep permanently  
**Usage**: `./test-cursor-mcp.sh`  
**Features**: MCP server testing, configuration examples

**Quality**: 8/10 - Clear, safe operations, valuable for users

## Utility Scripts (Conditional Keep)

### `apply-constants-to-go.py` 🔄
**Purpose**: Applies generated constants to Go source files  
**Status**: Needs refactoring - Currently has hardcoded mappings  
**Usage**: `python3 apply-constants-to-go.py`  
**Issue**: Brittle due to hardcoded tool mappings

**Quality**: 7/10 - Works but needs improvement  
**TODO**: Refactor to read from `generate-tool-constants.py` output

### `apply-constants-to-python.py` 🔄
**Purpose**: Applies generated constants to Python test files  
**Status**: Needs refactoring - Mixed responsibilities  
**Usage**: `python3 apply-constants-to-python.py`  
**Issue**: Monolithic, hardcoded mappings, mixed concerns

**Quality**: 6/10 - Complex, needs splitting  
**TODO**: Split into focused, reusable components

## Maintenance Guidelines

### Immediate Actions
1. **Add unit tests** to `generate-tool-constants.py`
2. **Document removal timeline** for utility scripts
3. **Ensure migration** is complete before removing scripts

### Short-term Improvements  
1. **Refactor** `apply-constants-to-go.py` to use generated constants
2. **Split** `apply-constants-to-python.py` into focused components
3. **Integrate** constant generation into build process

### Long-term Considerations
1. **Move test scripts** to `/scripts/testing/` subdirectory
2. **Add comprehensive error handling** across all scripts
3. **Document usage patterns** in project documentation

## Script Quality Matrix

| Script | Quality | Value | Action |
|--------|---------|-------|--------|
| `generate-tool-constants.py` | 9/10 | High | Keep + Improve |
| `test.sh` | 9/10 | High | Keep |
| `test-cursor-mcp.sh` | 8/10 | High | Keep |
| `apply-constants-to-go.py` | 7/10 | Medium | Refactor |
| `apply-constants-to-python.py` | 6/10 | Medium | Refactor |

## Usage Notes

- **Core scripts** should be maintained and improved
- **Utility scripts** need refactoring before long-term use
- **All scripts** should be tested before major changes
- **Generated files** should not be manually edited

---

*Last updated: 2025-07-22*  
*Analysis completed by: Claude Code Assistant*