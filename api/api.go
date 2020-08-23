package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const APIEndpoint = "https://api.telegram.org"

type APIResponse struct {
	Ok          bool               `json:"ok"`
	Result      json.RawMessage    `json:"result"`
	Description string             `json:"description,omitempty"`
	Parameters  ResponseParameters `json:"parameters,omitempty"`
	ErrorCode   int                `json:"error_code,omitempty"`
}

type Error struct {
	Code    int
	Message string
	ResponseParameters
}

func (e Error) Error() string {
	return fmt.Sprintf(
		"response from api. code: %d, message: %s, parameters: %v",
		e.Code, e.Message, e.ResponseParameters)
}

type API struct {
	cli      *http.Client
	endpoint string
}

func New(token string) *API {
	return &API{
		cli:      http.DefaultClient,
		endpoint: fmt.Sprintf("%s/bot%s", APIEndpoint, token),
	}
}

// MakeRequest makes a request to a specific endpoint with our token.
func (api *API) MakeRequest(method string, data interface{}) (*APIResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := api.cli.Post(
		fmt.Sprintf("%s/%s", api.endpoint, method),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return &apiResp, err
	}

	if !apiResp.Ok {
		return nil, Error{
			Code:               apiResp.ErrorCode,
			Message:            apiResp.Description,
			ResponseParameters: apiResp.Parameters,
		}
	}

	return &apiResp, nil
}
