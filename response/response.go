package response

// HTTPError represents response body if there is an error
type HTTPError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
