package users_handler

import (
	"fmt"
	"github.com/jameycribbs/pythia/db"
	"github.com/jameycribbs/pythia/global_vars"
	"github.com/justinas/nosurf"
	"html/template"
	"net/http"
	"path"
)

type IndexTemplateData struct {
	Users       []*db.User
	CurrentUser *db.User
}

type TemplateData struct {
	Rec               *db.User
	CurrentUser       *db.User
	DontShowLoginLink bool
	CsrfToken         string
}

func Index(w http.ResponseWriter, r *http.Request, throwAway string, gv *global_vars.GlobalVars, currentUser *db.User) {
	if (currentUser == nil) || (currentUser.Level != "admin") {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	var err error

	templateData := IndexTemplateData{}

	templateData.CurrentUser = currentUser

	templateData.Users, err = gv.MyDB.FindUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lp := path.Join("templates", "layouts", "layout.html")
	fp := path.Join("templates", "users", "index.html")

	tmpl, _ := template.ParseFiles(lp, fp)
	err = tmpl.ExecuteTemplate(w, "layout", templateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func View(w http.ResponseWriter, r *http.Request, fileId string, gv *global_vars.GlobalVars, currentUser *db.User) {
	if (currentUser == nil) || (currentUser.Level != "admin") {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	rec, err := gv.MyDB.FindUser(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := TemplateData{CurrentUser: currentUser, Rec: rec}

	renderTemplate(w, "view", &templateData)
}

func New(w http.ResponseWriter, r *http.Request, throwaway string, sv *global_vars.GlobalVars, currentUser *db.User) {
	if (currentUser == nil) || (currentUser.Level != "admin") {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	templateData := TemplateData{CurrentUser: currentUser, CsrfToken: nosurf.Token(r)}

	renderTemplate(w, "new", &templateData)
}

func Create(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars, currentUser *db.User) {
	if (currentUser == nil) || (currentUser.Level != "admin") {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fileId, err := saveFormDataToDb(gv.MyDB, "", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%v", fileId), http.StatusFound)
}

func Edit(w http.ResponseWriter, r *http.Request, fileId string, gv *global_vars.GlobalVars, currentUser *db.User) {
	if (currentUser == nil) || (currentUser.Level != "admin") {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	rec, err := gv.MyDB.FindUser(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := TemplateData{CurrentUser: currentUser, Rec: rec, CsrfToken: nosurf.Token(r)}

	renderTemplate(w, "edit", &templateData)
}

func Update(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars, currentUser *db.User) {
	if (currentUser == nil) || (currentUser.Level != "admin") {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fileId := r.FormValue("fileId")

	_, err := saveFormDataToDb(gv.MyDB, fileId, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%v", fileId), http.StatusFound)
}

func Delete(w http.ResponseWriter, r *http.Request, fileId string, gv *global_vars.GlobalVars, currentUser *db.User) {
	if (currentUser == nil) || (currentUser.Level != "admin") {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	rec, err := gv.MyDB.FindUser(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := TemplateData{CurrentUser: currentUser, Rec: rec, CsrfToken: nosurf.Token(r)}

	renderTemplate(w, "delete", &templateData)
}

func Destroy(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars, currentUser *db.User) {
	if (currentUser == nil) || (currentUser.Level != "admin") {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fileId := r.FormValue("fileId")

	err := gv.MyDB.DeleteUser(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusFound)
}

//=============================================================================
// Helper Functions
//=============================================================================
func renderTemplate(w http.ResponseWriter, templateName string, templateData *TemplateData) {
	lp := path.Join("templates", "layouts", "layout.html")
	fp := path.Join("templates", "users", templateName+".html")

	tmpl, _ := template.ParseFiles(lp, fp)
	err := tmpl.ExecuteTemplate(w, "layout", templateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func saveFormDataToDb(myDB *db.DB, fileId string, r *http.Request) (string, error) {
	name := r.FormValue("name")
	login := r.FormValue("login")
	password := r.FormValue("password")
	level := r.FormValue("level")

	rec := &db.User{FileId: fileId, Name: name, Login: login, Password: []byte(password), Level: level}

	returnedFileId, err := myDB.SaveUser(rec)
	if err != nil {
		return "", err
	}

	return returnedFileId, nil
}
