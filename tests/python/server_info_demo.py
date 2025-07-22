#!/usr/bin/env python3
"""
Demo script showing what the get_server_info tool would return
based on the implementation in the URL-DB MCP server.
"""

import json
from datetime import datetime

def get_server_info():
    """
    Simulates the get_server_info tool response based on the Go implementation.
    """
    
    # Based on MetadataManager.GetServerInfo() in internal/mcp/metadata.go
    server_info = {
        "name": "url-db",
        "version": "1.0.0",  # Default version from metadata.go
        "description": "URL 데이터베이스 MCP 서버",  # URL Database MCP Server
        "capabilities": [
            "resources",
            "tools", 
            "prompts",
            "sampling"
        ],
        "composite_key_format": "url-db:domain_name:id"
    }
    
    return server_info

def get_detailed_server_info():
    """
    Simulates the detailed server info that includes additional metadata.
    """
    
    # Based on MetadataManager.GetDetailedServerInfo() 
    detailed_info = {
        "name": "url-db",
        "version": "1.0.0",
        "description": "URL 데이터베이스 MCP 서버",
        "build_time": datetime.now().isoformat(),
        "git_commit": "unknown",  # Default when not provided
        "environment": "development",  # Default environment
        "go_version": "go1.21",  # Example Go version
        "platform": "darwin/amd64",  # Example platform
        "composite_key_format": "url-db:domain_name:id",
        "capabilities": [
            "resources",
            "tools",
            "prompts", 
            "sampling"
        ],
        "supported_operations": [
            "create_node",
            "get_node",
            "update_node",
            "delete_node",
            "list_nodes",
            "find_node_by_url",
            "batch_get_nodes",
            "list_domains",
            "create_domain",
            "get_node_attributes",
            "set_node_attributes"
        ],
        "limits": {
            "max_batch_size": 100,
            "max_page_size": 100,
            "max_url_length": 2048,
            "max_title_length": 255,
            "max_description_length": 1000,
            "max_domain_name_length": 50,
            "max_attribute_value_length": 2048
        }
    }
    
    return detailed_info

def list_available_tools():
    """
    Lists all available MCP tools based on the tool registry.
    """
    
    # Based on registerTools() in internal/mcp/tools.go
    tools = [
        {
            "name": "list_domains",
            "description": "List all domains in the URL database"
        },
        {
            "name": "create_domain", 
            "description": "Create a new domain"
        },
        {
            "name": "list_nodes",
            "description": "List nodes in a specific domain"
        },
        {
            "name": "create_node",
            "description": "Create a new node (URL) in a domain"
        },
        {
            "name": "get_node",
            "description": "Get a node by composite ID"
        },
        {
            "name": "update_node",
            "description": "Update a node's title and description"
        },
        {
            "name": "delete_node",
            "description": "Delete a node by composite ID"
        },
        {
            "name": "find_node_by_url",
            "description": "Find a node by URL in a domain"
        },
        {
            "name": "get_node_attributes",
            "description": "Get all attributes for a node"
        },
        {
            "name": "set_node_attributes",
            "description": "Set attributes for a node"
        },
        {
            "name": "get_server_info",
            "description": "Get server information and capabilities"
        }
    ]
    
    return tools

def main():
    print("URL-DB MCP Server Information Demo")
    print("=" * 50)
    
    print("\n1. Basic Server Info (get_server_info):")
    print(json.dumps(get_server_info(), indent=2))
    
    print("\n2. Detailed Server Info:")
    print(json.dumps(get_detailed_server_info(), indent=2))
    
    print("\n3. Available MCP Tools:")
    tools = list_available_tools()
    for tool in tools:
        print(f"  - {tool['name']}: {tool['description']}")
    
    print("\n4. Key Features:")
    print("  - Manages URLs organized by domains")
    print("  - Supports node attributes for metadata")
    print("  - Uses composite keys format: url-db:domain_name:id")
    print("  - Implements full MCP protocol with tools and resources")
    print("  - Auto-creates domains when needed")
    print("  - Supports batch operations for efficiency")

if __name__ == "__main__":
    main()