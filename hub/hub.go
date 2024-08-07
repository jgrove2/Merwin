package hub

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"sync"

	"github.com/jgrove2/Merwin/canvas"
	"github.com/jgrove2/Merwin/window"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	sync.RWMutex

	clients map[*Client]bool

	window     *window.Window
	broadcast  chan *canvas.BaseCanvas
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    map[*Client]bool{},
		window:     window.NewWindow(),
		broadcast:  make(chan *canvas.BaseCanvas),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.Lock()
			h.clients[client] = true
			h.Unlock()
			slog.Info("New client connected", "client_id", client.id)

			client.send <- getMessageTemplate(&client.userCanvas.Canvas)
		case client := <-h.unregister:
			h.Lock()
			if _, ok := h.clients[client]; ok {
				close(client.send)
				slog.Info("client unregistered", "client_id", client.id)

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
		case userEvent := <-h.window.EventList:
			slog.Info(userEvent.UserID, userEvent.Event, userEvent.Key)
		}

	}
}

func getMessageTemplate(data *canvas.BaseCanvas) []byte {
	tmpl, err := template.ParseFiles("templates/canvas.html")
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
	}

	// Render the template with the message as data.
	var renderedCanvas bytes.Buffer
	err = tmpl.Execute(&renderedCanvas, data)
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
	}

	return renderedCanvas.Bytes()
}
