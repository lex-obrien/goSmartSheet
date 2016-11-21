package goSmartSheet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/pkg/errors"
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
		return s, errors.Wrapf(err, "Failed to get sheet (ID: %v)", id)

	}
	defer body.Close()

	dec := json.NewDecoder(body)
	if err := dec.Decode(&s); err != nil {
		return s, errors.Wrap(err, "Failed to decode")
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
		return nil, errors.Wrap(err, "Failed to decode response")
	}

	if err = json.Unmarshal(resp.Data, &cols); err != nil {
		return nil, errors.Wrap(err, "Failed to decode columns")
	}

	return
}

//GetJSONString with return a Json string of the result
func (c *Client) GetJSONString(path string, prettify bool) (string, error) {
	body, err := c.Get(path)
	if err != nil {
		return "", errors.Wrap(err, "Failed to Get JSON String")
	}
	defer body.Close()

	buf := new(bytes.Buffer)
	var s string

	if prettify {
		var m json.RawMessage

		dec := json.NewDecoder(body)
		if err := dec.Decode(&m); err != nil {
			return "", errors.Wrap(err, "Failed to decode")
		}

		b, err := json.MarshalIndent(&m, "", "\t")
		if err != nil {
			return "", errors.Wrap(err, "Error during indent")
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
						return nil, errors.Wrapf(err, "Cannot retrieve columns: %v")
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
			return nil, errors.Errorf("Specified row option not yet implemented: %v", rowOpt)
		}
	}

	body, err := c.PostObject(fmt.Sprintf("sheets/%v/rows", sheetID), rows)
	if err != nil {
		return nil, err
	}

	return body, nil
}

//DeleteRowsFromSheet will delte the specified rowes from the specified sheet
func (c *Client) DeleteRowsFromSheet(sheetID string, rows []Row) (io.ReadCloser, error) {
	ids := []string{}
	for _, r := range rows {
		ids = append(ids, strconv.FormatInt(r.ID, 10))
	}

	return c.DeleteRowsIdsFromSheet(sheetID, ids)
}

//DeleteRowsIdsFromSheet will delete the specified rowIDs from the specified sheet
func (c *Client) DeleteRowsIdsFromSheet(sheetID string, ids []string) (io.ReadCloser, error) {
	path := fmt.Sprintf("sheets/%v/rows?ids=%v", sheetID, strings.Join(ids, ","))
	return c.Delete(path)
}

//UpdateRowsOnSheet will update the specified rows and data
func (c *Client) UpdateRowsOnSheet(sheetID string, rows []Row) (io.ReadCloser, error) {
	//clean the row to remove the data that cannot be sent accross

	// //the caller needs to pass in clean data right now
	// newRows := []Row{}
	// copy(newRows, rows)

	// for i := range newRows {
	// 	newRows[i].CreatedAt = nil
	// 	newRows[i].ModifiedAt = nil
	// 	newRows[i]. = nil
	// }

	return c.PutObject(fmt.Sprintf("sheets/%v/rows", sheetID), rows)
}

func encodeData(data interface{}) (io.Reader, error) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(data)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to encode")
	}

	log.Printf("Body:\n%v\n", string(b.Bytes()))

	return b, nil
}

//PostObject will post data as JSOn
func (c *Client) PostObject(path string, data interface{}) (io.ReadCloser, error) {

	b, err := encodeData(data)
	if err != nil {
		return nil, err
	}

	return c.Post(path, b)
}

//Post will send a POST request through the client
func (c *Client) Post(path string, body io.Reader) (io.ReadCloser, error) {
	return c.send("POST", path, body)
}

//PutObject will post data as JSOn
func (c *Client) PutObject(path string, data interface{}) (io.ReadCloser, error) {

	b, err := encodeData(data)
	if err != nil {
		return nil, err
	}
	return c.Put(path, b)
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

func (c *Client) send(verb string, p string, body io.Reader) (io.ReadCloser, error) {
	req, err := http.NewRequest(verb, path.Join(c.url, p), body)

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create %v request", verb)
	}

	log.Printf("URL: %v\n", req.URL)

	req.Header.Add("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to %v", verb)
	}

	//TODO: check resp.StatusCode?
	return resp.Body, nil
}
