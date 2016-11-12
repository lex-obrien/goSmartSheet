package goSmartSheet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

//Client is used to interact with the SamartSheet API
type Client struct {
	url    string
	apiKey string
	client *http.Client
}

//GetClient will return back a SmartSheet client based on the specified apiKey
//Currently, this will always point to the prouction API
func GetClient(apiKey string) *Client {
	//default to prod API
	api := &Client{url: "https://api.smartsheet.com/2.0", apiKey: apiKey}
	api.client = &http.Client{} //per docs clients should be made once, https://golang.org/pkg/net/http/

	return api
}

//GetSheetFilterCols returns a Sheet but filter to only the specified columns
//Columns are specified via the Column Id
func (c *Client) GetSheetFilterCols(id string, onlyTheseColumns []string) (Sheet, error) {
	filter := "columnIds=" + strings.Join(onlyTheseColumns, ",")
	return c.GetSheet(id, filter)
}

//GetSheet returns a sheet with the specified Id
func (c *Client) GetSheet(id, queryFilter string) (Sheet, error) {
	s := Sheet{}

	path := "sheets/" + id
	if queryFilter != "" {
		path += "?" + queryFilter
	}

	body, err := c.Get(path)
	if err != nil {
		return s, fmt.Errorf("Failed to get sheet (ID: %v): %v", id, err)

	}
	defer body.Close()

	dec := json.NewDecoder(body)
	if err := dec.Decode(&s); err != nil {
		return s, fmt.Errorf("Failed to decode: %v", err)
	}

	return s, nil
}

//GetColumns will return back the columns for the specified Sheet
func (c *Client) GetColumns(sheetID string) (cols []Column, err error) {
	path := fmt.Sprintf("sheets/%v/columns", sheetID)

	body, err := c.Get(path)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var resp PaginatedResponse
	//TODO: need generic handling and ability to read from pages to get all data... eventually
	dec := json.NewDecoder(body)
	if err = dec.Decode(&resp); err != nil {
		return nil, fmt.Errorf("Failed to decode response: %v", err)
	}

	if err = json.Unmarshal(resp.Data, &cols); err != nil {
		return nil, fmt.Errorf("Failed to decode columns: %v", err)
	}

	return
}

//GetJSONString with return a Json string of the result
func (c *Client) GetJSONString(path string, prettify bool) (string, error) {
	body, err := c.Get(path)
	if err != nil {
		return "", fmt.Errorf("Failed to Get JSON String: %v", err)
	}
	defer body.Close()

	buf := new(bytes.Buffer)
	var s string

	if prettify {
		var m json.RawMessage

		dec := json.NewDecoder(body)
		if err := dec.Decode(&m); err != nil {
			return "", fmt.Errorf("Failed to decode: %v\n", err)
		}

		b, err := json.MarshalIndent(&m, "", "\t")
		if err != nil {
			return "", fmt.Errorf("Error during indent: %v\n", err)
		}

		s = string(b)
	} else {
		buf.ReadFrom(body)
		s = buf.String()
	}

	return s, nil
}

//AddRowToSheet will add a single row of data to an existing smartsheet by ID based on the specified cellValues
func (c *Client) AddRowToSheet(sheetID string, rowOpt RowPostOptions, cellValues ...CellValue) (io.ReadCloser, error) {
	var r Row

	for i := range cellValues {
		c := Cell{Value: &cellValues[i]}
		r.Cells = append(r.Cells, c)
	}

	return c.AddRowsToSheet(sheetID, rowOpt, []Row{r}, NormalValidation)
}

//AddRowsToSheet will add the specified rows to a sheet based on ID
func (c *Client) AddRowsToSheet(sheetID string, rowOpt RowPostOptions, rows []Row, opt PostOptions) (io.ReadCloser, error) {

	//adjust each row to match values
	var sheetCols []Column
	var err error
	var colsPopulated bool
	for i := range rows {
		r := &rows[i]

		//adjust col IDs
		for j := range r.Cells {
			//only fill the col Id if they are missing
			if r.Cells[j].ColumnID == 0 {
				//columnId is missing, so we need to perform some validation

				if !colsPopulated {
					sheetCols, err = c.GetColumns(sheetID)
					colsPopulated = true
					if err != nil {
						return nil, fmt.Errorf("Cannot retrieve columns: %v\n", err)
					}

					//perform basic validation
					err = ValidateCellsInRow(r.Cells, sheetCols, opt)
					if err != nil {
						return nil, err
					}
				}

				r.Cells[j].ColumnID = sheetCols[j].ID
			}
		}

		//adjust row options
		switch rowOpt {
		case ToBottom:
			r.ToBottom = true
		case ToTop:
			r.ToTop = true
		default:
			return nil, fmt.Errorf("Specified row option not yet implemented: %v", rowOpt)
		}
	}

	body, err := c.PostObject(fmt.Sprintf("sheets/%v/rows", sheetID), rows)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func encodeData(data interface{}) (io.Reader, error) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode: %v", err)
	}

	log.Printf("Body:\n%v\n", string(b.Bytes()))

	return b, nil
}

//PostObject will post data as JSOn
func (c *Client) PostObject(path string, data interface{}) (io.ReadCloser, error) {

	b, err := encodeData(data)
	if err != nil {
		return c.Post(path, b)
	}

	return nil, err
}

//Post will send a POST request through the client
func (c *Client) Post(path string, body io.Reader) (io.ReadCloser, error) {
	return c.send("POST", path, body)
}

//PutObject will post data as JSOn
func (c *Client) PutObject(path string, data interface{}) (io.ReadCloser, error) {

	b, err := encodeData(data)
	if err != nil {
		return c.Put(path, b)
	}

	return nil, err
}

//Put will send a PUT request through the client
func (c *Client) Put(path string, body io.Reader) (io.ReadCloser, error) {
	return c.send("PUT", path, body)
}

//Delete will send a DELETE request through the client
func (c *Client) Delete(path string) (io.ReadCloser, error) {
	return c.send("DELETE", path, nil)
}

//Get will append the proper info to pull from the API
func (c *Client) Get(path string) (io.ReadCloser, error) {
	return c.send("GET", path, nil)
}

func (c *Client) send(verb string, path string, body io.Reader) (io.ReadCloser, error) {
	req, err := http.NewRequest(verb, c.url+"/"+path, body)

	if err != nil {
		return nil, fmt.Errorf("Failed to create %v request: %v", verb, err)
	}

	log.Printf("URL: %v\n", req.URL)

	req.Header.Add("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to %v: %v", verb, err)
	}

	//TODO: check resp.StatusCode?
	return resp.Body, nil
}