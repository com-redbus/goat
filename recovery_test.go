package goat

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)


type TestPanicHandler struct{}

func (h *TestPanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("Oh No its an error")
}

func Test_Recovery(t *testing.T) {
	h := &TestPanicHandler{}
	server := httptest.NewServer(Recovery(h))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Internal Server Error")
}
