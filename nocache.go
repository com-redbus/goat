package goat

import "net/http"

//NoCache middleware func which adds the no cache headers to the response
func NoCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=0, no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		//call the next handler
		next.ServeHTTP(w, r)
	})
}
