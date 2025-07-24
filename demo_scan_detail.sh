#!/bin/bash

# Detailed demo script for scan_all_content tool features

echo "=== URL-DB scan_all_content Detailed Demo ==="
echo ""

# Remove old database and start fresh
rm -f url-db.sqlite

echo "1. Creating demo domain with many URLs for pagination..."
echo '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"create_domain","arguments":{"name":"demo","description":"Demo domain for pagination testing"}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null | jq -r '.result.content[].text'

echo ""
echo "2. Adding 10 URLs to test pagination..."
for i in {1..10}; do
  echo '{"jsonrpc":"2.0","id":'$((1+i))',"method":"tools/call","params":{"name":"create_node","arguments":{"domain_name":"demo","url":"https://example.com/article-'$i'","title":"Article '$i': Long Title to Consume More Tokens for Testing Pagination","description":"This is a much longer description for article '$i' that contains more words to increase the token count. We want to test how the pagination works when we have content that uses more tokens per node."}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null | jq -r '.result.composite_id // empty' | xargs -I {} echo "Created: {}"
done

echo ""
echo "3. Creating attributes for compression demo..."
echo '{"jsonrpc":"2.0","id":20,"method":"tools/call","params":{"name":"create_domain_attribute","arguments":{"domain_name":"demo","name":"category","type":"tag"}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null >/dev/null
echo '{"jsonrpc":"2.0","id":21,"method":"tools/call","params":{"name":"create_domain_attribute","arguments":{"domain_name":"demo","name":"priority","type":"tag"}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null >/dev/null
echo '{"jsonrpc":"2.0","id":22,"method":"tools/call","params":{"name":"create_domain_attribute","arguments":{"domain_name":"demo","name":"status","type":"tag"}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null >/dev/null
echo "Created: category, priority, status attributes"

echo ""
echo "4. Adding duplicate attributes across nodes..."
# Most nodes get "tech" category and "high" priority
for i in {1..7}; do
  echo '{"jsonrpc":"2.0","id":'$((22+i))',"method":"tools/call","params":{"name":"set_node_attributes","arguments":{"composite_id":"url-db:demo:'$i'","attributes":[{"name":"category","value":"tech"},{"name":"priority","value":"high"},{"name":"status","value":"published"}]}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null >/dev/null
done
# Some variation
echo '{"jsonrpc":"2.0","id":30,"method":"tools/call","params":{"name":"set_node_attributes","arguments":{"composite_id":"url-db:demo:8","attributes":[{"name":"category","value":"business"},{"name":"priority","value":"medium"},{"name":"status","value":"draft"}]}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null >/dev/null
echo '{"jsonrpc":"2.0","id":31,"method":"tools/call","params":{"name":"set_node_attributes","arguments":{"composite_id":"url-db:demo:9","attributes":[{"name":"category","value":"tech"},{"name":"priority","value":"low"},{"name":"status","value":"published"}]}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null >/dev/null
echo "Attributes added (7 nodes with same values, 2 with variations)"

echo ""
echo "5. Testing page navigation with small token limit..."
echo "--- Page 1 (500 tokens max) ---"
echo '{"jsonrpc":"2.0","id":40,"method":"tools/call","params":{"name":"scan_all_content","arguments":{"domain_name":"demo","max_tokens_per_page":500,"page":1,"include_attributes":false}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null | jq -r '.result.result | "Total nodes: \(.metadata.total_nodes)\nCurrent page: \(.pagination.current_page)/\(.pagination.total_pages)\nItems on page: \(.items | length)\nHas more pages: \(.pagination.has_more)\nCurrent tokens: \(.pagination.current_tokens)"'

echo ""
echo "--- Page 2 (500 tokens max) ---"
echo '{"jsonrpc":"2.0","id":41,"method":"tools/call","params":{"name":"scan_all_content","arguments":{"domain_name":"demo","max_tokens_per_page":500,"page":2,"include_attributes":false}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null | jq -r '.result.result | "Current page: \(.pagination.current_page)/\(.pagination.total_pages)\nItems on page: \(.items | length)\nHas previous: \(.pagination.has_previous)\nHas more: \(.pagination.has_more)"'

echo ""
echo "6. Comparing attribute compression..."
echo "--- WITHOUT compression ---"
echo '{"jsonrpc":"2.0","id":50,"method":"tools/call","params":{"name":"scan_all_content","arguments":{"domain_name":"demo","max_tokens_per_page":8000,"page":1,"include_attributes":true,"compress_attributes":false}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null | jq -r '.result.result | "Total attribute values: \([.items[].attributes[] | select(. != null)] | length)\nEstimated tokens: \(.pagination.current_tokens)"'

echo ""
echo "--- WITH compression ---"
echo '{"jsonrpc":"2.0","id":51,"method":"tools/call","params":{"name":"scan_all_content","arguments":{"domain_name":"demo","max_tokens_per_page":8000,"page":1,"include_attributes":true,"compress_attributes":true}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null | jq -r '.result.result | "Duplicates removed: \(.metadata.attribute_summary.total_duplicates_removed)\nUnique category values: \(.metadata.attribute_summary.unique_values.category | length)\nUnique priority values: \(.metadata.attribute_summary.unique_values.priority | length)\nEstimated tokens: \(.pagination.current_tokens)"'

echo ""
echo "7. Showing compressed attribute summary..."
echo '{"jsonrpc":"2.0","id":52,"method":"tools/call","params":{"name":"scan_all_content","arguments":{"domain_name":"demo","max_tokens_per_page":8000,"page":1,"include_attributes":true,"compress_attributes":true}}}' | ./bin/url-db -mcp-mode=stdio 2>/dev/null | jq '.result.result.metadata.attribute_summary | {
  most_common_category: .most_common_values.category,
  most_common_priority: .most_common_values.priority,
  category_values: .unique_values.category,
  priority_values: .unique_values.priority,
  value_counts: .value_counts
}'

echo ""
echo "=== Demo Complete ==="
echo "Key features demonstrated:"
echo "1. Page-based navigation (1,2,3... pages)"
echo "2. Attribute compression removes duplicates"
echo "3. Token-based pagination for AI context optimization"
echo "4. Compression statistics and unique value tracking"