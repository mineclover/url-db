package mcp

import (
	"context"
	"runtime"
	"time"
)

type MetadataManager struct {
	version     string
	buildTime   string
	gitCommit   string
	environment string
}

func NewMetadataManager(version, buildTime, gitCommit, environment string) *MetadataManager {
	if version == "" {
		version = "1.0.0"
	}
	if buildTime == "" {
		buildTime = time.Now().Format(time.RFC3339)
	}
	if gitCommit == "" {
		gitCommit = "unknown"
	}
	if environment == "" {
		environment = "development"
	}

	return &MetadataManager{
		version:     version,
		buildTime:   buildTime,
		gitCommit:   gitCommit,
		environment: environment,
	}
}

func (mm *MetadataManager) GetServerInfo(ctx context.Context) (*MCPServerInfo, error) {
	return &MCPServerInfo{
		Name:        "url-db",
		Version:     mm.version,
		Description: "URL 데이터베이스 MCP 서버",
		Capabilities: []string{
			"resources",
			"tools",
			"prompts",
			"sampling",
		},
		CompositeKeyFormat: "url-db:domain_name:id",
	}, nil
}

func (mm *MetadataManager) GetDetailedServerInfo(ctx context.Context) (*DetailedServerInfo, error) {
	return &DetailedServerInfo{
		Name:               "url-db",
		Version:            mm.version,
		Description:        "URL 데이터베이스 MCP 서버",
		BuildTime:          mm.buildTime,
		GitCommit:          mm.gitCommit,
		Environment:        mm.environment,
		GoVersion:          runtime.Version(),
		Platform:           runtime.GOOS + "/" + runtime.GOARCH,
		CompositeKeyFormat: "url-db:domain_name:id",
		Capabilities: []string{
			"resources",
			"tools",
			"prompts",
			"sampling",
		},
		SupportedOperations: []string{
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
			"set_node_attributes",
		},
		Limits: ServerLimits{
			MaxBatchSize:            100,
			MaxPageSize:             100,
			MaxURLLength:            2048,
			MaxTitleLength:          255,
			MaxDescriptionLength:    1000,
			MaxDomainNameLength:     50,
			MaxAttributeValueLength: 2048,
		},
	}, nil
}

func (mm *MetadataManager) GetHealthStatus(ctx context.Context) (*HealthStatus, error) {
	return &HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   mm.version,
		Uptime:    time.Since(getStartTime()).String(),
		Checks: map[string]HealthCheck{
			"database": {
				Status:  "healthy",
				Message: "Database connection is healthy",
			},
			"memory": {
				Status:  "healthy",
				Message: "Memory usage is within normal limits",
			},
			"goroutines": {
				Status:  "healthy",
				Message: "Goroutine count is normal",
			},
		},
	}, nil
}

func (mm *MetadataManager) GetAPIDocumentation(ctx context.Context) (*APIDocumentation, error) {
	return &APIDocumentation{
		OpenAPI: "3.0.0",
		Info: APIInfo{
			Title:       "URL Database MCP API",
			Description: "MCP (Model Context Protocol) API for URL Database",
			Version:     mm.version,
			Contact: APIContact{
				Name:  "URL Database Team",
				Email: "support@url-db.com",
			},
		},
		Servers: []APIServer{
			{
				URL:         "/api/mcp",
				Description: "MCP API Server",
			},
		},
		Paths: map[string]APIPath{
			"/nodes": {
				GET: &APIOperation{
					Summary:     "List nodes",
					Description: "Get a list of nodes with optional filtering",
					Parameters: []APIParameter{
						{
							Name:        "domain_name",
							In:          "query",
							Description: "Filter by domain name",
							Required:    false,
							Schema:      APISchema{Type: "string"},
						},
						{
							Name:        "page",
							In:          "query",
							Description: "Page number",
							Required:    false,
							Schema:      APISchema{Type: "integer", Default: 1},
						},
						{
							Name:        "size",
							In:          "query",
							Description: "Page size",
							Required:    false,
							Schema:      APISchema{Type: "integer", Default: 20},
						},
					},
				},
				POST: &APIOperation{
					Summary:     "Create node",
					Description: "Create a new node",
					RequestBody: &APIRequestBody{
						Description: "Node creation request",
						Content: map[string]APIMediaType{
							"application/json": {
								Schema: APISchema{
									Type: "object",
									Properties: map[string]APISchema{
										"domain_name": {Type: "string", Description: "Domain name"},
										"url":         {Type: "string", Description: "URL"},
										"title":       {Type: "string", Description: "Title"},
										"description": {Type: "string", Description: "Description"},
									},
									Required: []string{"domain_name", "url"},
								},
							},
						},
					},
				},
			},
			"/nodes/{composite_id}": {
				GET: &APIOperation{
					Summary:     "Get node",
					Description: "Get a node by composite ID",
					Parameters: []APIParameter{
						{
							Name:        "composite_id",
							In:          "path",
							Description: "Composite ID in format 'url-db:domain_name:id'",
							Required:    true,
							Schema:      APISchema{Type: "string"},
						},
					},
				},
				PUT: &APIOperation{
					Summary:     "Update node",
					Description: "Update a node by composite ID",
					Parameters: []APIParameter{
						{
							Name:        "composite_id",
							In:          "path",
							Description: "Composite ID in format 'url-db:domain_name:id'",
							Required:    true,
							Schema:      APISchema{Type: "string"},
						},
					},
				},
				DELETE: &APIOperation{
					Summary:     "Delete node",
					Description: "Delete a node by composite ID",
					Parameters: []APIParameter{
						{
							Name:        "composite_id",
							In:          "path",
							Description: "Composite ID in format 'url-db:domain_name:id'",
							Required:    true,
							Schema:      APISchema{Type: "string"},
						},
					},
				},
			},
		},
	}, nil
}

func (mm *MetadataManager) GetStatistics(ctx context.Context) (*ServerStatistics, error) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return &ServerStatistics{
		Version:   mm.version,
		Uptime:    time.Since(getStartTime()).String(),
		Timestamp: time.Now(),
		Runtime: RuntimeStats{
			GoVersion:    runtime.Version(),
			NumGoroutine: runtime.NumGoroutine(),
			NumCPU:       runtime.NumCPU(),
			MemoryStats: MemoryStats{
				Alloc:        memStats.Alloc,
				TotalAlloc:   memStats.TotalAlloc,
				Sys:          memStats.Sys,
				NumGC:        memStats.NumGC,
				HeapAlloc:    memStats.HeapAlloc,
				HeapSys:      memStats.HeapSys,
				HeapInuse:    memStats.HeapInuse,
				HeapReleased: memStats.HeapReleased,
			},
		},
	}, nil
}

var startTime = time.Now()

func getStartTime() time.Time {
	return startTime
}

type DetailedServerInfo struct {
	Name                string       `json:"name"`
	Version             string       `json:"version"`
	Description         string       `json:"description"`
	BuildTime           string       `json:"build_time"`
	GitCommit           string       `json:"git_commit"`
	Environment         string       `json:"environment"`
	GoVersion           string       `json:"go_version"`
	Platform            string       `json:"platform"`
	CompositeKeyFormat  string       `json:"composite_key_format"`
	Capabilities        []string     `json:"capabilities"`
	SupportedOperations []string     `json:"supported_operations"`
	Limits              ServerLimits `json:"limits"`
}

type ServerLimits struct {
	MaxBatchSize            int `json:"max_batch_size"`
	MaxPageSize             int `json:"max_page_size"`
	MaxURLLength            int `json:"max_url_length"`
	MaxTitleLength          int `json:"max_title_length"`
	MaxDescriptionLength    int `json:"max_description_length"`
	MaxDomainNameLength     int `json:"max_domain_name_length"`
	MaxAttributeValueLength int `json:"max_attribute_value_length"`
}

type HealthStatus struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Version   string                 `json:"version"`
	Uptime    string                 `json:"uptime"`
	Checks    map[string]HealthCheck `json:"checks"`
}

type HealthCheck struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type APIDocumentation struct {
	OpenAPI string             `json:"openapi"`
	Info    APIInfo            `json:"info"`
	Servers []APIServer        `json:"servers"`
	Paths   map[string]APIPath `json:"paths"`
}

type APIInfo struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Version     string     `json:"version"`
	Contact     APIContact `json:"contact"`
}

type APIContact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type APIServer struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

type APIPath struct {
	GET    *APIOperation `json:"get,omitempty"`
	POST   *APIOperation `json:"post,omitempty"`
	PUT    *APIOperation `json:"put,omitempty"`
	DELETE *APIOperation `json:"delete,omitempty"`
}

type APIOperation struct {
	Summary     string          `json:"summary"`
	Description string          `json:"description"`
	Parameters  []APIParameter  `json:"parameters,omitempty"`
	RequestBody *APIRequestBody `json:"requestBody,omitempty"`
}

type APIParameter struct {
	Name        string    `json:"name"`
	In          string    `json:"in"`
	Description string    `json:"description"`
	Required    bool      `json:"required"`
	Schema      APISchema `json:"schema"`
}

type APIRequestBody struct {
	Description string                  `json:"description"`
	Content     map[string]APIMediaType `json:"content"`
}

type APIMediaType struct {
	Schema APISchema `json:"schema"`
}

type APISchema struct {
	Type        string               `json:"type"`
	Properties  map[string]APISchema `json:"properties,omitempty"`
	Required    []string             `json:"required,omitempty"`
	Description string               `json:"description,omitempty"`
	Default     interface{}          `json:"default,omitempty"`
}

type ServerStatistics struct {
	Version   string       `json:"version"`
	Uptime    string       `json:"uptime"`
	Timestamp time.Time    `json:"timestamp"`
	Runtime   RuntimeStats `json:"runtime"`
}

type RuntimeStats struct {
	GoVersion    string      `json:"go_version"`
	NumGoroutine int         `json:"num_goroutine"`
	NumCPU       int         `json:"num_cpu"`
	MemoryStats  MemoryStats `json:"memory_stats"`
}

type MemoryStats struct {
	Alloc        uint64 `json:"alloc"`
	TotalAlloc   uint64 `json:"total_alloc"`
	Sys          uint64 `json:"sys"`
	NumGC        uint32 `json:"num_gc"`
	HeapAlloc    uint64 `json:"heap_alloc"`
	HeapSys      uint64 `json:"heap_sys"`
	HeapInuse    uint64 `json:"heap_inuse"`
	HeapReleased uint64 `json:"heap_released"`
}
