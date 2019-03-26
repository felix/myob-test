package server

import (
	"encoding/json"
	"io"
	"net/http"
	"sync/atomic"
)

type meta struct {
	Version       string `json:"version"`
	Description   string `json:"description"`
	Author        string `json:"author"`
	TotalRequests int64  `json:"total_requests"`
}

// Server defines the simple API server.
type Server struct {
	meta meta
}

// helloHandler serves "Hello World" as plain text.
func (s *Server) helloHandler(w http.ResponseWriter, _ *http.Request) {
	s.updateMeta()
	// Set explicitly
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "Hello world\n")
}

// healthHandler serves a status 204 indicating server is alive.
func (s *Server) healthHandler(w http.ResponseWriter, _ *http.Request) {
	s.updateMeta()
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "OK\n")
}

// metaHandler sends JSON stats for server.
func (s *Server) metaHandler(w http.ResponseWriter, _ *http.Request) {
	s.updateMeta()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.meta)
}

func (s *Server) updateMeta() {
	atomic.AddInt64(&s.meta.TotalRequests, 1)
}
