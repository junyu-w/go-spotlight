package server

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "To be finished!")
}

func StartServer() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
