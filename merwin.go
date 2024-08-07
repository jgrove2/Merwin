package merwin

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/jgrove2/Merwin/hub"
)

type Merwin struct {
	Hub hub.Hub
}

func NewMerwin(logDir string) *Merwin {
	return &Merwin{
		Hub: *hub.NewHub(),
	}
}

func (m *Merwin) Run() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting Hub...")
	go m.Hub.Run()
	slog.Info("Hub is running")

	http.HandleFunc("/", m.serveIndex)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWs(&m.Hub, w, r)
	})
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func (m *Merwin) serveIndex(w http.ResponseWriter, r *http.Request) {
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
