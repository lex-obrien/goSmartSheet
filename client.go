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
		log.Fatalf("Failed: %v\n", err)
		return s, err
	}
	defer body.Close()

	dec := json.NewDecoder(body)
	if err := dec.Decode(&s); err != nil {
		log.Fatalf("Failed to decode: %v\n", err)
		return s, err
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
	//TODO: need generic handling and ability to read from pages to get all datc... eventually
	dec := json.NewDecoder(body)
	if err = dec.Decode(&resp); err != nil {
		log.Fatalf("Failed to decode: %v\n", err)
		return
	}

	if err = json.Unmarshal(resp.Data, &cols); err != nil {
		log.Fatalf("Failed to decode data: %v\n", err)
		return
	}

	return
}

//GetJSONString with return a Json string of the result
func (c *Client) GetJSONString(path string, prettify bool) (string, error) {
	body, err := c.Get(path)
	if err != nil {
		log.Fatalf("Failed: %v\n", err)
		return "", err
	}
	defer body.Close()

	buf := new(bytes.Buffer)
	var s string

	if prettify {
		var m json.RawMessage

		dec := json.NewDecoder(body)
		if err := dec.Decode(&m); err != nil {
			log.Fatalf("Failed to decode: %v\n", err)
			return "", err
		}

		b, err := json.MarshalIndent(&m, "", "\t")
		if err != nil {
			log.Fatalf("Error during indent: %v\n", err)
			return "", err
		}

		s = string(b)
	} else {
		buf.ReadFrom(body)
		s = buf.String()
	}

	return s, nil
}

//Get will append the proper info to pull from the API
func (c *Client) Get(path string) (io.ReadCloser, error) {

	req, err := http.NewRequest("GET", c.url+"/"+path, nil)
	if err != nil {
		log.Fatalln("Failed: %v", err)
		return nil, err
	}

	log.Printf("URL: %v\n", req.URL)

	req.Header.Add("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatalln("Failed: %v", err)
		return nil, err
	}

	return resp.Body, nil
}

func (c *Client) AddRowToSheet(sheetID string, opt RowPostOptions, cellValues ...CellValue) (io.ReadCloser, error) {

	//TODO: validate length of cols with cells, match types, etc
	//right now this assumes the consumer is putting them in the correct order
	var r Row

	cols, err := c.GetColumns(sheetID)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	//TODO: apply multi-row validation logic here as well

	if len(cols) != len(cellValues) {
		log.Fatalf("Data Value length must match columns in sheet\n")
		return nil, nil //TODO: make an actual error
	}

	for i, col := range cols {
		c := Cell{ColumnID: col.ID}
		c.Value = &cellValues[i]
		r.Cells = append(r.Cells, c)
	}

	switch opt {
	case ToBottom:
		r.ToBottom = true
	case ToTop:
		r.ToTop = true
	case Above:
		log.Fatal("Above not implemented yet")
	}

	body, err := c.PostSingleObject(fmt.Sprintf("sheets/%v/rows", sheetID), r)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	return body, nil
}

func (c *Client) AddRowsToSheet(sheetID string, rowOpt RowPostOptions, rows []Row, opt PostOptions) (io.ReadCloser, error) {

	//TODO: validate length of cols with cells, match types, etc
	//right now this assumes the consumer is putting them in the correct order

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
						log.Fatalf("Error: %v\n", err)
						return nil, err //TODO: make an actual error
					}

					//perform basic validation
					err = ValidateCellsInRow(r.Cells, sheetCols, opt)
					if err != nil {
						log.Fatalf("Error: %v\n", err)
						return nil, err //TODO: make an actual error
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
		case Above:
			log.Fatal("Above not implemented yet")
		}
	}

	body, err := c.PostSingleObject(fmt.Sprintf("sheets/%v/rows", sheetID), rows)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	return body, nil
}

//Pending deletion
// func (c *Client) PostObjects(path string, jsonWrapper string, data ...interface{}) (io.ReadCloser, error) {

// 	//build data array
// 	log.Println("post")
// 	jsonData := make([]interface{}, len(data))
// 	for i, d := range data {
// 		b := new(bytes.Buffer)
// 		err := json.NewEncoder(b).Encode(d)
// 		if err != nil {
// 			log.Fatalf("Failed to encode: %v\n", err)
// 		}
// 		jsonData[i] = b.String()
// 	}
// 	//apply format to wrapper
// 	bodyData := new(bytes.Buffer)
// 	fmt.Fprintf(bodyData, jsonWrapper, jsonData...)

// 	log.Printf("Body:\n%v\n", string(bodyData.Bytes()))

// 	return c.Post(path, bodyData)
// }

func (c *Client) PostSingleObject(path string, data interface{}) (io.ReadCloser, error) {

	//build data array
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(data)
	if err != nil {
		log.Fatalf("Failed to encode: %v\n", err)
	}

	log.Printf("Body:\n%v\n", string(b.Bytes()))

	return c.Post(path, b)
}

func (c *Client) Post(path string, body io.Reader) (io.ReadCloser, error) {

	req, err := http.NewRequest("POST", c.url+"/"+path, body)

	if err != nil {
		log.Fatalln("Failed: %v", err)
		return nil, err
	}

	log.Printf("URL: %v\n", req.URL)

	req.Header.Add("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatalln("Failed: %v", err)
		return nil, err
	}

	//TODO: check resp.StatusCode?
	return resp.Body, nil
}

//TODO: consolidate PUT and post
func (c *Client) Put(path string, body io.Reader) (io.ReadCloser, error) {

	req, err := http.NewRequest("PUT", c.url+"/"+path, body)

	if err != nil {
		log.Fatalln("Failed: %v", err)
		return nil, err
	}

	log.Printf("URL: %v\n", req.URL)

	req.Header.Add("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatalln("Failed: %v", err)
		return nil, err
	}

	//TODO: check resp.StatusCode?
	return resp.Body, nil
}
