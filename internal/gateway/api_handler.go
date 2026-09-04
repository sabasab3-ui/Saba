package gateway

import (
	"context"
	"encoding/json"
	"fmt"
)

// APIHandler handles incoming API requests
type APIHandler struct {
	gateway *Gateway
}

// TaskRequest represents an incoming task request
type TaskRequest struct {
	TaskID      string            `json:"task_id"`
	Type        string            `json:"type"`     // agent, automation, analysis
	Payload     map[string]interface{} `json:"payload"`
	Mode        string            `json:"mode"`
	UserID      string            `json:"user_id"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// TaskResponse represents the API response
type TaskResponse struct {
	TaskID     string          `json:"task_id"`
	Status     string          `json:"status"`
	Result     interface{}     `json:"result"`
	Error      string          `json:"error,omitempty"`
	Duration   int64           `json:"duration_ms"`
	Timestamp  string          `json:"timestamp"`
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(gateway *Gateway) *APIHandler {
	return &APIHandler{
		gateway: gateway,
	}
}

// HandleTask processes an incoming task request
func (ah *APIHandler) HandleTask(ctx context.Context, req *TaskRequest) (*TaskResponse, error) {
	if req.TaskID == "" {
		return nil, fmt.Errorf("task_id is required")
	}

	if req.Type == "" {
		return nil, fmt.Errorf("type is required")
	}

	response := &TaskResponse{
		TaskID:    req.TaskID,
		Status:    "processing",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	switch req.Type {
	case "agent":
		// Route to agent handler
		result, err := ah.gateway.ProcessAgentTask(ctx, req)
		if err != nil {
			response.Status = "failed"
			response.Error = err.Error()
		} else {
			response.Status = "completed"
			response.Result = result
		}

	case "automation":
		// Route to automation handler
		result, err := ah.gateway.ProcessAutomation(ctx, req)
		if err != nil {
			response.Status = "failed"
			response.Error = err.Error()
		} else {
			response.Status = "completed"
			response.Result = result
		}

	case "analysis":
		// Route to analysis handler
		result, err := ah.gateway.PerformAnalysis(ctx, req)
		if err != nil {
			response.Status = "failed"
			response.Error = err.Error()
		} else {
			response.Status = "completed"
			response.Result = result
		}

	default:
		response.Status = "failed"
		response.Error = fmt.Sprintf("unknown type: %s", req.Type)
	}

	return response, nil
}

// HandleBulkTasks processes multiple tasks
func (ah *APIHandler) HandleBulkTasks(ctx context.Context, requests []*TaskRequest) ([]*TaskResponse, error) {
	responses := make([]*TaskResponse, len(requests))

	for i, req := range requests {
		resp, err := ah.HandleTask(ctx, req)
		if err != nil {
			responses[i] = &TaskResponse{
				TaskID: req.TaskID,
				Status: "error",
				Error:  err.Error(),
			}
		} else {
			responses[i] = resp
		}
	}

	return responses, nil
}

// GetTaskStatus retrieves the status of a task
func (ah *APIHandler) GetTaskStatus(ctx context.Context, taskID string) (*TaskResponse, error) {
	// Implementation would query task store
	return &TaskResponse{
		TaskID:    taskID,
		Status:    "completed",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}
