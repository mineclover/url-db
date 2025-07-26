package mcp

import (
	"context"
	"encoding/json"
	"fmt"
)

// handleToolCall executes a tool call
func (h *MCPProtocolHandler) handleToolCall(ctx context.Context, req *JSONRPCRequest) *JSONRPCResponse {
	var params struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return h.createErrorResponse(req.ID, InvalidParams, "Invalid tool call parameters", err.Error())
	}

	// Use tool name directly without namespace
	toolName := params.Name

	var result interface{}
	var err error

	switch toolName {
	case "get_server_info":
		return h.handleGetServerInfo(req)
	case "list_domains":
		result, err = h.toolHandler.handleListDomains(ctx, params.Arguments)
	case "create_domain":
		result, err = h.toolHandler.handleCreateDomain(ctx, params.Arguments)
	case "list_nodes":
		result, err = h.toolHandler.handleListNodes(ctx, params.Arguments)
	case "create_node":
		result, err = h.toolHandler.handleCreateNode(ctx, params.Arguments)
	case "get_node":
		result, err = h.toolHandler.handleGetNode(ctx, params.Arguments)
	case "update_node":
		result, err = h.toolHandler.handleUpdateNode(ctx, params.Arguments)
	case "delete_node":
		result, err = h.toolHandler.handleDeleteNode(ctx, params.Arguments)
	case "find_node_by_url":
		result, err = h.toolHandler.handleFindNodeByURL(ctx, params.Arguments)
	case "scan_all_content":
		result, err = h.toolHandler.handleScanAllContent(ctx, params.Arguments)
	case "get_node_attributes":
		result, err = h.toolHandler.handleGetNodeAttributes(ctx, params.Arguments)
	case "set_node_attributes":
		result, err = h.toolHandler.handleSetNodeAttributes(ctx, params.Arguments)
	case "list_domain_attributes":
		result, err = h.toolHandler.handleListDomainAttributes(ctx, params.Arguments)
	case "create_domain_attribute":
		result, err = h.toolHandler.handleCreateDomainAttribute(ctx, params.Arguments)
	case "get_domain_attribute":
		result, err = h.toolHandler.handleGetDomainAttribute(ctx, params.Arguments)
	case "update_domain_attribute":
		result, err = h.toolHandler.handleUpdateDomainAttribute(ctx, params.Arguments)
	case "delete_domain_attribute":
		result, err = h.toolHandler.handleDeleteDomainAttribute(ctx, params.Arguments)
	case "create_dependency":
		result, err = h.toolHandler.handleCreateDependency(ctx, params.Arguments)
	case "list_node_dependencies":
		result, err = h.toolHandler.handleListNodeDependencies(ctx, params.Arguments)
	case "list_node_dependents":
		result, err = h.toolHandler.handleListNodeDependents(ctx, params.Arguments)
	case "delete_dependency":
		result, err = h.toolHandler.handleDeleteDependency(ctx, params.Arguments)
	case "filter_nodes_by_attributes":
		result, err = h.toolHandler.handleFilterNodesByAttributes(ctx, params.Arguments)
	case "get_node_with_attributes":
		result, err = h.toolHandler.handleGetNodeWithAttributes(ctx, params.Arguments)
	case "list_templates":
		result, err = h.toolHandler.handleListTemplates(ctx, params.Arguments)
	case "create_template":
		result, err = h.toolHandler.handleCreateTemplate(ctx, params.Arguments)
	case "get_template":
		result, err = h.toolHandler.handleGetTemplate(ctx, params.Arguments)
	case "update_template":
		result, err = h.toolHandler.handleUpdateTemplate(ctx, params.Arguments)
	case "delete_template":
		result, err = h.toolHandler.handleDeleteTemplate(ctx, params.Arguments)
	case "clone_template":
		result, err = h.toolHandler.handleCloneTemplate(ctx, params.Arguments)
	case "generate_template_scaffold":
		result, err = h.toolHandler.handleGenerateTemplateScaffold(ctx, params.Arguments)
	case "validate_template":
		result, err = h.toolHandler.handleValidateTemplate(ctx, params.Arguments)
	default:
		return h.createErrorResponse(req.ID, MethodNotFound, fmt.Sprintf("Tool not found: %s", params.Name), nil)
	}

	// Handle the response
	if err != nil {
		return h.createErrorResponse(req.ID, InternalError, "Tool execution failed", err.Error())
	}

	return h.createSuccessResponse(req.ID, result)
}
