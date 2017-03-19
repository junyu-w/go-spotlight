package server

import (
	"fmt"
	"net/http"
	"regexp"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	dir := r.FormValue("dir")
	queryTerms := r.FormValue("query")
	res, err := runQuery(dir, queryTerms)
	if err != nil {
		http.Error(w, err.Error(), 501)
	}
	fmt.Fprintf(w, res)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}

// validation
var validPath = regexp.MustCompile("^/(search)$")

func StartServer(port string) {
	http.HandleFunc("/search", makeHandler(searchHandler))
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
