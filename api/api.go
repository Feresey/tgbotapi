package api

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	OK          bool                `json:"ok"`
	Result      json.RawMessage     `json:"result"`
	Description *string             `json:"description,omitempty"`
	Parameters  *ResponseParameters `json:"parameters,omitempty"`
}

type API struct {
	cli    *http.Client
	apiURL string
}
