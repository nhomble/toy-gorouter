package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func ping(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Printf("Received message\n")
		w.Write([]byte("pong"))
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Printf("healthy\n")
		w.Write([]byte("healthy"))
	}
}

func main() {
	var port string
	flag.StringVar(&port, "port", "", "port to server http requests")
	flag.Parse()

	http.Handle("/api/v1/ping", http.HandlerFunc(ping))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal(err)
	}
}
