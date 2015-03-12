package user_handler

import (
	"fmt"
	"github.com/jameycribbs/pythia/db"
	"html/template"
	"net/http"
	"path"
)

type IndexData struct {
	SearchTagsString string
	Users            []*db.User
}

func Index(w http.ResponseWriter, r *http.Request, throwAway string, myDB *db.DB) {
	var err error

	indexData := IndexData{}

	indexData.SearchTagsString = r.FormValue("searchTags")

	indexData.Answers, err = myDB.FindUsers(indexData.SearchTagsString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lp := path.Join("tmpl", "layout.html")
	fp := path.Join("tmpl", "index.html")

	tmpl, _ := template.ParseFiles(lp, fp)
	err = tmpl.ExecuteTemplate(w, "layout", indexData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AnswerViewHandler(w http.ResponseWriter, r *http.Request, fileId string, myDB *db.DB) {
	a, err := myDB.FindAnswer(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "view", a)
}

func AnswerNewHandler(w http.ResponseWriter, r *http.Request, throwaway string, myDB *db.DB) {
	renderTemplate(w, "new", nil)
}

func AnswerCreateHandler(w http.ResponseWriter, r *http.Request, throwaway string, myDB *db.DB) {
	fileId, err := saveFormDataToDb(myDB, "", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/view/%v", fileId), http.StatusFound)
}

func AnswerEditHandler(w http.ResponseWriter, r *http.Request, fileId string, myDB *db.DB) {
	a, err := myDB.FindAnswer(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "edit", a)
}

func AnswerSaveHandler(w http.ResponseWriter, r *http.Request, fileId string, myDB *db.DB) {
	_, err := saveFormDataToDb(myDB, fileId, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/view/%v", fileId), http.StatusFound)
}

func AnswerDeleteHandler(w http.ResponseWriter, r *http.Request, fileId string, myDB *db.DB) {
	err := myDB.DeleteAnswer(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

//=============================================================================
// Helper Functions
//=============================================================================
func renderTemplate(w http.ResponseWriter, templateName string, a *db.Answer) {
	lp := path.Join("tmpl", "layout.html")
	fp := path.Join("tmpl", templateName+".html")

	tmpl, _ := template.ParseFiles(lp, fp)
	err := tmpl.ExecuteTemplate(w, "layout", a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func saveFormDataToDb(myDB *db.DB, fileId string, r *http.Request) (string, error) {
	question := r.FormValue("question")
	answer := r.FormValue("answer")
	tags := r.FormValue("tags")

	a := &db.Answer{FileId: fileId, Question: question, Answer: answer, Tags: tags}

	returnedFileId, err := myDB.SaveAnswer(a)
	if err != nil {
		return "", err
	}

	return returnedFileId, nil
}
