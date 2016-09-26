package goSmartSheet

import "encoding/json"

//PaginatedResponse Returned by the API in certain scenarios
type PaginatedResponse struct {
	PageNumber int             `json:"pageNumber"`
	PageSize   int             `json:"pageSize"`
	TotalPages int             `json:"totalPages"`
	TotalCount int             `json:"totalCount"`
	Data       json.RawMessage `json:"data"`
}

//Response is the generic response from the SmartSheet API
type Response struct {
	Message    string          `json:"message"`
	ResultCode int             `json:"resultCode"`
	Result     json.RawMessage `json:"result"`
	Version    int             `json:"version"`
}
