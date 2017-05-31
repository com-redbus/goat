package goat

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"compress/gzip"

	"io/ioutil"

	"github.com/stretchr/testify/assert"
)

type TestHandler struct{}

func (h *TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Its a test handler")
}

func Test_Middleware_Default(t *testing.T) {
	h := &TestHandler{}
	commonMiddlewares := CommonMiddlewares()
	server := httptest.NewServer(commonMiddlewares.Then(h))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Internal Server Error")
}

func Test_Middleware_Panic(t *testing.T) {
	h := &TestPanicHandler{}
	commonMiddlewares := CommonMiddlewares()
	server := httptest.NewServer(commonMiddlewares.Then(h))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Internal Server Error")
}

func Test_Middleware_NoCache(t *testing.T) {
	h := &TestNoCacheHandler{}
	commonMiddlewares := CommonMiddlewares()
	server := httptest.NewServer(commonMiddlewares.Then(h))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "max-age=0, no-cache, no-store, must-revalidate", resp.Header.Get("Cache-Control"), "No Cache not working")
}

func Test_Middleware_Compression(t *testing.T) {
	h := &TestCompressionHandler{}
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.foo/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	commonMiddlewares := CommonMiddlewares()
	compressionAddedMiddleware := commonMiddlewares.Append(Compression)
	x := compressionAddedMiddleware.Then(h)
	x.ServeHTTP(rr, req)
	resp := rr.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status Code not same")
	assert.Equal(t, "gzip", resp.Header.Get("Content-Encoding"), "No gzip Content Encoding Header")
	assert.Equal(t, "Accept-Encoding", resp.Header.Get("Vary"), "No gzip Vary Header")

	reader, _ := gzip.NewReader(resp.Body)
	defer reader.Close()
	b, _ := ioutil.ReadAll(reader)
	assert.Equal(t, "this is compression test", string(b), "No gzip string doesnt match")
}

func Test_Middleware_Monitor(t *testing.T) {
	h := &TestMonitorHandler{}
	m := NewMonitor()
	commonMiddlewares := CommonMiddlewares()
	monitAddedMiddleware := commonMiddlewares.Append(m.Monitor)
	server := httptest.NewServer(monitAddedMiddleware.Then(h))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(m.Get())
	b, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status Code does not match")
	assert.Equal(t, "this is monitor test", string(b), "Response string doesnt match")
}
