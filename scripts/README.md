# Scripts Directory

This directory contains automation scripts for the URL-DB MCP server project.

## Available Scripts

### `generate-tool-constants.py` ⭐
**Purpose**: Generates consistent tool constants for Go and Python from YAML specification  
**Status**: Core infrastructure - Keep permanently  
**Usage**: `python3 generate-tool-constants.py`  
**Outputs**: 
- `generated/tool_constants.go` - Go constants
- `generated/tool_constants.py` - Python constants  
- `generated/mcp_tools_schema.json` - JSON schema

**Quality**: 9/10 - Excellent architecture, essential for consistency

## Usage Guidelines

### Core Script Usage
```bash
# Generate tool constants from YAML specification
python3 scripts/generate-tool-constants.py
```

### Generated Files
- **`generated/tool_constants.go`**: Go constants for MCP tools
- **`generated/tool_constants.py`**: Python constants for testing
- **`generated/mcp_tools_schema.json`**: JSON schema for validation

### Manual Application
Since automated application scripts were removed due to reliability issues, constants should be applied manually:

1. **Go Code**: Copy constants from `generated/tool_constants.go` to `internal/interfaces/mcp/tool_constants.go`
2. **Python Tests**: Copy constants from `generated/tool_constants.py` when Python test infrastructure is added

## Maintenance Guidelines

### Current Status
1. ✅ **Core script** `generate-tool-constants.py` is working perfectly
2. ✅ **Clean directory** with only essential scripts
3. ✅ **Accurate documentation** reflecting actual state

### Future Improvements
1. **Add unit tests** to `generate-tool-constants.py`
2. **Integrate constant generation** into build process
3. **Create testing infrastructure** when needed
4. **Add comprehensive error handling** to core script

### Long-term Considerations
1. **Move test scripts** to `/scripts/testing/` subdirectory when created
2. **Document usage patterns** in project documentation
3. **Automate constant application** when Go code structure stabilizes

## Script Quality Matrix

| Script | Quality | Value | Status |
|--------|---------|-------|--------|
| `generate-tool-constants.py` | 9/10 | High | ✅ Keep + Improve |

## Usage Notes

- **Core script** `generate-tool-constants.py` should be maintained and improved
- **Generated files** should not be manually edited
- **Constants should be applied manually** until automated solution is stable
- **Test scripts will be created** when testing infrastructure is needed

---

*Last updated: 2025-07-23*  
*Status: Clean and functional*