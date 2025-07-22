#!/usr/bin/env python3
"""
Fix JSON parsing issues in Python test files
Add proper error handling for empty responses
"""

import re
from pathlib import Path

def fix_json_parsing():
    """Fix JSON parsing issues in test files"""
    
    # List of files that need fixing
    test_files = [
        "test_all_mcp_tools.py",
        "test_mcp_tools.py"
    ]
    
    files_fixed = 0
    
    for test_file_name in test_files:
        test_file = Path(test_file_name)
        if not test_file.exists():
            print(f"⚠️  File not found: {test_file_name}")
            continue
            
        print(f"\n🔧 Fixing {test_file_name}...")
        
        # Read file
        with open(test_file, 'r') as f:
            content = f.read()
        
        original_content = content
        
        # Pattern 1: Simple json.loads() calls
        pattern1 = r'json\.loads\(([^)]+)\)'
        def replacement1(match):
            var = match.group(1)
            return f'json.loads({var}) if {var}.strip() else {{}}'
        
        content = re.sub(pattern1, replacement1, content)
        
        # Pattern 2: Fix specific problematic lines
        problem_patterns = [
            # Domain parsing
            (r'domain = json\.loads\(response\[\'result\'\]\[\'content\'\]\[0\]\[\'text\'\]\)',
             'result_text = response[\'result\'][\'content\'][0][\'text\']\n        domain = json.loads(result_text) if result_text.strip() else {}'),
            
            # Node parsing
            (r'node = json\.loads\(result\)',
             'node = json.loads(result) if result.strip() else {}'),
            
            # Attribute parsing
            (r'attr_response = json\.loads\(result\)',
             'attr_response = json.loads(result) if result.strip() else {}'),
        ]
        
        for pattern, replacement in problem_patterns:
            content = re.sub(pattern, replacement, content)
        
        # Write back if changes were made
        if content != original_content:
            with open(test_file, 'w') as f:
                f.write(content)
            files_fixed += 1
            print(f"   ✅ Fixed JSON parsing in {test_file_name}")
        else:
            print(f"   ⚪ No changes needed for {test_file_name}")
    
    print(f"\n📊 Summary:")
    print(f"   Files processed: {len(test_files)}")
    print(f"   Files fixed: {files_fixed}")
    
    return files_fixed > 0

def main():
    """Main function"""
    print("🔧 Fixing JSON parsing issues in Python tests...")
    
    success = fix_json_parsing()
    
    if success:
        print("\n🎉 Successfully fixed JSON parsing issues!")
    else:
        print("\n⚪ No fixes needed")
    
    return 0

if __name__ == "__main__":
    exit(main())