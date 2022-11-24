package sse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"teniditter-server/cmd/api/routes"
	"teniditter-server/cmd/global/console"
	"teniditter-server/cmd/global/utils"
)

type Client struct {
	IP     string
	Events chan any
	Done   chan bool
}

var clients = map[string]*Client{}

func SSEHandler(w http.ResponseWriter, r *http.Request) {
	ip, err := routes.GetIP(r, true)
	if err != nil {
		fmt.Fprintf(w, "data: CLOSING\n\n")
		return
	}

	client := &Client{IP: ip, Events: make(chan any), Done: make(chan bool)}
	client.addToGlobal()

	defer func() {
		close(client.Events)
		close(client.Done)
		client.RemoveFromGlobal()
	}()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		select {
		case ev := <-client.Events:
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(ev)
			if err != nil {
				continue
			}

			_, err = fmt.Fprintf(w, "data: %v\n\n", buf.String())
			if err != nil {
				console.Log(fmt.Sprintf("Failed to sent SSE to %s", r.RemoteAddr), console.Error)
				fmt.Fprintf(w, "data: CLOSING\n\n")
				return
			}

			if f, ok := w.(http.Flusher); ok {
				log.Println("flushing", r.RemoteAddr)
				f.Flush()
			}
		case close := <-client.Done:
			if close {
				fmt.Fprintf(w, "data: CLOSING\n\n")
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
