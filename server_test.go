package server

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestHandlers(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{"/", "Hello world\n"},
		{"/health", "OK\n"},
		{"/meta", `{"version":"0.0.0","description":"MYOB Technical Test","author":"Felix Hanley","total_requests":3}
`},
	}

	s, _ := New("localhost:8080")

	for _, tt := range tests {
		req := httptest.NewRequest("GET", "http://localhost:8080"+tt.url, nil)
		w := httptest.NewRecorder()

		s.srv.Handler.ServeHTTP(w, req)

		resp := w.Result()
		if resp.StatusCode != 200 {
			t.Fatalf("Expected 200, got %d", resp.StatusCode)
		}
		actual, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Got error: %s", err)
		}
		if string(actual) != tt.expected {
			t.Fatalf("Expected %q, got %q", tt.expected, actual)
		}
	}
}
