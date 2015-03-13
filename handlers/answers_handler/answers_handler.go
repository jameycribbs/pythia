package answers_handler

import (
	"fmt"
	"github.com/jameycribbs/pythia/db"
	"html/template"
	"net/http"
	"path"
)

type IndexData struct {
	SearchTagsString string
	Answers          []*db.Answer
}

func Index(w http.ResponseWriter, r *http.Request, throwAway string, myDB *db.DB) {
	var err error

	indexData := IndexData{}

	indexData.SearchTagsString = r.FormValue("searchTags")

	indexData.Answers, err = myDB.FindAnswers(indexData.SearchTagsString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lp := path.Join("templates", "layouts", "layout.html")
	fp := path.Join("templates", "answers", "index.html")

	tmpl, _ := template.ParseFiles(lp, fp)
	err = tmpl.ExecuteTemplate(w, "layout", indexData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func View(w http.ResponseWriter, r *http.Request, fileId string, myDB *db.DB) {
	rec, err := myDB.FindAnswer(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "view", rec)
}

func New(w http.ResponseWriter, r *http.Request, throwaway string, myDB *db.DB) {
	renderTemplate(w, "new", nil)
}

func Create(w http.ResponseWriter, r *http.Request, throwaway string, myDB *db.DB) {
	fileId, err := saveFormDataToDb(myDB, "", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/answer/view/%v", fileId), http.StatusFound)
}

func Edit(w http.ResponseWriter, r *http.Request, fileId string, myDB *db.DB) {
	rec, err := myDB.FindAnswer(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "edit", rec)
}

func Save(w http.ResponseWriter, r *http.Request, fileId string, myDB *db.DB) {
	_, err := saveFormDataToDb(myDB, fileId, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/answer/view/%v", fileId), http.StatusFound)
}

func Delete(w http.ResponseWriter, r *http.Request, fileId string, myDB *db.DB) {
	err := myDB.DeleteAnswer(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/answers", http.StatusFound)
}

//=============================================================================
// Helper Functions
//=============================================================================
func renderTemplate(w http.ResponseWriter, templateName string, rec *db.Answer) {
	lp := path.Join("templates", "layouts", "layout.html")
	fp := path.Join("templates", "answers", templateName+".html")

	tmpl, _ := template.ParseFiles(lp, fp)
	err := tmpl.ExecuteTemplate(w, "layout", rec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func saveFormDataToDb(myDB *db.DB, fileId string, r *http.Request) (string, error) {
	question := r.FormValue("question")
	answer := r.FormValue("answer")
	tags := r.FormValue("tags")

	rec := &db.Answer{FileId: fileId, Question: question, Answer: answer, Tags: tags}

	returnedFileId, err := myDB.SaveAnswer(rec)
	if err != nil {
		return "", err
	}

	return returnedFileId, nil
}
