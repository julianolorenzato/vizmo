package main

import (
	"fmt"
	"net/http"
)

func main() {
	// configure the songs directory name and port
	const dirName = "hls"
	const port = 8080

	http.Handle("/", addHeaders(http.FileServer(http.Dir(dirName))))

	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

// addHeaders will act as middleware to give us CORS support
func addHeaders(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	}
}
