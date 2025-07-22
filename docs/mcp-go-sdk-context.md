# MCP Go SDK Context

## Required Context: MCP Go SDK

**Repository**: https://github.com/modelcontextprotocol/go-sdk

### Overview

The Model Context Protocol (MCP) Go SDK provides a Go implementation for building MCP servers and clients. This SDK is the official reference implementation for the MCP protocol in Go.

### Key Components

1. **Core MCP Package** (`mcp`)
   - Client and server APIs
   - Protocol message handling
   - Type-safe tool invocation

2. **JSON Schema Package** (`jsonschema`)
   - JSON Schema validation
   - Type definitions

3. **JSON-RPC Package** (`jsonrpc`)
   - Custom transport implementations
   - Stdin/stdout communication

### Design Patterns

#### Client-Server Communication
```go
// Client connection pattern
client := mcp.NewClient(&mcp.Implementation{
    Name: "client-name", 
    Version: "version"
}, nil)

// Server tool definition
func SayHi(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[HiParams]) (*mcp.CallToolResultFor[any], error) {
    // Tool implementation
}
```

#### Key Features
- Flexible transport mechanisms
- Type-safe tool definitions  
- Context-based async operations
- Error handling and content return mechanisms

### Current Implementation Status

**Note**: The Go SDK is marked as "unreleased" and "unstable" as of 2024, so this project implements a custom MCP JSON-RPC 2.0 server based on the official specification rather than using the SDK directly.

### Implementation Alignment

Our URL-DB MCP server follows the same patterns and protocols as the Go SDK:

1. **JSON-RPC 2.0 Protocol**: Full compliance with MCP specification
2. **Tool Registry**: Similar pattern to SDK's tool management
3. **Resource System**: Aligned with SDK's resource handling
4. **Type Safety**: Go's type system for parameter validation
5. **Context Handling**: Proper context propagation throughout

### References

- [MCP Specification](https://spec.modelcontextprotocol.io/)
- [MCP Go SDK Repository](https://github.com/modelcontextprotocol/go-sdk)
- [JSON-RPC 2.0 Specification](https://www.jsonrpc.org/specification)

This context should be considered when evaluating MCP server implementations, testing scenarios, and integration patterns.