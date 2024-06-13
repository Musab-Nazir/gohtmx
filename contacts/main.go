package main

import (
	"log"
	"net/http"
)

func main() {
	runServer()
}

func runServer() {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}
