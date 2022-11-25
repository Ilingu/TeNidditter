package sse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"teniditter-server/cmd/api/routes"
	"teniditter-server/cmd/global/console"
	"teniditter-server/cmd/global/utils"
	utils_concurrency "teniditter-server/cmd/global/utils/concurrency"
)

type Client struct {
	IP     string
	Events chan any
	Close  chan bool

	// this field allows to keep track of the integrety of the connection, it represent the minimum numbers of events that will be sent through this connection
	//
	// note that the connection cannot be closed until ther sum of SentEvent and ErrorEvent is inferior than MinOpsNumber
	//
	// The minimum number of events should be at least 1
	MinOpsNumber chan uint64

	SentNumber  uint64
	ErrorNumber uint64
}

var clients = map[string]*Client{}

func SSEHandler(w http.ResponseWriter, r *http.Request) {
	ip, err := routes.GetIP(r, true)
	if err != nil {
		fmt.Fprintf(w, "data: CLOSING\n\n")
		return
	}

	client := &Client{IP: ip, Events: make(chan any), Close: make(chan bool), MinOpsNumber: make(chan uint64, 1)}
	client.addToGlobal()

	defer func() {
		close(client.Events)
		close(client.Close)
		client.RemoveFromGlobal()
	}()

	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ALLOWED_ORIGIN"))
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flush := func() {
		if f, ok := w.(http.Flusher); ok {
			log.Println("flushing", r.RemoteAddr)
			f.Flush()

			atomic.AddUint64(&client.SentNumber, 1)
		} else {
			atomic.AddUint64(&client.ErrorNumber, 1)
		}
	}

	minOpsNum := <-client.MinOpsNumber
	if minOpsNum < 1 {
		fmt.Fprintf(w, "data: CLOSING\n\n")
		flush()
		return
	}

	sentTracker := utils_concurrency.NewMultipleRoutineWaitGroup()
	sentTracker.Add(int(minOpsNum))

	// Wait events
	for {
		select {
		case ev := <-client.Events:
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(ev)
			if err != nil {
				atomic.AddUint64(&client.ErrorNumber, 1)
				sentTracker.Done()
				continue
			}

			if _, err = fmt.Fprintf(w, "data: %v\n\n", buf.String()); err != nil {
				console.Log(fmt.Sprintf("Failed to sent SSE to %s", r.RemoteAddr), console.Error)
				atomic.AddUint64(&client.ErrorNumber, 1)
				sentTracker.Done()
				continue
			}

			flush()
			sentTracker.Done()

		case close := <-client.Close:
			sentTracker.Wait()
			if close && client.SentNumber+client.ErrorNumber >= minOpsNum {
				fmt.Fprintf(w, "data: CLOSING\n\n")
				flush()
				return
			}
		}
	}
}

func GetClient(ip string) (*Client, bool) {
	client, ok := clients[ip]
	return client, ok
}

func (client *Client) addToGlobal() {
	log.Printf("New SSE Connection: [%s]\n", client.IP)

	clients[client.IP] = client
	go dispatchNewConnEvent(client.IP)
}
func (client *Client) RemoveFromGlobal() {
	log.Printf("SSE Connection Disconnected: [%s]\n", client.IP)
	delete(clients, client.IP)
}

var listeners = map[string]func(ip string){}

func ListenToNewConn(cb func(ip string)) string {
	id, _ := utils.GenerateRandomChars(4)
	listeners[id] = cb
	return id
}
func UnlistenNewConn(id string) {
	delete(listeners, id)
}
func dispatchNewConnEvent(ip string) {
	for _, cb := range listeners {
		go cb(ip)
	}
}
