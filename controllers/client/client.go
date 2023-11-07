package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"stoicdynamics.com/relax/types"
)

var (
	requestsHistory []types.RequestLog = []types.RequestLog{}
)

func MakeRequest(requestParams types.Request) (types.Response, error) {
	bodyReader := bytes.NewReader([]byte(requestParams.Body))
	request, err := http.NewRequest(string(requestParams.Verb), requestParams.Url, bodyReader)

	headers := strings.Split(requestParams.Headers, "\n")

	for _, header := range headers {
		parts := strings.Split(header, ": ")
		if len(parts) > 1 {
			request.Header.Add(parts[0], parts[1])
		}
	}

	if err != nil {
		return types.Response{}, errors.New(fmt.Sprintf("Error parsing HTTP request: %v\n", err))
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return types.Response{}, errors.New(fmt.Sprintf("Error executing request: %v\n", err))
	}
	defer response.Body.Close()

	// Read and print the response
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return types.Response{}, errors.New(fmt.Sprintf("Error reading response: %v\n", err))
	}

	responseHeader := ""

	for key, values := range response.Header {
		for _, value := range values {
			responseHeader += fmt.Sprintf("%s: %s\n", key, value)
		}
	}

	return types.Response{
		Status: response.Status,
		Header: responseHeader,
		Body:   string(responseBody),
	}, nil
}
