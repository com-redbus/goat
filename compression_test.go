package goat

import (
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"fmt"
)

type TestCompressionHandler struct{}

func (h *TestCompressionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "this is compression test")
}
func Test_Compression(t *testing.T) {
	h := &TestCompressionHandler{}
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.foo/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	x := Compression(h)
	x.ServeHTTP(rr, req)
	resp := rr.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Internal Server Error")
	assert.Equal(t, "gzip", resp.Header.Get("Content-Encoding"), "No gzip Content Encoding Header")
	assert.Equal(t, "Accept-Encoding", resp.Header.Get("Vary"), "No gzip Vary Header")

	reader, _ := gzip.NewReader(resp.Body)
	defer reader.Close()
	b, _ := ioutil.ReadAll(reader)
	assert.Equal(t, "this is compression test", string(b), "No gzip string doesnt match")
}
