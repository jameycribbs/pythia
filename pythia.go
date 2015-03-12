package main

import (
	"fmt"
	"github.com/jameycribbs/pythia/db"
	"github.com/jameycribbs/pythia/handlers/answers_handler"
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

	http.HandleFunc("/index/", makeHandler(answers_handler.Index, myDB))
	http.HandleFunc("/view/", makeHandler(answers_handler.View, myDB))
	http.HandleFunc("/new/", makeHandler(answers_handler.New, myDB))
	http.HandleFunc("/create/", makeHandler(answers_handler.Create, myDB))
	http.HandleFunc("/edit/", makeHandler(answers_handler.Edit, myDB))
	http.HandleFunc("/save/", makeHandler(answers_handler.Save, myDB))
	http.HandleFunc("/delete/", makeHandler(answers_handler.Delete, myDB))
	http.HandleFunc("/", makeHandler(answers_handler.AnswerIndex, myDB))

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
