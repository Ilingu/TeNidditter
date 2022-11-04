package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResponsePayload struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Error   string `json:"error"`
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
		if resp.Success {
			resp.Data = data[0]
		} else if errorStr, ok := data[0].(string); ok {
			resp.Error = errorStr
		}
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

func (c EchoWrapper) InjectSubs(subs []string) {
	stringifiedSubs, err := json.Marshal(subs)
	if err != nil {
		stringifiedSubs = []byte{}
	}

	c.Response().Header().Set("TedditSubs", string(stringifiedSubs))
}

func (c EchoWrapper) SetPublicCache(maxage int) {
	if maxage == 0 {
		maxage = 1800 // 30min in seconds
	}
	c.Response().Header().Set("Cache-Control", fmt.Sprintf("public,max-age=%d", maxage))
}
func (c EchoWrapper) SetAuthCache(maxage ...int) {
	var ma = 129600 // 1.5d in seconds
	if len(maxage) == 1 {
		ma = maxage[0]
	}
	c.Response().Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", ma))
}
