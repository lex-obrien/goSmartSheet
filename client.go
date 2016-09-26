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

//SmartSheetAPI is used to interact with the SS API
type Client struct {
	url    string
	apiKey string
	client *http.Client
}

func GetClient(apiKey string) *Client {
	//default to prod API
	api := &Client{url: "https://api.smartsheet.com/2.0", apiKey: apiKey}
	api.client = &http.Client{} //per docs clientss should be made once, https://golang.org/pkg/net/http/

	return api
}

func (c *Client) GetSheetFilterCols(id string, onlyTheseColumns []string) (Sheet, error) {
	filter := "columnIds=" + strings.Join(onlyTheseColumns, ",")
	return c.GetSheet(id, filter)
}

///GetSheet returns a sheet with the specified Id
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

func (c *Client) AddRowToSheet(sheetId string, cellValues ...CellValue) (io.ReadCloser, error) {

	//TODO: validate length of cols with cells, match types, etc
	//right now this assumes the consumer is putting them in the correct order

	var cells []Cell

	cols, err := c.GetColumns(sheetId)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	if len(cols) != len(cellValues) {
		log.Fatalf("Data Value length must match columsn in sheet\n")
		return nil, nil //TODO: make an actual error
	}

	for i, col := range cols {
		c := Cell{ColumnID: col.ID}
		c.Value = &cellValues[i]
		cells = append(cells, c)
	}

	//TODO: make this just use a real row...
	//turns out their API made more sense than I thought... this is just a row, nothing special, probably dont need my PostObjs method...
	body, err := c.PostObjects("sheets/597019279550340/rows", `[{"toBottom":true, "cells": %v }]`, cells)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	return body, nil
}

func (c *Client) PostObjects(path string, jsonWrapper string, data ...interface{}) (io.ReadCloser, error) {

	//build data array
	log.Println("post")
	jsonData := make([]interface{}, len(data))
	for i, d := range data {
		b := new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(d)
		if err != nil {
			log.Fatalf("Failed to encode: %v\n", err)
		}
		jsonData[i] = b.String()
	}
	//apply format to wrapper
	bodyData := new(bytes.Buffer)
	fmt.Fprintf(bodyData, jsonWrapper, jsonData...)

	log.Printf("Body:\n%v\n", string(bodyData.Bytes()))

	return c.Post(path, bodyData)
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
