package goat

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestNoCacheHandler struct{}

func (h *TestNoCacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "This is a test string")
}
func Test_NoCache(t *testing.T) {
	h := &TestNoCacheHandler{}
	server := httptest.NewServer(NoCache(h))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "max-age=0, no-cache, no-store, must-revalidate", resp.Header.Get("Cache-Control"), "No Cache not working")
}
