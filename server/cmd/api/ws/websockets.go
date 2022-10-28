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

var globalWsConns = map[string]map[int]WebsocketConn{}

// create a new websocket connection between client and server
func (c EchoWrapper) NewWsConn(key string, returns chan *WebsocketConn) {
	ConnId := len(globalWsConns[key])
	if ConnId == 0 {
		globalWsConns[key] = map[int]WebsocketConn{}
	}

	connHandler := make(chan *websocket.Conn)
	go func() {
		wsConn := <-connHandler

		wst := WebsocketConn{key, wsConn}
		if _, exist := globalWsConns[key][ConnId]; exist {
			for i := ConnId; ; i++ {
				if _, exist := globalWsConns[key][i]; !exist {
					ConnId = i
					break
				}
			}
		}
		log.Printf("Conn Id #[%s_%d] connected\n", key, ConnId)
		globalWsConns[key][ConnId] = wst

		returns <- &wst
		close(returns)
	}()

	websocket.Handler(func(ws *websocket.Conn) {
		connHandler <- ws

		defer func() {
			delete(globalWsConns[key], ConnId) // remove conn from global
			ws.Close()                         // Close WS Connection

			log.Printf("Conn Id #[%s_%d] disconnected\n", key, ConnId)
		}()

		clientOK := make(chan bool)
		interrupt := make(chan any)
		go Ping(ws, clientOK, interrupt)

		for {
			select {
			case <-interrupt:
				return // Client stopped to respond -> quit connection
			default:
				// Read
				var msg string
				if err := websocket.Message.Receive(ws, &msg); err == nil {
					clientOK <- true
				}
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
}

func Ping(ws *websocket.Conn, ok chan bool, interrupt chan any) {
	for {
		err := websocket.Message.Send(ws, "PING")
		if err != nil {
			close(interrupt)
			return
		}

		select {
		case <-time.After(5 * time.Second): // Wait 5s for the client to answer
			close(interrupt) // Client didn't answer, closing ws conn
			return
		case <-ok:
			// time.Sleep(time.Minute) // Client answered, backing off 1min
			time.Sleep(10 * time.Second) // Client answered, backing off 1min
		}
	}
}

func GetWsConn(key string) (map[int]WebsocketConn, error) {
	if wst, exist := globalWsConns[key]; exist {
		return wst, nil
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
