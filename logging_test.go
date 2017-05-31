package goat

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var notFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
})

func Test_Logger(t *testing.T) {
	ts := httptest.NewServer(Logger(notFoundHandler))
	defer ts.Close()

	var u bytes.Buffer
	u.WriteString(ts.URL)
	u.WriteString("/foo")
	res, err := http.Get(u.String())
	assert.NoError(t, err, "msg")
	if res != nil {
		defer res.Body.Close()
	}

	assert.Equal(t, http.StatusNotFound, res.StatusCode, "Not Found")
}
