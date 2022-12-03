package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

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

func (c EchoWrapper) InjectJsonHeader(headerName string, data any) {
	stringifiedSubs, err := json.Marshal(data)
	if err != nil {
		stringifiedSubs = []byte{}
	}

	c.Response().Header().Set(headerName, string(stringifiedSubs))
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

func GetIP(r *http.Request, onlyRemoteAddr bool) (string, error) {
	if !onlyRemoteAddr {
		//Get IP from the X-REAL-IP header
		ip := r.Header.Get("X-REAL-IP")
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}

		//Get IP from X-FORWARDED-FOR header
		ips := r.Header.Get("X-FORWARDED-FOR")
		splitIps := strings.Split(ips, ",")
		for _, ip := range splitIps {
			netIP := net.ParseIP(ip)
			if netIP != nil {
				return ip, nil
			}
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	if netIP := net.ParseIP(ip); netIP == nil {
		return "", errors.New("no valid ip found")
	}

	return ip, nil
}
