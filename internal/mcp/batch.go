package mcp

import (
	"context"
	"sync"
	"time"

	"github.com/url-db/internal/models"
)

type BatchProcessor struct {
	mcpService MCPService
	maxWorkers int
	timeout    time.Duration
}

type BatchResult struct {
	CompositeID string
	Node        *models.MCPNode
	Error       error
}

func NewBatchProcessor(mcpService MCPService, maxWorkers int, timeout time.Duration) *BatchProcessor {
	if maxWorkers <= 0 {
		maxWorkers = 10
	}
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	return &BatchProcessor{
		mcpService: mcpService,
		maxWorkers: maxWorkers,
		timeout:    timeout,
	}
}

func (bp *BatchProcessor) BatchGetNodes(ctx context.Context, compositeIDs []string) (*models.BatchMCPNodeResponse, error) {
	if len(compositeIDs) == 0 {
		return &models.BatchMCPNodeResponse{
			Nodes:    []models.MCPNode{},
			NotFound: []string{},
		}, nil
	}

	ctx, cancel := context.WithTimeout(ctx, bp.timeout)
	defer cancel()

	jobs := make(chan string, len(compositeIDs))
	results := make(chan BatchResult, len(compositeIDs))

	var wg sync.WaitGroup
	workerCount := min(bp.maxWorkers, len(compositeIDs))

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go bp.worker(ctx, jobs, results, &wg)
	}

	for _, compositeID := range compositeIDs {
		jobs <- compositeID
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	return bp.collectResults(results, len(compositeIDs))
}

func (bp *BatchProcessor) worker(ctx context.Context, jobs <-chan string, results chan<- BatchResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for compositeID := range jobs {
		select {
		case <-ctx.Done():
			results <- BatchResult{
				CompositeID: compositeID,
				Error:       ctx.Err(),
			}
			return
		default:
			node, err := bp.mcpService.GetNode(ctx, compositeID)
			results <- BatchResult{
				CompositeID: compositeID,
				Node:        node,
				Error:       err,
			}
		}
	}
}

func (bp *BatchProcessor) collectResults(results <-chan BatchResult, expectedCount int) (*models.BatchMCPNodeResponse, error) {
	nodes := make([]models.MCPNode, 0, expectedCount)
	notFound := make([]string, 0)

	for result := range results {
		if result.Error != nil {
			notFound = append(notFound, result.CompositeID)
		} else if result.Node != nil {
			nodes = append(nodes, *result.Node)
		}
	}

	return &models.BatchMCPNodeResponse{
		Nodes:    nodes,
		NotFound: notFound,
	}, nil
}

func (bp *BatchProcessor) BatchCreateNodes(ctx context.Context, requests []models.CreateMCPNodeRequest) (*BatchCreateResponse, error) {
	if len(requests) == 0 {
		return &BatchCreateResponse{
			Nodes:  []models.MCPNode{},
			Failed: []BatchCreateFailure{},
		}, nil
	}

	ctx, cancel := context.WithTimeout(ctx, bp.timeout)
	defer cancel()

	jobs := make(chan models.CreateMCPNodeRequest, len(requests))
	results := make(chan BatchCreateResult, len(requests))

	var wg sync.WaitGroup
	workerCount := min(bp.maxWorkers, len(requests))

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go bp.createWorker(ctx, jobs, results, &wg)
	}

	for _, request := range requests {
		jobs <- request
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	return bp.collectCreateResults(results, len(requests))
}

func (bp *BatchProcessor) createWorker(ctx context.Context, jobs <-chan models.CreateMCPNodeRequest, results chan<- BatchCreateResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for request := range jobs {
		select {
		case <-ctx.Done():
			results <- BatchCreateResult{
				Request: request,
				Error:   ctx.Err(),
			}
			return
		default:
			node, err := bp.mcpService.CreateNode(ctx, &request)
			results <- BatchCreateResult{
				Request: request,
				Node:    node,
				Error:   err,
			}
		}
	}
}

func (bp *BatchProcessor) collectCreateResults(results <-chan BatchCreateResult, expectedCount int) (*BatchCreateResponse, error) {
	nodes := make([]models.MCPNode, 0, expectedCount)
	failed := make([]BatchCreateFailure, 0)

	for result := range results {
		if result.Error != nil {
			failed = append(failed, BatchCreateFailure{
				Request: result.Request,
				Error:   result.Error.Error(),
			})
		} else if result.Node != nil {
			nodes = append(nodes, *result.Node)
		}
	}

	return &BatchCreateResponse{
		Nodes:  nodes,
		Failed: failed,
	}, nil
}

func (bp *BatchProcessor) BatchUpdateNodes(ctx context.Context, updates []BatchUpdateRequest) (*BatchUpdateResponse, error) {
	if len(updates) == 0 {
		return &BatchUpdateResponse{
			Nodes:  []models.MCPNode{},
			Failed: []BatchUpdateFailure{},
		}, nil
	}

	ctx, cancel := context.WithTimeout(ctx, bp.timeout)
	defer cancel()

	jobs := make(chan BatchUpdateRequest, len(updates))
	results := make(chan BatchUpdateResult, len(updates))

	var wg sync.WaitGroup
	workerCount := min(bp.maxWorkers, len(updates))

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go bp.updateWorker(ctx, jobs, results, &wg)
	}

	for _, update := range updates {
		jobs <- update
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	return bp.collectUpdateResults(results, len(updates))
}

func (bp *BatchProcessor) updateWorker(ctx context.Context, jobs <-chan BatchUpdateRequest, results chan<- BatchUpdateResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for request := range jobs {
		select {
		case <-ctx.Done():
			results <- BatchUpdateResult{
				CompositeID: request.CompositeID,
				Error:       ctx.Err(),
			}
			return
		default:
			node, err := bp.mcpService.UpdateNode(ctx, request.CompositeID, &request.UpdateRequest)
			results <- BatchUpdateResult{
				CompositeID: request.CompositeID,
				Node:        node,
				Error:       err,
			}
		}
	}
}

func (bp *BatchProcessor) collectUpdateResults(results <-chan BatchUpdateResult, expectedCount int) (*BatchUpdateResponse, error) {
	nodes := make([]models.MCPNode, 0, expectedCount)
	failed := make([]BatchUpdateFailure, 0)

	for result := range results {
		if result.Error != nil {
			failed = append(failed, BatchUpdateFailure{
				CompositeID: result.CompositeID,
				Error:       result.Error.Error(),
			})
		} else if result.Node != nil {
			nodes = append(nodes, *result.Node)
		}
	}

	return &BatchUpdateResponse{
		Nodes:  nodes,
		Failed: failed,
	}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type BatchCreateResult struct {
	Request models.CreateMCPNodeRequest
	Node    *models.MCPNode
	Error   error
}

type BatchCreateResponse struct {
	Nodes  []models.MCPNode      `json:"nodes"`
	Failed []BatchCreateFailure `json:"failed"`
}

type BatchCreateFailure struct {
	Request models.CreateMCPNodeRequest `json:"request"`
	Error   string                      `json:"error"`
}

type BatchUpdateRequest struct {
	CompositeID   string                    `json:"composite_id"`
	UpdateRequest models.UpdateNodeRequest `json:"update_request"`
}

type BatchUpdateResult struct {
	CompositeID string
	Node        *models.MCPNode
	Error       error
}

type BatchUpdateResponse struct {
	Nodes  []models.MCPNode      `json:"nodes"`
	Failed []BatchUpdateFailure `json:"failed"`
}

type BatchUpdateFailure struct {
	CompositeID string `json:"composite_id"`
	Error       string `json:"error"`
}

type BatchDeleteRequest struct {
	CompositeIDs []string `json:"composite_ids"`
}

type BatchDeleteResponse struct {
	Deleted []string              `json:"deleted"`
	Failed  []BatchDeleteFailure `json:"failed"`
}

type BatchDeleteFailure struct {
	CompositeID string `json:"composite_id"`
	Error       string `json:"error"`
}

func (bp *BatchProcessor) BatchDeleteNodes(ctx context.Context, compositeIDs []string) (*BatchDeleteResponse, error) {
	if len(compositeIDs) == 0 {
		return &BatchDeleteResponse{
			Deleted: []string{},
			Failed:  []BatchDeleteFailure{},
		}, nil
	}

	ctx, cancel := context.WithTimeout(ctx, bp.timeout)
	defer cancel()

	jobs := make(chan string, len(compositeIDs))
	results := make(chan BatchDeleteResult, len(compositeIDs))

	var wg sync.WaitGroup
	workerCount := min(bp.maxWorkers, len(compositeIDs))

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go bp.deleteWorker(ctx, jobs, results, &wg)
	}

	for _, compositeID := range compositeIDs {
		jobs <- compositeID
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	return bp.collectDeleteResults(results, len(compositeIDs))
}

func (bp *BatchProcessor) deleteWorker(ctx context.Context, jobs <-chan string, results chan<- BatchDeleteResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for compositeID := range jobs {
		select {
		case <-ctx.Done():
			results <- BatchDeleteResult{
				CompositeID: compositeID,
				Error:       ctx.Err(),
			}
			return
		default:
			err := bp.mcpService.DeleteNode(ctx, compositeID)
			results <- BatchDeleteResult{
				CompositeID: compositeID,
				Error:       err,
			}
		}
	}
}

func (bp *BatchProcessor) collectDeleteResults(results <-chan BatchDeleteResult, expectedCount int) (*BatchDeleteResponse, error) {
	deleted := make([]string, 0, expectedCount)
	failed := make([]BatchDeleteFailure, 0)

	for result := range results {
		if result.Error != nil {
			failed = append(failed, BatchDeleteFailure{
				CompositeID: result.CompositeID,
				Error:       result.Error.Error(),
			})
		} else {
			deleted = append(deleted, result.CompositeID)
		}
	}

	return &BatchDeleteResponse{
		Deleted: deleted,
		Failed:  failed,
	}, nil
}

type BatchDeleteResult struct {
	CompositeID string
	Error       error
}