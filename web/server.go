package web

import (
	"fmt"
	"github.com/DrakeW/go-spotlight/api"
	"html/template"
	"net/http"
	"regexp"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	dir := r.FormValue("dir")
	queryTerms := r.FormValue("query")
	res, err := api.RunQuery(dir, queryTerms)
	if err != nil {
		http.Error(w, err.Error(), 501)
	}
	fmt.Fprintf(w, res)
}

func handler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index")
}

// template caching
var templates = template.Must(template.ParseFiles("web/index.html"))

func renderTemplate(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
var validPath = regexp.MustCompile("^/(search|view)$")

func StartServer(port string) {
	http.HandleFunc("/view", makeHandler(handler))
	http.HandleFunc("/search", makeHandler(searchHandler))
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
