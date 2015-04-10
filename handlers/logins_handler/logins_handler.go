package logins_handler

import (
	"fmt"
	"github.com/jameycribbs/pythia/global_vars"
	"github.com/jameycribbs/pythia/models"
	"github.com/justinas/nosurf"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"path"
)

type TemplateData struct {
	Msg               string
	CurrentUser       *models.User
	DontShowLoginLink bool
	CsrfToken         string
}

func New(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars, currentUser *models.User) {
	templateData := TemplateData{CurrentUser: currentUser, DontShowLoginLink: true, CsrfToken: nosurf.Token(r)}
	renderTemplate(w, "new", &templateData)
}

func Create(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars, currentUser *models.User) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	user, err := loginUser(login, password, gv)
	if err != nil {
		templateData := TemplateData{Msg: "Login unsuccessful", CurrentUser: currentUser, CsrfToken: nosurf.Token(r)}
		renderTemplate(w, "new", &templateData)
		return
	}

	session, _ := gv.SessionStore.Get(r, "pythia")
	session.Values["user"] = user.FileId
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
}

func Logout(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars, currentUser *models.User) {
	session, _ := gv.SessionStore.Get(r, "pythia")
	delete(session.Values, "user")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
}

//=============================================================================
// Helper Functions
//=============================================================================
func loginUser(login string, password string, gv *global_vars.GlobalVars) (models.User, error) {
	var user models.User

	id, err := gv.MyDB.FindFirstIdForField("users", "login", login)
	if err != nil {
		return user, err
	}

	err = gv.MyDB.Find("users", &user, id)
	if err != nil {
		return user, err
	}

	fmt.Println(user)
	fmt.Println(password)

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil
}

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
