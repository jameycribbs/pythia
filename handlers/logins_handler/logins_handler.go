package logins_handler

import (
	"github.com/jameycribbs/pythia/db"
	"github.com/jameycribbs/pythia/global_vars"
	"html/template"
	"net/http"
	"path"
)

type TemplateData struct {
	Msg               string
	CurrentUser       *db.User
	DontShowLoginLink bool
}

func New(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars, currentUser *db.User) {
	templateData := TemplateData{CurrentUser: currentUser, DontShowLoginLink: true}

	renderTemplate(w, "new", &templateData)
}

func Create(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars, currentUser *db.User) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	user, err := gv.MyDB.LoginUser(login, password)
	if err != nil {
		templateData := TemplateData{Msg: "Login unsuccessful", CurrentUser: currentUser}
		renderTemplate(w, "new", &templateData)
		return
	}

	session, _ := gv.SessionStore.Get(r, "pythia")
	session.Values["user"] = user.FileId
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
}

func Logout(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars, currentUser *db.User) {
	session, _ := gv.SessionStore.Get(r, "pythia")
	delete(session.Values, "user")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
}

//=============================================================================
// Helper Functions
//=============================================================================
func renderTemplate(w http.ResponseWriter, templateName string, templateData *TemplateData) {
	lp := path.Join("templates", "layouts", "layout.html")
	fp := path.Join("templates", "logins", templateName+".html")

	tmpl, _ := template.ParseFiles(lp, fp)
	err := tmpl.ExecuteTemplate(w, "layout", templateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
