#!/usr/bin/env python3
"""
Fix Python test files by removing incorrect imports and applying correct constants
"""

import re
from pathlib import Path

def fix_python_imports():
    """Fix Python import statements and apply correct constants"""
    
    # Define correct mapping
    tool_mappings = {
        'LIST_DOMAINS_TOOL': 'LIST_DOMAINS',
        'CREATE_DOMAIN_TOOL': 'CREATE_DOMAIN',
        'LIST_NODES_TOOL': 'LIST_NODES',
        'CREATE_NODE_TOOL': 'CREATE_NODE',
        'GET_NODE_TOOL': 'GET_NODE',
        'UPDATE_NODE_TOOL': 'UPDATE_NODE',
        'DELETE_NODE_TOOL': 'DELETE_NODE',
        'FIND_NODE_BY_URL_TOOL': 'FIND_NODE_BY_URL',
        'GET_NODE_ATTRIBUTES_TOOL': 'GET_NODE_ATTRIBUTES',
        'SET_NODE_ATTRIBUTES_TOOL': 'SET_NODE_ATTRIBUTES',
        'LIST_DOMAIN_ATTRIBUTES_TOOL': 'LIST_DOMAIN_ATTRIBUTES',
        'CREATE_DOMAIN_ATTRIBUTE_TOOL': 'CREATE_DOMAIN_ATTRIBUTE',
        'GET_DOMAIN_ATTRIBUTE_TOOL': 'GET_DOMAIN_ATTRIBUTE',
        'UPDATE_DOMAIN_ATTRIBUTE_TOOL': 'UPDATE_DOMAIN_ATTRIBUTE',
        'DELETE_DOMAIN_ATTRIBUTE_TOOL': 'DELETE_DOMAIN_ATTRIBUTE',
        'GET_NODE_WITH_ATTRIBUTES_TOOL': 'GET_NODE_WITH_ATTRIBUTES',
        'FILTER_NODES_BY_ATTRIBUTES_TOOL': 'FILTER_NODES_BY_ATTRIBUTES',
        'GET_SERVER_INFO_TOOL': 'GET_SERVER_INFO',
    }
    
    # List of correct constants
    correct_constants = list(tool_mappings.values())
    
    # Process test files
    test_dir = Path("tests/python")
    test_files = list(test_dir.glob("*.py"))
    
    files_fixed = 0
    
    for test_file in test_files:
        if test_file.name == "tool_constants.py":
            continue
            
        print(f"\n🔧 Fixing {test_file.name}...")
        
        # Read file
        with open(test_file, 'r') as f:
            content = f.read()
        
        original_content = content
        
        # Remove old incorrect import lines
        import_pattern = r'^from tool_constants import.*$'
        content = re.sub(import_pattern, '', content, flags=re.MULTILINE)
        
        # Replace incorrect constant names with correct ones
        for old_const, new_const in tool_mappings.items():
            content = content.replace(old_const, new_const)
        
        # Add correct import at the top
        if any(const in content for const in correct_constants):
            # Find where to insert import
            lines = content.split('\n')
            
            # Remove empty lines at the top
            while lines and lines[0].strip() == '':
                lines.pop(0)
            
            # Find import section
            insert_pos = 0
            for i, line in enumerate(lines):
                if line.startswith('import ') or line.startswith('from '):
                    insert_pos = i + 1
                elif line.strip() == '' and insert_pos > 0:
                    break
            
            # Create import line
            used_constants = [const for const in correct_constants if const in content]
            if used_constants:
                import_line = f"from tool_constants import {', '.join(sorted(used_constants))}"
                lines.insert(insert_pos, import_line)
                content = '\n'.join(lines)
                print(f"   📝 Added import: {len(used_constants)} constants")
        
        # Write back if changes were made
        if content != original_content:
            with open(test_file, 'w') as f:
                f.write(content)
            files_fixed += 1
            print(f"   ✅ Fixed {test_file.name}")
        else:
            print(f"   ⚪ No changes needed for {test_file.name}")
    
    print(f"\n📊 Summary:")
    print(f"   Files processed: {len(test_files) - 1}")
    print(f"   Files fixed: {files_fixed}")
    
    return files_fixed > 0

def main():
    """Main function"""
    print("🔧 Fixing Python import statements...")
    
    success = fix_python_imports()
    
    if success:
        print("\n🎉 Successfully fixed Python imports!")
    else:
        print("\n⚪ No fixes needed")
    
    return 0

if __name__ == "__main__":
    exit(main())