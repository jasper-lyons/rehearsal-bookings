package handlers

import (
	"encoding/json"
	"bytes"
	"net/http"
	"io"
)

type Api struct {
	baseUrl string
	headers map[string]string
}

func NewApi(baseUrl string, headers map[string]string) Api {
	return Api { baseUrl: baseUrl, headers: headers }
}

func (a* Api) Request(method string, path string, body interface{}) (string, error) {
		requestBody, err := json.Marshal(body)
		if err != nil {
			return "", err
		}
		bodyReader := bytes.NewReader(requestBody)

		req, err := http.NewRequest(method, a.baseUrl + path, bodyReader)	
		if err != nil {
			return "", err
		}

		for k, v := range a.headers {
			req.Header.Set(k, v)
		}

		client := http.DefaultClient
		response, err := client.Do(req)
		if err != nil {
			return "", err
		}

		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			return "", err
		}

		return string(responseBody), nil
}

func (a* Api) Post(path string, body interface{}) (string, error) {
	return a.Request(http.MethodPost, path, body) 
}
