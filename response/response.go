package response

// HTTP represents response body
type HTTP struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
}
