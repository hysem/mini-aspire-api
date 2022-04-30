package response

// APIResponse holds the response details
type APIResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}
