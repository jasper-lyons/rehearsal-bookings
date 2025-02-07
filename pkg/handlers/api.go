package handlers

import (
	"encoding/json"
	"bytes"
	"net/http"
	"io"
	"fmt"
)

type Api struct {
	baseUrl string
	headers map[string]string
}

type ApiResponse struct {
	Body string
	Status int
}

func NewApi(baseUrl string, headers map[string]string) Api {
	return Api { baseUrl: baseUrl, headers: headers }
}

func (a* Api) Request(method string, path string, body interface{}) (*ApiResponse, error) {
		requestBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader := bytes.NewReader(requestBody)

		req, err := http.NewRequest(method, a.baseUrl + path, bodyReader)	
		if err != nil {
			return nil, err
		}

		for k, v := range a.headers {
			req.Header.Set(k, v)
		}

		client := http.DefaultClient
		response, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		return &ApiResponse {
			Body: string(responseBody),
			Status: response.StatusCode,
		}, nil
}

func (a* Api) Post(path string, body interface{}) (*ApiResponse, error) {
	return a.Request(http.MethodPost, path, body) 
}

func (a* Api) Get(path string, args... any) (*ApiResponse, error) {
	return a.Request(http.MethodGet, fmt.Sprintf(path, args...), nil)
}

func (a *Api) Delete(path string, args... any) (*ApiResponse, error) {
	return a.Request(http.MethodDelete, fmt.Sprintf(path, args...), nil)
}
