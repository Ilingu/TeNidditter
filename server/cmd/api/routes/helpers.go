package routes

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResponsePayload struct {
	Success bool `json:"success"`
	Code    int  `json:"code"`
	Data    any  `json:"data"`
}

type EchoWrapper struct {
	echo.Context
}

func newRespPayload(code uint, data ...any) *ResponsePayload {
	resp := ResponsePayload{Code: int(code)}
	if code >= 200 && code < 400 {
		resp.Success = true
	}
	if len(data) == 1 {
		resp.Data = data[0]
	}
	return &resp
}

func (c EchoWrapper) HandleResp(code uint, data ...any) error {
	resp := newRespPayload(code, data...)
	return c.JSON(resp.Code, resp)
}

func (c EchoWrapper) HandleRespBlob(code uint, data ...any) error {
	resp := newRespPayload(code, data...)

	blob, err := json.Marshal(resp)
	if err != nil {
		return c.JSON(http.StatusCreated, resp)
	}

	return c.JSONBlob(resp.Code, blob)
}
