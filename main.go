package main

import (
	"net/http"
)

func main() {
	s := NewServer()
	go s.Listen()
	http.ListenAndServe(":8080", nil)
}
