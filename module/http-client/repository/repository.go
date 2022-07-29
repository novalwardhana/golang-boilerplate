package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/novalwardhana/golang-boilerplate/config/env"
	"github.com/novalwardhana/golang-boilerplate/module/http-client/model"
)

type repository struct {
}

type Repository interface {
	Create(payload *model.Person) <-chan model.Result
	GetData(page, limit int) <-chan model.Result
	BulkInsert(filedir, filename string) <-chan model.Result
	DownloadCSV() <-chan model.Result
}

func NewRepository() Repository {
	return &repository{}
}

// CrudCreate:
func (r *repository) Create(payload *model.Person) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Prepare payload to byte */
		payloadByte, err := json.Marshal(payload)
		if err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Prepare http client */
		httpClient := http.Client{
			Timeout: 5 * time.Minute,
		}

		/* Prepare http request */
		httpHeader := http.Header{}
		httpHeader.Add("Content-Type", "application/json")
		httpRequest := http.Request{}
		httpRequest.Header = httpHeader
		httpRequest.URL, _ = url.Parse(fmt.Sprintf("%s/%s/%s", os.Getenv(env.EnvHTTPClientURL), "crud", "create"))
		httpRequest.Method = "POST"
		httpRequest.Body = ioutil.NopCloser(bytes.NewBuffer(payloadByte))

		/* Process */
		httpResponse, err := httpClient.Do(&httpRequest)
		if err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Check response */
		if httpResponse.StatusCode != 200 {
			result <- model.Result{Error: errors.New("Failed process create new data")}
			return
		}
		httpResponseBody, err := ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Check response body */
		var response = new(model.Response)
		if err := json.Unmarshal(httpResponseBody, response); err != nil {
			result <- model.Result{Error: err}
			return
		}
		if response.Status != 200 {
			result <- model.Result{Error: errors.New(response.Message)}
			return
		}

		result <- model.Result{}
	}()
	return result
}

// GetData:
func (r *repository) GetData(page, limit int) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Prepare http client */
		httpClient := http.Client{
			Timeout: 50 * time.Second,
		}

		/* Prepare http request */
		httpHeader := http.Header{}
		httpHeader.Add("Content-Type", "application/json")
		httpRequest := http.Request{}
		httpRequest.Header = httpHeader
		httpRequest.URL, _ = url.Parse(fmt.Sprintf("%s/%s/%s?page=%d&limit=%d", os.Getenv(env.EnvHTTPClientURL), "crud", "get-data", page, limit))
		httpRequest.Method = "GET"
		httpRequest.Body = nil

		/* Process*/
		httpResponse, err := httpClient.Do(&httpRequest)
		if err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Check response */
		if httpResponse.StatusCode != 200 {
			result <- model.Result{Error: errors.New("Failed get data")}
			return
		}
		httpResponseBody, err := ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Check response body */
		var response = new(model.Response)
		if err := json.Unmarshal(httpResponseBody, response); err != nil {
			result <- model.Result{Error: err}
			return
		}
		if response.Status != 200 {
			result <- model.Result{Error: errors.New("Failed get data")}
			return
		}

		result <- model.Result{Data: response.Data}
	}()
	return result
}

// BulkInsert:
func (r *repository) BulkInsert(filedir, filename string) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Open file */
		file, err := os.Open(filepath.Join(filedir, filename))
		defer file.Close()
		if err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Prepare http client */
		httpClient := http.Client{
			Timeout: 10 * time.Second,
		}

		/* Prepare http headers */
		httpHeader := http.Header{}
		httpHeader.Add("Content-Type", "application/json")
		httpRequest := http.Request{}
		httpRequest.Header = httpHeader
		httpRequest.URL, _ = url.Parse(fmt.Sprintf("%s/%s/%s", os.Getenv(env.EnvHTTPClientURL), "advance-crud", "bulk-insert"))
		httpRequest.Method = "POST"

		/* Create byte buffer */
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		fieldFile, err := writer.CreateFormFile("file", file.Name())
		if err != nil {
			result <- model.Result{Error: err}
			return
		}
		_, err = io.Copy(fieldFile, file)
		if err != nil {
			result <- model.Result{Error: err}
			return
		}
		httpRequest.Body = ioutil.NopCloser(body)

		/* Process */
		httpResponse, err := httpClient.Do(&httpRequest)
		if err != nil {
			result <- model.Result{Error: err}
			return
		}
		if httpResponse.StatusCode != 200 {
			result <- model.Result{Error: errors.New("Failed process bulk insert")}
			return
		}

		/* Check response body */
		httpResponseBody, err := ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Parse response body */
		var response = new(model.Response)
		if err := json.Unmarshal(httpResponseBody, response); err != nil {
			result <- model.Result{Error: err}
			return
		}

		result <- model.Result{}
	}()
	return result
}

// ExportCSV:
func (r *repository) DownloadCSV() <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* http client */
		httpClient := http.Client{
			Timeout: 10 * time.Second,
		}

		/* http header */
		httpHeader := http.Header{}
		//httpHeader.Add("Content-Type", "application/json")

		/* http request */
		httpRequest := http.Request{}
		httpRequest.Header = httpHeader
		httpRequest.URL, _ = url.Parse(fmt.Sprintf("%s/%s/%s", os.Getenv(env.EnvHTTPClientURL), "advance-crud", "download-csv"))
		httpRequest.Method = "GET"

		/* Process */
		httpResponse, err := httpClient.Do(&httpRequest)
		if err != nil {
			result <- model.Result{Error: err}
			return
		}
		if httpResponse.StatusCode != 200 {
			result <- model.Result{Error: errors.New("Failed export csv")}
			return
		}

		/* Get response body */
		httpResponseBody, err := ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			result <- model.Result{Error: err}
			return
		}

		result <- model.Result{Data: httpResponseBody}

	}()
	return result
}
