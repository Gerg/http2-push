package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	address  = ":8080"
	PAGE_URL = "/"
	PAGE     = `<!DOCTYPE html>
<html>
<script type="text/javascript" src="script.js" defer></script>
<title>HTTP/2 Push example</title>
<h1>HTTP/2 Push example</h1>
<p id="main"></p>
</html>`
	SCRIPT_URL = "/script.js"
	SCRIPT     = `(function() {
	document.querySelector("#main").innerHTML = 'HTTP/2 Push example content';	
})();`
)

func script(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript")
	io.WriteString(w, SCRIPT)
}

func page(w http.ResponseWriter, r *http.Request) {
	if pusher, ok := w.(http.Pusher); ok {
		// HTTP/2 Push is supported
		w.Header().Set("Link", "<script.js>; rel=preload; as=script")
		if err := pusher.Push(SCRIPT_URL, nil); err != nil {
			fmt.Printf("Failed to push: %v\n", err)
		}
	}
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, PAGE)
}

func main() {
	http.HandleFunc(PAGE_URL, page)
	http.HandleFunc(SCRIPT_URL, script)
	server := &http.Server{
		Addr:         address,
		Handler:      h2c.NewHandler(http.DefaultServeMux, &http2.Server{}),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Printf("Listening on port %s\n", address)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
