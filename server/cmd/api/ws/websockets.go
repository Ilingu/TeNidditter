package ws

import (
	"errors"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type EchoWrapper struct {
	echo.Context
}

type WebsocketConn struct {
	Key    string
	WsConn *websocket.Conn
	// callback func(msg any)
}

var globalWsConns map[string][]WebsocketConn

// create a new websocket connection between client and server
func (c *EchoWrapper) NewWsConn(key string) *WebsocketConn {
	ch := make(chan *websocket.Conn)
	go func() {
		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()
			ws.SetDeadline(time.Now().Add(2 * time.Hour))

			ch <- ws
			for {
				// Read
				var msg any
				if err := websocket.Message.Receive(ws, &msg); err == nil {
					// cb(msg)
					log.Println(msg)
				}
			}
		}).ServeHTTP((*c).Response(), (*c).Request())
	}()

	wsConn := <-ch

	wst := WebsocketConn{key, wsConn}
	globalWsConns[key] = append(globalWsConns[key], wst)

	return &wst
}

func GetWsConn(key string) (*[]WebsocketConn, error) {
	if wst, exist := globalWsConns[key]; exist {
		return &wst, nil
	}
	return nil, errors.New("websocket connection not found")
}

func (wst *WebsocketConn) CloseConn() bool {
	if err := wst.WsConn.Close(); err != nil {
		return false
	}

	delete(globalWsConns, wst.Key)
	return true
}
