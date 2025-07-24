#!/bin/bash

# Demo script for scan_all_content tool with page-based navigation and compression

echo "=== URL-DB scan_all_content Demo ==="
echo ""

# Initialize the MCP session
echo "1. Initializing MCP session..."
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null | jq '.'

echo ""
echo "2. Creating a test domain..."
echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"create_domain","arguments":{"name":"demo","description":"Demo domain for scan testing"}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null | jq '.result'

echo ""
echo "3. Adding some test URLs..."
for i in {1..5}; do
  echo '{"jsonrpc":"2.0","id":'$((2+i))',"method":"tools/call","params":{"name":"create_node","arguments":{"domain_name":"demo","url":"https://example.com/page'$i'","title":"Page '$i' Title","description":"This is the description for page '$i'"}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null | jq -r '.result.composite_id // empty'
done

echo ""
echo "4. Creating domain attributes..."
echo '{"jsonrpc":"2.0","id":8,"method":"tools/call","params":{"name":"create_domain_attribute","arguments":{"domain_name":"demo","name":"category","type":"tag","description":"Content category"}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null | jq '.'
echo '{"jsonrpc":"2.0","id":9,"method":"tools/call","params":{"name":"create_domain_attribute","arguments":{"domain_name":"demo","name":"priority","type":"tag","description":"Content priority"}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null | jq '.'

echo ""
echo "5. Adding attributes to nodes (with duplicates for compression demo)..."
# Add same category to multiple nodes
for i in {1..3}; do
  echo '{"jsonrpc":"2.0","id":'$((9+i))',"method":"tools/call","params":{"name":"set_node_attributes","arguments":{"composite_id":"url-db:demo:'$i'","attributes":[{"name":"category","value":"tech"},{"name":"priority","value":"high"}]}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null | jq -r '.result // empty'
done

echo ""
echo "6. Scanning all content WITHOUT compression (page 1)..."
echo '{"jsonrpc":"2.0","id":20,"method":"tools/call","params":{"name":"scan_all_content","arguments":{"domain_name":"demo","max_tokens_per_page":2000,"page":1,"include_attributes":true,"compress_attributes":false}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null | jq '.result.result | {total_nodes: .metadata.total_nodes, current_page: .pagination.current_page, has_more: .pagination.has_more, items_count: (.items | length)}'

echo ""
echo "7. Scanning all content WITH compression (page 1)..."
echo '{"jsonrpc":"2.0","id":21,"method":"tools/call","params":{"name":"scan_all_content","arguments":{"domain_name":"demo","max_tokens_per_page":2000,"page":1,"include_attributes":true,"compress_attributes":true}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null | jq '.result.result | {compressed: .metadata.compressed_output, duplicates_removed: .metadata.attribute_summary.total_duplicates_removed, unique_values: .metadata.attribute_summary.unique_values}'

echo ""
echo "8. Testing page navigation (page 2 if available)..."
echo '{"jsonrpc":"2.0","id":22,"method":"tools/call","params":{"name":"scan_all_content","arguments":{"domain_name":"demo","max_tokens_per_page":500,"page":2,"include_attributes":true}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null | jq '.result.result | {current_page: .pagination.current_page, has_previous: .pagination.has_previous, items_on_page: (.items | length)}'

echo ""
echo "=== Demo Complete ==="