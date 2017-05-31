package goat

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

//gziResponseWriter wraps the io.Writer and ResponseWrirer
type gzipResponseWriter struct {
	io.Writer
	ResponseWriter
}

//Need to implement Write func because io.Writer interface needs this func
func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

//Handler struct provides pool of gzipWriter which can be reused many times from the pool
//Dont exactly understand why but saw it in https://github.com/NYTimes/gziphandler
type Handler struct {
	pool sync.Pool
	next http.Handler
}

//ServeHTTP func needs to ge implemented because http.handler interface needs this method otherwise Handler will not become a http.Handler
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		//if header does not have Accept-Encoding header just move on
		h.next.ServeHTTP(w, r)
		return
	}

	if w.Header().Get("Content-Encoding") == "gzip" {
		//if already compressed move on
		h.next.ServeHTTP(w, r)
		return
	}

	//inspired by https://github.com/NYTimes/gziphandler
	//get gzip Writer from pool
	gz := h.pool.Get().(*gzip.Writer)
	//dont forget to put to back into the pool after writing
	defer h.pool.Put(gz)
	//Reset the responseWriter  to original state , this allows to resuse a writer rather than creating a new one
	gz.Reset(w)

	//set the required headers
	headers := w.Header()
	headers.Set("Content-Encoding", "gzip")
	headers.Set("Vary", "Accept-Encoding")

	//wrap responseWriter to our response writer
	nrw := NewResponseWriter(w)
	//created gzipResponseWriter to pass to next handler
	grw := gzipResponseWriter{
		gz,
		nrw,
	}
	h.next.ServeHTTP(grw, r)
	//close writer
	gz.Close()
}

func Compression(next http.Handler) http.Handler {
	handler := &Handler{
		next: next,
	}
	handler.pool.New = func() interface{} {
		//write the compressed data to the writer using Default Compression
		gz, err := gzip.NewWriterLevel(ioutil.Discard, gzip.DefaultCompression)
		if err != nil {
			panic(err)
		}
		return gz
	}
	return handler
}
