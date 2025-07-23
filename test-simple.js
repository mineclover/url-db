const { spawn } = require('child_process');

console.log('Starting MCP server test...');

const server = spawn('/Users/junwoobang/mcp/url-db/bin/url-db', [
  '-mcp-mode=stdio',
  '-db-path=/tmp/test-mcp.db'
]);

let responseBuffer = '';

server.stdout.on('data', (data) => {
  console.log('STDOUT:', data.toString());
  responseBuffer += data.toString();
  
  // Try to parse complete JSON responses
  const lines = responseBuffer.split('\n');
  for (let i = 0; i < lines.length - 1; i++) {
    if (lines[i].trim()) {
      try {
        const json = JSON.parse(lines[i]);
        console.log('Parsed response:', JSON.stringify(json, null, 2));
      } catch (e) {
        console.log('Failed to parse:', lines[i]);
      }
    }
  }
  responseBuffer = lines[lines.length - 1];
});

server.stderr.on('data', (data) => {
  console.log('STDERR:', data.toString());
});

server.on('close', (code) => {
  console.log(`Server exited with code ${code}`);
});

// Send initialize request
const initRequest = {
  jsonrpc: "2.0",
  id: 1,
  method: "initialize",
  params: {
    protocolVersion: "2024-11-05",
    capabilities: {
      tools: {}
    }
  }
};

console.log('Sending initialize request...');
server.stdin.write(JSON.stringify(initRequest) + '\n');

// Send tools/list request after a delay
setTimeout(() => {
  const toolsRequest = {
    jsonrpc: "2.0",
    id: 2,
    method: "tools/list",
    params: {}
  };
  
  console.log('Sending tools/list request...');
  server.stdin.write(JSON.stringify(toolsRequest) + '\n');
  
  // Close stdin after another delay
  setTimeout(() => {
    console.log('Closing stdin...');
    server.stdin.end();
  }, 1000);
}, 1000);