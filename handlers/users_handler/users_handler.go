package users_handler

import (
	"fmt"
	"github.com/jameycribbs/pythia/db"
	"html/template"
	"net/http"
	"path"
)

func Index(w http.ResponseWriter, r *http.Request, throwAway string, myDB *db.DB) {
	users, err := myDB.FindUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lp := path.Join("templates", "answers", "layout.html")
	fp := path.Join("templates", "answers", "index.html")

	tmpl, _ := template.ParseFiles(lp, fp)
	err = tmpl.ExecuteTemplate(w, "layout", users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func View(w http.ResponseWriter, r *http.Request, fileId string, myDB *db.DB) {
	user, err := myDB.FindUser(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "view", user)
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

	http.Redirect(w, r, fmt.Sprintf("/view/%v", fileId), http.StatusFound)
}

func Edit(w http.ResponseWriter, r *http.Request, fileId string, myDB *db.DB) {
	user, err := myDB.FindUser(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "edit", user)
}

func Save(w http.ResponseWriter, r *http.Request, fileId string, myDB *db.DB) {
	_, err := saveFormDataToDb(myDB, fileId, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/view/%v", fileId), http.StatusFound)
}

func Delete(w http.ResponseWriter, r *http.Request, fileId string, myDB *db.DB) {
	err := myDB.DeleteUser(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

//=============================================================================
// Helper Functions
//=============================================================================
func renderTemplate(w http.ResponseWriter, templateName string, user *db.User) {
	lp := path.Join("templates", "users", "layout.html")
	fp := path.Join("templates", "users", templateName+".html")

	tmpl, _ := template.ParseFiles(lp, fp)
	err := tmpl.ExecuteTemplate(w, "layout", user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func saveFormDataToDb(myDB *db.DB, fileId string, r *http.Request) (string, error) {
	name := r.FormValue("name")
	login := r.FormValue("login")
	password := r.FormValue("password")

	user := &db.User{FileId: fileId, Name: name, Login: login, Password: password}

	returnedFileId, err := myDB.SaveUser(user)
	if err != nil {
		return "", err
	}

	return returnedFileId, nil
}
