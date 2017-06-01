package goat

import (
	"bytes"
	"html"
	"html/template"
	"log"
	"net/http"
	"time"
)

//logger template is the type of string that will get logged to the console
var loggerTemplate = "{{.StartTime}} || {{.Status}} || \t {{.Duration}} | {{.HostName}} | {{.Method}} | {{.Path}} \n"

//loggerStruct stores the value of the logs
type loggerStruct struct {
	StartTime string
	Status    int
	Duration  time.Duration
	HostName  string
	Method    string
	Path      string
}

//Logger func handler for logging middleware
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		//wrap the response writer to get the status code
		//cant access status code from http.ResponseWriter
		nrw := NewResponseWriter(w)
		//call the next handler
		next.ServeHTTP(nrw, r)
		//response := w.(ResponseWriter)

		ls := &loggerStruct{
			StartTime: start.Format(time.RFC3339),
			Status:    nrw.Status(),
			Duration:  time.Since(start),
			HostName:  r.Host,
			Method:    r.Method,
			Path:      r.URL.Path,
		}

		t := template.New("logger_template")
		tem := template.Must(t.Parse(loggerTemplate))

		buf := &bytes.Buffer{}
		tem.Execute(buf, ls)

		log.Println(html.UnescapeString(buf.String()))
	})
}
