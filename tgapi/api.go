package tgapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const APIEndpoint = "https://api.telegram.org"

type Response struct {
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

func NewWithEndpointAndClient(token, endpoint string, cli *http.Client) *API {
	return &API{
		cli:      cli,
		endpoint: fmt.Sprintf("%s/bot%s", endpoint, token),
	}
}

func New(token string) *API {
	return NewWithEndpointAndClient(token, APIEndpoint, http.DefaultClient)
}

// MakeRequest makes a request to a specific endpoint with our token.
func (api *API) MakeRequest(method string, data interface{}) (*Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(fmt.Sprintf("%s/%s", api.endpoint, method))
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodPost,
		Body:   ioutil.NopCloser(bytes.NewReader(body)),
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		URL: u,
		// context?
	}
	resp, err := api.cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp Response
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
