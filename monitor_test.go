package goat

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"log"

	"github.com/stretchr/testify/assert"
)

type TestMonitorHandler struct{}

func (h *TestMonitorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "this is monitor test")
}

func Test_Monitor(t *testing.T) {
	h := &TestMonitorHandler{}
	m := NewMonitor()
	server := httptest.NewServer(m.Monitor(h))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(m.Get())
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status Code does not match")
}
