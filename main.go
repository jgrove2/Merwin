package main

import (
	"log"
	"net/http"

	"github.com/jgrove2/browser_game_engine/hub"
)

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "method not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, "templates/index.html")
}

func main() {
	mainHub := hub.NewHub()
	go mainHub.Run()

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWs(mainHub, w, r)
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
