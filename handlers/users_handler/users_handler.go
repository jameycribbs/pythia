package users_handler

import (
	"fmt"
	"github.com/jameycribbs/pythia/global_vars"
	"github.com/jameycribbs/pythia/models"
	"github.com/justinas/nosurf"
	"html/template"
	"net/http"
	"path"
)

type IndexTemplateData struct {
	Users       []*models.User
	CurrentUser *models.User
}

type TemplateData struct {
	Rec               *models.User
	CurrentUser       *models.User
	DontShowLoginLink bool
	CsrfToken         string
}

func Index(w http.ResponseWriter, r *http.Request, throwAway string, gv *global_vars.GlobalVars, currentUser *models.User) {
	if (currentUser == nil) || (currentUser.Level != "admin") {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	var err error

	templateData := IndexTemplateData{}

	templateData.CurrentUser = currentUser

	ids, err := gv.MyDB.FindAllIds("users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, id := range ids {
		user := models.User{}

		err = gv.MyDB.Find("users", &user, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		templateData.Users = append(templateData.Users, &user)
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

func View(w http.ResponseWriter, r *http.Request, fileId string, gv *global_vars.GlobalVars, currentUser *models.User) {
	if (currentUser == nil) || (currentUser.Level != "admin") {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	var rec models.User

	err := gv.MyDB.Find("users", &rec, fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := TemplateData{CurrentUser: currentUser, Rec: &rec}

	renderTemplate(w, "view", &templateData)
}

func New(w http.ResponseWriter, r *http.Request, throwaway string, sv *global_vars.GlobalVars, currentUser *models.User) {
	if (currentUser == nil) || (currentUser.Level != "admin") {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	templateData := TemplateData{CurrentUser: currentUser, CsrfToken: nosurf.Token(r)}

	renderTemplate(w, "new", &templateData)
}

func Create(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars, currentUser *models.User) {
	if (currentUser == nil) || (currentUser.Level != "admin") {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	name := r.FormValue("name")
	login := r.FormValue("login")
	password := r.FormValue("password")
	level := r.FormValue("level")

	rec := models.User{Name: name, Login: login, Password: []byte(password), Level: level}

	fileId, err := gv.MyDB.Create("users", rec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%v", fileId), http.StatusFound)
}

func Edit(w http.ResponseWriter, r *http.Request, fileId string, gv *global_vars.GlobalVars, currentUser *models.User) {
	var rec models.User

	if (currentUser == nil) || (currentUser.Level != "admin") {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	err := gv.MyDB.Find("users", &rec, fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := TemplateData{CurrentUser: currentUser, Rec: &rec, CsrfToken: nosurf.Token(r)}

	renderTemplate(w, "edit", &templateData)
}

func Update(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars, currentUser *models.User) {
	if (currentUser == nil) || (currentUser.Level != "admin") {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fileId := r.FormValue("fileId")

	name := r.FormValue("name")
	login := r.FormValue("login")
	password := r.FormValue("password")
	level := r.FormValue("level")

	rec := models.User{FileId: fileId, Name: name, Login: login, Password: []byte(password), Level: level}

	err := gv.MyDB.Update("users", rec, fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%v", fileId), http.StatusFound)
}

func Delete(w http.ResponseWriter, r *http.Request, fileId string, gv *global_vars.GlobalVars, currentUser *models.User) {
	var rec models.User

	if (currentUser == nil) || (currentUser.Level != "admin") {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	err := gv.MyDB.Find("users", &rec, fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := TemplateData{CurrentUser: currentUser, Rec: &rec, CsrfToken: nosurf.Token(r)}

	renderTemplate(w, "delete", &templateData)
}

func Destroy(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars, currentUser *models.User) {
	if (currentUser == nil) || (currentUser.Level != "admin") {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fileId := r.FormValue("fileId")

	err := gv.MyDB.Delete("users", fileId)
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
