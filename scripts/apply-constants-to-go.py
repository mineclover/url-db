#!/usr/bin/env python3
"""
Apply generated constants to Go code
Replace hardcoded tool names in tools.go with constants
"""

import re
from pathlib import Path

def apply_constants_to_tools_go():
    """Apply constants to internal/mcp/tools.go"""
    
    tools_file = Path("internal/mcp/tools.go")
    
    if not tools_file.exists():
        print(f"❌ File not found: {tools_file}")
        return False
    
    # Read the file
    with open(tools_file, 'r') as f:
        content = f.read()
    
    # Define mapping of hardcoded names to constants
    tool_mappings = {
        '"list_domains"': 'ListDomainsTool',
        '"create_domain"': 'CreateDomainTool',
        '"list_nodes"': 'ListNodesTool',
        '"create_node"': 'CreateNodeTool',
        '"get_node"': 'GetNodeTool',
        '"update_node"': 'UpdateNodeTool',
        '"delete_node"': 'DeleteNodeTool',
        '"find_node_by_url"': 'FindNodeByUrlTool',
        '"get_node_attributes"': 'GetNodeAttributesTool',
        '"set_node_attributes"': 'SetNodeAttributesTool',
        '"list_domain_attributes"': 'ListDomainAttributesTool',
        '"create_domain_attribute"': 'CreateDomainAttributeTool',
        '"get_domain_attribute"': 'GetDomainAttributeTool',
        '"update_domain_attribute"': 'UpdateDomainAttributeTool',
        '"delete_domain_attribute"': 'DeleteDomainAttributeTool',
        '"get_node_with_attributes"': 'GetNodeWithAttributesTool',
        '"filter_nodes_by_attributes"': 'FilterNodesByAttributesTool',
        '"get_server_info"': 'GetServerInfoTool',
    }
    
    # Apply replacements
    original_content = content
    for hardcoded, constant in tool_mappings.items():
        # Replace in Name fields (with proper spacing)
        pattern = r'Name:\s*' + re.escape(hardcoded)
        replacement = f'Name: {constant}'
        content = re.sub(pattern, replacement, content)
    
    # Check if any changes were made
    if content == original_content:
        print("⚠️  No hardcoded tool names found to replace in Name fields")
        return False
    
    # Write back to file
    with open(tools_file, 'w') as f:
        f.write(content)
    
    # Count replacements
    changes = sum(1 for hardcoded in tool_mappings if hardcoded in original_content)
    print(f"✅ Applied constants to {tools_file}")
    print(f"   Replaced {changes} hardcoded tool names with constants")
    
    return True

def main():
    """Main function"""
    print("🔧 Applying generated constants to Go code...")
    
    success = apply_constants_to_tools_go()
    
    if success:
        print("\n🎉 Successfully applied constants to Go code!")
        print("\nNext steps:")
        print("1. Compile and test the Go code")
        print("2. Apply constants to Python test files")
        print("3. Run full test suite")
    else:
        print("\n❌ Failed to apply constants")
        return 1
    
    return 0

if __name__ == "__main__":
    exit(main())