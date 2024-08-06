package hub

import (
	"bytes"
	"html/template"
	"log"
	"sync"
	"github.com/jgrove2/browser_game_engine/canvas"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	sync.RWMutex

	clients map[*Client]bool

	broadcast  chan *canvas.BaseCanvas
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    map[*Client]bool{},
		broadcast:  make(chan *canvas.BaseCanvas),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

/*
The select statement inside the loop is used to attempt a non-blocking send operation. If the client's send channel is ready to receive the message (i.e., it's not blocked), case client.send <- msg: will execute, sending the message to the client.

If the client's send channel is not ready to receive the message (i.e., it's blocked because the client is not ready to receive data or the channel is full), the default: case will execute. This will close the client's send channel and remove the client from the h.clients map, effectively disconnecting the client.

This code is a common pattern in Go for handling multiple clients and ensuring that if one client is slow or unresponsive, it doesn't block the entire system from sending messages to other clients.
*/
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.Lock()
			h.clients[client] = true
			h.Unlock()

			log.Printf("client registered %s", client.id)

			client.send <- getMessageTemplate(&client.userCanvas.Canvas)
		case client := <-h.unregister:
			h.Lock()
			if _, ok := h.clients[client]; ok {
				close(client.send)
				log.Printf("client unregistered %s", client.id)
				delete(h.clients, client)
			}
			h.Unlock()
		case updatedCanvas := <-h.broadcast:
			h.RLock()
			for client := range h.clients {
				select {
				case client.send <- getMessageTemplate(updatedCanvas):
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.RUnlock()
		}
	}
}

func getMessageTemplate(data *canvas.BaseCanvas) []byte {
	tmpl, err := template.ParseFiles("templates/canvas.html")
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}

	// Render the template with the message as data.
	var renderedCanvas bytes.Buffer
	err = tmpl.Execute(&renderedCanvas, data)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}

	return renderedCanvas.Bytes()
}
