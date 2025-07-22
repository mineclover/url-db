#!/usr/bin/env python3
"""
Demo script showing what the get_server_info tool would return
based on the implementation in the URL-DB MCP server.
"""

import json
from datetime import datetime
from tool_constants import CREATE_DOMAIN, CREATE_NODE, DELETE_NODE, FIND_NODE_BY_URL, GET_NODE, GET_NODE_ATTRIBUTES, GET_SERVER_INFO, LIST_DOMAINS, LIST_NODES, SET_NODE_ATTRIBUTES, UPDATE_NODE


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
            CREATE_NODE,
            GET_NODE,
            UPDATE_NODE,
            DELETE_NODE,
            LIST_NODES,
            FIND_NODE_BY_URL,
            "batch_get_nodes",
            LIST_DOMAINS,
            CREATE_DOMAIN,
            GET_NODE_ATTRIBUTES,
            SET_NODE_ATTRIBUTES
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
            "name": LIST_DOMAINS,
            "description": "List all domains in the URL database"
        },
        {
            "name": CREATE_DOMAIN, 
            "description": "Create a new domain"
        },
        {
            "name": LIST_NODES,
            "description": "List nodes in a specific domain"
        },
        {
            "name": CREATE_NODE,
            "description": "Create a new node (URL) in a domain"
        },
        {
            "name": GET_NODE,
            "description": "Get a node by composite ID"
        },
        {
            "name": UPDATE_NODE,
            "description": "Update a node's title and description"
        },
        {
            "name": DELETE_NODE,
            "description": "Delete a node by composite ID"
        },
        {
            "name": FIND_NODE_BY_URL,
            "description": "Find a node by URL in a domain"
        },
        {
            "name": GET_NODE_ATTRIBUTES,
            "description": "Get all attributes for a node"
        },
        {
            "name": SET_NODE_ATTRIBUTES,
            "description": "Set attributes for a node"
        },
        {
            "name": GET_SERVER_INFO,
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