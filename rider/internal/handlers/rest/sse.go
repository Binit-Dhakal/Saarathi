package rest

import (
	"net/http"
	"sync"

	"github.com/Binit-Dhakal/Saarathi/rider/internal/application"
)

type TripUpdateHandler struct {
	updateSvc   application.RiderUpdateService
	connections map[string]*Client
	mu          sync.RWMutex
}

func NewTripUpdateHandler(updateSvc application.RiderUpdateService) *TripUpdateHandler {
	return &TripUpdateHandler{
		updateSvc:   updateSvc,
		connections: map[string]*Client{},
	}
}

func (t *TripUpdateHandler) TripUpdate(w http.ResponseWriter, r *http.Request) {
	tripID := r.URL.Query().Get("tripId")
	if tripID == "" {
		http.Error(w, "missing tripId", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	client := &Client{
		id:         tripID,
		send:       make(chan any),
		flusher:    flusher,
		w:          w,
		reqContext: r.Context(),
	}

	client.Start()

	t.mu.Lock()
	t.connections[tripID] = client
	t.mu.Unlock()

	<-client.done

	t.mu.Lock()
	delete(t.connections, tripID)
	t.mu.Lock()
}

func (t *TripUpdateHandler) NotifyRider(tripID string, payload []byte) {
	t.mu.RLock()
	client, ok := t.connections[tripID]
	t.mu.RUnlock()

	if ok {
		select {
		case client.send <- payload:
			return
		default:
			// Client channel is backed up; drop the message
		}
	}
}
