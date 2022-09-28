package routes

import (
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

func (c EchoWrapper) HandleResp(code uint, data ...any) error {
	resp := ResponsePayload{Code: int(code)}
	if code >= 200 && code < 400 {
		resp.Success = true
	}
	if len(data) == 1 {
		resp.Data = data[0]
	}

	return c.JSON(resp.Code, resp)
}