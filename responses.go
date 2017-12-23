package goSmartSheet

import (
	"encoding/json"
	"fmt"
	"io"

	"time"

	"github.com/pkg/errors"
)

//PaginatedResponse Returned by the API in certain scenarios
type PaginatedResponse struct {
	PageNumber int             `json:"pageNumber"`
	PageSize   int             `json:"pageSize"`
	TotalPages int             `json:"totalPages"`
	TotalCount int             `json:"totalCount"`
	Data       json.RawMessage `json:"data"`
}

/*
2017/05/22 00:32:09 Status Code: 200
2017/05/22 00:32:09 {
	"message": "SUCCESS",
	"resultCode": 0,
	"result": [
		{
			"id": 5135307589871492,
			"rowNumber": 1,
			"expanded": true,
			"createdAt": "2017-02-20T21:43:27Z",
			"modifiedAt": "2017-05-22T05:30:27Z",
*/

//Response is the generic response from the SmartSheet API
type Response struct {
	Message    string `json:"message"`
	ResultCode int    `json:"resultCode"`
	Version    int    `json:"version"`
}

//RowAlterResponse is the generic response when altering rows from the SmartSheet API
type RowAlterResponse struct {
	Response
	Result []RowResponse `json:"result"`
}

//RowResponse is the individual response for each row of data altered in the call to the SS API
type RowResponse struct {
	ID         string    `json:"id"`
	RowNumber  int       `json:"rowNumber"`
	Expanded   bool      `json:"expanded"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
}

//ResultResponse is the generic result including a returned object
type ResultResponse struct {
	Response
	Result      json.RawMessage   `json:"result"`
	FailedItems []BulkItemFailure `json:"failedItems,omitempty"`
}

//BulkItemFailure is a failure when a bulk operation occurs
type BulkItemFailure struct {
	Index   int       `json:"index"`
	Failure ErrorItem `json:"error"`
	RowID   int       `json:"rowId"`
}

/*
Example Response:

2017/05/22 00:27:25 Status Code: 400
2017/05/22 00:27:25 {
  "errorCode" : 1042,
  "message" : "The value for cell in column 7192636189632388, update test 1, did not conform to the strict requirements for type PICKLIST.",
  "refId" : "sa3bo1hh1g8",
  "detail" : {
    "index" : 0,
    "rowId" : 1651223794345860
  }
}

*/

//ErrorItem reprsents a single failure during an operation
type ErrorItem struct {
	ErrorCode  int               `json:"errorCode"`
	Message    string            `json:"message"`
	RefID      string            `json:"refId,omitempty"`
	Details    []ErrorItemDetail `json:"details,omitempty"`
	StatusCode int               `json:"-"` //not part of SS API, but used to provide extra context
}

//String returns a string representation of an ErrorItem for output purposes
func (e *ErrorItem) String() string {
	return fmt.Sprintf("Error Code: %v, Message: %v, RefId: %v", e.ErrorCode, e.Message, e.RefID)
}

//Error returns a string representation of an ErrorItem for output purposes
func (e *ErrorItem) Error() string {
	return fmt.Sprintf("Error (%v): %s", e.StatusCode, e.String())
}

//ErrorItemDecode translates the SmartSheet ErrorItem into a Go erorr
func ErrorItemDecode(statusCode int, bodyDec *json.Decoder) error {
	e := &ErrorItem{}
	if err := bodyDec.Decode(e); err != nil {
		return errors.Wrap(err, "Failed to decode into ErrorItem")
	}
	e.StatusCode = statusCode
	return e
}

//ErrorItemDecodeFromReader translates the SmartSheet ErrorItem into a Go erorr taking a ReadCloser
func ErrorItemDecodeFromReader(statusCode int, body io.ReadCloser) error {
	bodyDec := json.NewDecoder(body)
	defer body.Close()

	return ErrorItemDecode(statusCode, bodyDec)
}

//ErrorItemDetail is the detail for a single failure
type ErrorItemDetail struct {
	Index int   `json:"index"`
	RowID int64 `json:"rowId"`
}
