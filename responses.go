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

//BulkOperationFailure is the generic response from the SmartSheet API
type BulkOperationFailure struct {
	Message     string            `json:"message"`
	ResultCode  int               `json:"resultCode"`
	Result      json.RawMessage   `json:"result"`
	Version     int               `json:"version"`
	FailedItems []BulkItemFailure `json:"failedItems"`
}

//BulkItemFailure is a failure when a bulk operation occurs
type BulkItemFailure struct {
	Index   int       `json:"index"`
	Failure ErrorItem `json:"error"`
	RowID   int       `json:"rowId"`
}

//ErrorItem is an erorr item from the SmartSheet API
type ErrorItem struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message"`
}

//SingleOperationFailure reprsents a single failure during an operation
type SingleOperationFailure struct {
	ErrorCode int                   `json:"errorCode"`
	Message   string                `json:"message"`
	RefID     string                `json:"refId"`
	Details   []SingleFailureDetail `json:"details"`
}

//SingleFailureDetail is the detail for a single failure
type SingleFailureDetail struct {
	Index int   `json:"index"`
	RowID int64 `json:"rowId"`
}
