package tgapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
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
	endpoint     string
	fileEndpoint string
}

func NewWithEndpointAndClient(token, endpoint, fileEndpoint string, cli *http.Client) *API {
	return &API{
		cli:          cli,
		endpoint:     fmt.Sprintf("%s/bot%s", endpoint, token),
		fileEndpoint: fmt.Sprintf("%s/bot%s", fileEndpoint, token),
	}
}

func New(token string) *API {
	return NewWithEndpointAndClient(token, APIEndpoint, FileEndpoint, http.DefaultClient)
}

func decodeAPIResponse(r io.Reader) (*Response, error) {
	var apiResp Response
	if err := json.NewDecoder(r).Decode(&apiResp); err != nil {
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

	resp, err := api.cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return decodeAPIResponse(resp.Body)
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

	resp, err := api.cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return decodeAPIResponse(resp.Body)
}

// func (tbg BaseRequester) Post(l *zap.SugaredLogger, token string, method string, params url.Values, data map[string]PostFile) (json.RawMessage, error) {
// 	endpoint := tbg.ApiUrl
// 	if endpoint == "" {
// 		endpoint = ApiUrl
// 	}

// 	b := bytes.Buffer{}
// 	w := multipart.NewWriter(&b)
// 	defer w.Close()

// 	for field, x := range data {
// 		fileName := x.FileName
// 		if fileName == "" {
// 			fileName = "unnamed_file"
// 		}

// 		part, err := w.CreateFormFile(field, fileName)
// 		if err != nil {
// 			return nil, err
// 		}

// 		_, err = io.Copy(part, x.File)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	err := w.Close()
// 	if err != nil {
// 		return nil, err
// 	}

// 	req, err := http.NewRequest("POST", endpoint+token+"/"+method, &b)
// 	if err != nil {
// 		l.Debugw("failed to execute POST request",
// 			zapcore.Field{
// 				Key:    "method",
// 				Type:   zapcore.StringType,
// 				String: method,
// 			},
// 			zap.Error(err))
// 		return nil, errors.Wrapf(err, "client error executing POST request to %v", method)
// 	}
// 	req.URL.RawQuery = params.Encode()
// 	req.Header.Set("Content-Type", w.FormDataContentType())

// 	l.Debugf("POST request with body: %+v", b)
// 	l.Debugf("executing POST: %+v", req)
// 	resp, err := tbg.Client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	l.Debugf("successful POST request: %+v", resp)

// 	var r response
// 	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
// 		l.Debugw("failed to deserialize POST response body",
// 			zapcore.Field{
// 				Key:    "method",
// 				Type:   zapcore.StringType,
// 				String: method,
// 			},
// 			zap.Error(err))
// 		return nil, errors.Wrapf(err, "decoding error while reading POST request to %v", method)
// 	}
// 	if !r.Ok {
// 		l.Debugw("error from POST",
// 			zapcore.Field{
// 				Key:    "method",
// 				Type:   zapcore.StringType,
// 				String: method,
// 			},
// 			zapcore.Field{
// 				Key:     "error_code",
// 				Type:    zapcore.Int64Type,
// 				Integer: int64(r.ErrorCode),
// 			},
// 			zapcore.Field{
// 				Key:    "description",
// 				Type:   zapcore.StringType,
// 				String: r.Description,
// 			},
// 		)
// 		return nil, &TelegramError{
// 			Method:      method,
// 			Values:      params,
// 			Code:        r.ErrorCode,
// 			Description: r.Description,
// 		}
// 	}

// 	l.Debugw("obtained POST result",
// 		zapcore.Field{
// 			Key:    "method",
// 			Type:   zapcore.StringType,
// 			String: method,
// 		},
// 		zapcore.Field{
// 			Key:    "result",
// 			Type:   zapcore.StringType,
// 			String: string(r.Result),
// 		},
// 	)

// 	return r.Result, nil
// }
