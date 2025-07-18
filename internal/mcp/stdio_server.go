package mcp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"url-db/internal/models"
)

// StdioServer implements MCP protocol over stdin/stdout
type StdioServer struct {
	service MCPService
	reader  *bufio.Reader
	writer  io.Writer
}

// MCPRequest represents an MCP protocol request
type MCPRequest struct {
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params,omitempty"`
	ID     interface{}            `json:"id,omitempty"`
}

// MCPResponse represents an MCP protocol response
type MCPResponse struct {
	Result interface{} `json:"result,omitempty"`
	Error  *MCPError   `json:"error,omitempty"`
	ID     interface{} `json:"id,omitempty"`
}

// NewStdioServer creates a new MCP stdio server
func NewStdioServer(service MCPService) *StdioServer {
	return &StdioServer{
		service: service,
		reader:  bufio.NewReader(os.Stdin),
		writer:  os.Stdout,
	}
}

// Start begins the stdio MCP session
func (s *StdioServer) Start() error {
	log.Println("MCP stdio server ready for commands")
	log.Println("Available commands: list_domains, list_nodes, create_node, get_node, update_node, delete_node, server_info, quit")
	
	for {
		fmt.Fprint(s.writer, "> ")
		
		line, err := s.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("EOF received, ending session")
				return nil
			}
			return fmt.Errorf("error reading input: %w", err)
		}
		
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Handle quit command
		if line == "quit" || line == "exit" {
			log.Println("Goodbye!")
			return nil
		}
		
		// Parse and handle the request
		if err := s.handleRequest(line); err != nil {
			fmt.Fprintf(s.writer, "Error: %v\n", err)
		}
	}
}

// handleRequest processes a single MCP request
func (s *StdioServer) handleRequest(input string) error {
	ctx := context.Background()
	
	// Split input into command and arguments
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}
	
	command := parts[0]
	args := parts[1:]
	
	switch command {
	case "list_domains":
		return s.handleListDomains(ctx)
		
	case "list_nodes":
		if len(args) < 1 {
			return fmt.Errorf("list_nodes requires domain_name argument")
		}
		return s.handleListNodes(ctx, args[0])
		
	case "create_node":
		if len(args) < 2 {
			return fmt.Errorf("create_node requires domain_name and url arguments")
		}
		title := ""
		if len(args) > 2 {
			title = strings.Join(args[2:], " ")
		}
		return s.handleCreateNode(ctx, args[0], args[1], title)
		
	case "get_node":
		if len(args) < 1 {
			return fmt.Errorf("get_node requires composite_id argument")
		}
		return s.handleGetNode(ctx, args[0])
		
	case "update_node":
		if len(args) < 2 {
			return fmt.Errorf("update_node requires composite_id and title arguments")
		}
		title := strings.Join(args[1:], " ")
		return s.handleUpdateNode(ctx, args[0], title)
		
	case "delete_node":
		if len(args) < 1 {
			return fmt.Errorf("delete_node requires composite_id argument")
		}
		return s.handleDeleteNode(ctx, args[0])
		
	case "server_info":
		return s.handleServerInfo(ctx)
		
	case "help":
		return s.handleHelp()
		
	default:
		return fmt.Errorf("unknown command: %s. Type 'help' for available commands", command)
	}
}

func (s *StdioServer) handleListDomains(ctx context.Context) error {
	response, err := s.service.ListDomains(ctx)
	if err != nil {
		return err
	}
	
	fmt.Fprintf(s.writer, "Domains (%d):\n", len(response.Domains))
	for _, domain := range response.Domains {
		fmt.Fprintf(s.writer, "  - %s: %s (nodes: %d)\n", domain.Name, domain.Description, domain.NodeCount)
	}
	return nil
}

func (s *StdioServer) handleListNodes(ctx context.Context, domainName string) error {
	response, err := s.service.ListNodes(ctx, domainName, 1, 20, "")
	if err != nil {
		return err
	}
	
	fmt.Fprintf(s.writer, "Nodes in domain '%s' (%d):\n", domainName, response.TotalCount)
	for _, node := range response.Nodes {
		fmt.Fprintf(s.writer, "  - %s: %s\n", node.CompositeID, node.URL)
		if node.Title != "" {
			fmt.Fprintf(s.writer, "    Title: %s\n", node.Title)
		}
		if node.Description != "" {
			fmt.Fprintf(s.writer, "    Description: %s\n", node.Description)
		}
	}
	return nil
}

func (s *StdioServer) handleCreateNode(ctx context.Context, domainName, url, title string) error {
	req := &models.CreateMCPNodeRequest{
		DomainName:  domainName,
		URL:         url,
		Title:       title,
		Description: "",
	}
	
	node, err := s.service.CreateNode(ctx, req)
	if err != nil {
		return err
	}
	
	fmt.Fprintf(s.writer, "Created node: %s\n", node.CompositeID)
	fmt.Fprintf(s.writer, "  URL: %s\n", node.URL)
	if node.Title != "" {
		fmt.Fprintf(s.writer, "  Title: %s\n", node.Title)
	}
	return nil
}

func (s *StdioServer) handleGetNode(ctx context.Context, compositeID string) error {
	node, err := s.service.GetNode(ctx, compositeID)
	if err != nil {
		return err
	}
	
	fmt.Fprintf(s.writer, "Node: %s\n", node.CompositeID)
	fmt.Fprintf(s.writer, "  Domain: %s\n", node.DomainName)
	fmt.Fprintf(s.writer, "  URL: %s\n", node.URL)
	if node.Title != "" {
		fmt.Fprintf(s.writer, "  Title: %s\n", node.Title)
	}
	if node.Description != "" {
		fmt.Fprintf(s.writer, "  Description: %s\n", node.Description)
	}
	fmt.Fprintf(s.writer, "  Created: %s\n", node.CreatedAt)
	fmt.Fprintf(s.writer, "  Updated: %s\n", node.UpdatedAt)
	return nil
}

func (s *StdioServer) handleUpdateNode(ctx context.Context, compositeID, title string) error {
	req := &models.UpdateNodeRequest{
		Title:       title,
		Description: "",
	}
	
	node, err := s.service.UpdateNode(ctx, compositeID, req)
	if err != nil {
		return err
	}
	
	fmt.Fprintf(s.writer, "Updated node: %s\n", node.CompositeID)
	fmt.Fprintf(s.writer, "  Title: %s\n", node.Title)
	return nil
}

func (s *StdioServer) handleDeleteNode(ctx context.Context, compositeID string) error {
	err := s.service.DeleteNode(ctx, compositeID)
	if err != nil {
		return err
	}
	
	fmt.Fprintf(s.writer, "Deleted node: %s\n", compositeID)
	return nil
}

func (s *StdioServer) handleServerInfo(ctx context.Context) error {
	info, err := s.service.GetServerInfo(ctx)
	if err != nil {
		return err
	}
	
	fmt.Fprintf(s.writer, "Server Info:\n")
	fmt.Fprintf(s.writer, "  Name: %s\n", info.Name)
	fmt.Fprintf(s.writer, "  Version: %s\n", info.Version)
	fmt.Fprintf(s.writer, "  Description: %s\n", info.Description)
	fmt.Fprintf(s.writer, "  Composite Key Format: %s\n", info.CompositeKeyFormat)
	fmt.Fprintf(s.writer, "  Capabilities: %v\n", info.Capabilities)
	return nil
}

func (s *StdioServer) handleHelp() error {
	fmt.Fprintf(s.writer, "Available commands:\n")
	fmt.Fprintf(s.writer, "  list_domains                           - List all domains\n")
	fmt.Fprintf(s.writer, "  list_nodes <domain_name>               - List nodes in domain\n")
	fmt.Fprintf(s.writer, "  create_node <domain> <url> [title]     - Create new node\n")
	fmt.Fprintf(s.writer, "  get_node <composite_id>                - Get node details\n")
	fmt.Fprintf(s.writer, "  update_node <composite_id> <title>     - Update node title\n")
	fmt.Fprintf(s.writer, "  delete_node <composite_id>             - Delete node\n")
	fmt.Fprintf(s.writer, "  server_info                            - Show server information\n")
	fmt.Fprintf(s.writer, "  help                                   - Show this help\n")
	fmt.Fprintf(s.writer, "  quit                                   - Exit the session\n")
	fmt.Fprintf(s.writer, "\n")
	fmt.Fprintf(s.writer, "Example: create_node example.com https://example.com/page \"My Page\"\n")
	return nil
}