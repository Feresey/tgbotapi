package tgapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

const APIEndpoint = "https://api.telegram.org"
const FileEndpoint = "https://api.telegram.org/file"

type Response struct {
	OK          bool               `json:"ok"`
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
		"api response. code: %d, message: %s, parameters: %v",
		e.Code, e.Message, e.ResponseParameters)
}

type API struct {
	cli          *http.Client
	token        string
	endpoint     string
	fileEndpoint string
}

func NewWithEndpointAndClient(token, endpoint, fileEndpoint string, cli *http.Client) *API {
	return &API{
		cli:          cli,
		token:        token,
		endpoint:     fmt.Sprintf("%s/bot%s", endpoint, token),
		fileEndpoint: fmt.Sprintf("%s/bot%s", fileEndpoint, token),
	}
}

func New(token string) *API {
	return NewWithEndpointAndClient(token, APIEndpoint, FileEndpoint, http.DefaultClient)
}

func (api *API) decodeAPIResponse(req *http.Request) (*Response, error) {
	resp, err := api.cli.Do(req)
	if err != nil {
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			urlErr.URL = strings.Replace(urlErr.URL, api.token, "TOKEN", 1)
			return nil, urlErr
		}
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp Response
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return &apiResp, err
	}

	if !apiResp.OK {
		return &apiResp, Error{
			Code:               apiResp.ErrorCode,
			Message:            apiResp.Description,
			ResponseParameters: apiResp.Parameters,
		}
	}

	return &apiResp, nil
}

// MakeRequest makes a request to a specific endpoint with our token.
func (api *API) MakeRequest(ctx context.Context, method string, data interface{}) (*Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s", api.endpoint, method)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return api.decodeAPIResponse(req)
}

func (api *API) UploadFile(
	ctx context.Context,
	values url.Values,
	method string,
	filetype string,
	file *InputFile,
) (*Response, error) {
	b := new(bytes.Buffer)
	w := multipart.NewWriter(b)
	defer w.Close()

	wr, err := w.CreateFormFile(filetype, file.Name)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(wr, file.Reader)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/%s", api.endpoint, method),
		b,
	)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = values.Encode()
	req.Header.Set("Content-Type", w.FormDataContentType())

	return api.decodeAPIResponse(req)
}
