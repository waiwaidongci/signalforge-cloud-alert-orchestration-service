package http

import (
	"time"
)

type ingestRequest struct {
	ExternalID  string            `json:"external_id"`
	Resource    string            `json:"resource"`
	Severity    string            `json:"severity"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	Timestamp   time.Time         `json:"timestamp"`
}

type batchIDsRequest struct {
	IDs   []string `json:"ids"`
	Actor string   `json:"actor"`
}
