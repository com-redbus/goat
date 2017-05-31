package goat

import (
	"errors"
	"net/http"
)

//Recovery middleware for catching panics in code globally
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			//get hold of the panic
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("unknown error")
				}
				//need to decide what type of format to send to the response
				//will work for now
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		//call the next handler
		next.ServeHTTP(w, r)
	})
}
