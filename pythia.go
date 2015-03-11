package main

import (
	"fmt"
	"github.com/jameycribbs/pythia/db"
	"github.com/jameycribbs/pythia/handlers"
	"net/http"
	"regexp"
)

var validPath = regexp.MustCompile("^/(new|create|edit|save|delete|view|index)*/*([a-zA-Z0-9]*)$")

func main() {
	myDB, err := db.OpenDB("data")
	if err != nil {
		fmt.Println("Database initialization failed:", err)
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/index/", makeHandler(handlers.AnswerIndexHandler, myDB))
	http.HandleFunc("/view/", makeHandler(handlers.AnswerViewHandler, myDB))
	http.HandleFunc("/new/", makeHandler(handlers.AnswerNewHandler, myDB))
	http.HandleFunc("/create/", makeHandler(handlers.AnswerCreateHandler, myDB))
	http.HandleFunc("/edit/", makeHandler(handlers.AnswerEditHandler, myDB))
	http.HandleFunc("/save/", makeHandler(handlers.AnswerSaveHandler, myDB))
	http.HandleFunc("/delete/", makeHandler(handlers.AnswerDeleteHandler, myDB))
	http.HandleFunc("/", makeHandler(handlers.AnswerIndexHandler, myDB))

	http.ListenAndServe(":8080", nil)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string, *db.DB), myDB *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2], myDB)
	}
}
