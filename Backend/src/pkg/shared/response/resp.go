package response

import "github.com/socialpay/socialpay/src/pkg/shared/pagination"

// ERROR RESPONSE
type ErrorResponse struct {
	Success bool     `json:"success"`
	Error   ApiError `json:"error"`
}

// SUCCESS RESPONSE (generic success without pagination)
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data" swaggertype:"object"`
}

// PAGINATED SUCCESS RESPONSE
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data" swaggertype:"object"`
	Pagination pagination.PaginationInfo  `json:"pagination"`
}

// API ERROR OBJ
type ApiError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// PAGINATION METADATA
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}
