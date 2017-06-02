package goat

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_XSS(t *testing.T) {
	h := &TestNoCacheHandler{}
	server := httptest.NewServer(XSS(h))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "1; mode=block", resp.Header.Get("X-XSS-Protection"), "XSS not working")
}
