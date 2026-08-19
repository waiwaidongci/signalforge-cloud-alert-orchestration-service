package http

type batchIncidentRequest struct {
	IDs   []string `json:"ids"`
	Actor string   `json:"actor"`
}
