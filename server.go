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
	meta   meta
	listen string
	srv    *http.Server
	done   chan (bool)
}

// New creates a new Server
func New(listen string, opts ...Option) (*Server, error) {
	s := new(Server)
	s.meta.Version = "0.0.0"
	s.meta.Description = "MYOB Technical Test"
	s.meta.Author = "Felix Hanley"

	s.listen = listen

	// Set variadic options passed
	for _, option := range opts {
		err := option(s)
		if err != nil {
			return nil, err
		}
	}
	s.done = make(chan (bool), 1)

	// Create, don't run
	s.srv = &http.Server{Addr: s.listen}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.helloHandler)
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/meta", s.metaHandler)
	s.srv.Handler = mux

	return s, nil
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

// Stop the server.
func (s *Server) Stop() error {
	s.done <- true
	return nil
}

// Run starts the server in the current thread.
func (s *Server) Run() {
	s.srv.ListenAndServe()
	select {
	case <-s.done:
	}
}
