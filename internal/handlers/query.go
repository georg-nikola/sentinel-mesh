package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// QueryHandler handles query requests
type QueryHandler struct{}

// NewQueryHandler creates a new query handler
func NewQueryHandler() *QueryHandler {
	return &QueryHandler{}
}

// QueryRangeRequest represents a range query request
type QueryRangeRequest struct {
	Query     string    `json:"query"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Step      string    `json:"step"`
}

// QueryInstantRequest represents an instant query request
type QueryInstantRequest struct {
	Query string    `json:"query"`
	Time  time.Time `json:"time"`
}

// QueryResponse represents query response
type QueryResponse struct {
	Status string                 `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

// QueryRange handles range queries
func (h *QueryHandler) QueryRange(w http.ResponseWriter, r *http.Request) {
	var req QueryRangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Placeholder response - would query actual time-series DB
	response := QueryResponse{
		Status: "success",
		Data: map[string]interface{}{
			"resultType": "matrix",
			"result": []map[string]interface{}{
				{
					"metric": map[string]string{
						"__name__": "up",
						"job":      "prometheus",
					},
					"values": [][]interface{}{
						{req.StartTime.Unix(), "1"},
						{req.EndTime.Unix(), "1"},
					},
				},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// QueryInstant handles instant queries
func (h *QueryHandler) QueryInstant(w http.ResponseWriter, r *http.Request) {
	var req QueryInstantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Placeholder response
	response := QueryResponse{
		Status: "success",
		Data: map[string]interface{}{
			"resultType": "vector",
			"result": []map[string]interface{}{
				{
					"metric": map[string]string{
						"__name__": "up",
						"job":      "prometheus",
					},
					"value": []interface{}{req.Time.Unix(), "1"},
				},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
