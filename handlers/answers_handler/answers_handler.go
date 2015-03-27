package answers_handler

import (
	"fmt"
	"github.com/jameycribbs/pythia/db"
	"github.com/jameycribbs/pythia/global_vars"
	"github.com/justinas/nosurf"
	"html/template"
	"net/http"
	"path"
	"strings"
	"time"
)

type IndexTemplateData struct {
	SearchTagsString  string
	Answers           []*db.Answer
	CurrentUser       *db.User
	DontShowLoginLink bool
	CurrentUserAdmin  bool
	CsrfToken         string
	AvailableTags     []string
}

type TemplateData struct {
	Rec               *db.Answer
	CurrentUser       *db.User
	DontShowLoginLink bool
	CsrfToken         string
}

func Index(w http.ResponseWriter, r *http.Request, throwAway string, gv *global_vars.GlobalVars, currentUser *db.User) {
	var err error

	funcMap := template.FuncMap{
		"panelClass": func(i int) string {
			if i == 0 {
				return "collapse in"
			} else {
				return "collapse"
			}
		},
		"tagBadges": func(tagString string) template.HTML {
			var formattedString string

			for _, tag := range strings.Split(tagString, " ") {
				formattedString = formattedString + "<span class='label label-primary'>" + tag + "</span> "
			}

			return template.HTML(formattedString)
		}}

	templateData := IndexTemplateData{}

	templateData.AvailableTags = gv.MyDB.AvailableAnswerTags

	templateData.CurrentUser = currentUser

	if (currentUser != nil) && (currentUser.Level == "admin") {
		templateData.CurrentUserAdmin = true
	} else {
		templateData.CurrentUserAdmin = false
	}

	if r.FormValue("searchTags") != "" {
		templateData.SearchTagsString = r.FormValue("searchTags")

		templateData.Answers, err = gv.MyDB.FindAnswers(templateData.SearchTagsString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	templateData.CsrfToken = nosurf.Token(r)

	lp := path.Join("templates", "layouts", "layout.html")
	fp := path.Join("templates", "answers", "index.html")

	tmpl := template.New("idx").Funcs(funcMap)

	tmpl, err = tmpl.ParseFiles(lp, fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", templateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func View(w http.ResponseWriter, r *http.Request, fileId string, gv *global_vars.GlobalVars, currentUser *db.User) {
	rec, err := gv.MyDB.FindAnswer(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := TemplateData{CurrentUser: currentUser, Rec: rec}

	renderTemplate(w, "view", &templateData)
}

func New(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars, currentUser *db.User) {
	if currentUser == nil {
		http.Redirect(w, r, "/answers", http.StatusFound)
		return
	}

	templateData := TemplateData{CurrentUser: currentUser, CsrfToken: nosurf.Token(r)}

	renderTemplate(w, "new", &templateData)
}

func Create(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars, currentUser *db.User) {
	if currentUser == nil {
		http.Redirect(w, r, "/answers", http.StatusFound)
		return
	}

	question := r.FormValue("question")
	answer := r.FormValue("answer")
	tags := r.FormValue("tags")

	rec := &db.Answer{Question: question, Answer: answer, Tags: tags, CreatedById: currentUser.FileId, CreatedAt: time.Now(),
		UpdatedById: currentUser.FileId, UpdatedAt: time.Now()}

	fileId, err := gv.MyDB.SaveAnswer(rec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/answers/%v", fileId), http.StatusFound)
}

func Edit(w http.ResponseWriter, r *http.Request, fileId string, gv *global_vars.GlobalVars, currentUser *db.User) {
	if currentUser == nil {
		http.Redirect(w, r, "/answers", http.StatusFound)
		return
	}

	rec, err := gv.MyDB.FindAnswer(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := TemplateData{CurrentUser: currentUser, Rec: rec, CsrfToken: nosurf.Token(r)}
	renderTemplate(w, "edit", &templateData)
}

func Update(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars, currentUser *db.User) {
	if currentUser == nil {
		http.Redirect(w, r, "/answers", http.StatusFound)
		return
	}

	fileId := r.FormValue("fileId")
	question := r.FormValue("question")
	answer := r.FormValue("answer")
	tags := r.FormValue("tags")

	rec := &db.Answer{FileId: fileId, Question: question, Answer: answer, Tags: tags, UpdatedById: currentUser.FileId,
		UpdatedAt: time.Now()}

	_, err := gv.MyDB.SaveAnswer(rec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/answers/%v", fileId), http.StatusFound)
}

func Delete(w http.ResponseWriter, r *http.Request, fileId string, gv *global_vars.GlobalVars, currentUser *db.User) {
	if currentUser == nil {
		http.Redirect(w, r, "/answers", http.StatusFound)
		return
	}

	rec, err := gv.MyDB.FindAnswer(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := TemplateData{CurrentUser: currentUser, Rec: rec, CsrfToken: nosurf.Token(r)}
	renderTemplate(w, "delete", &templateData)
}

func Destroy(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars, currentUser *db.User) {
	if currentUser == nil {
		http.Redirect(w, r, "/answers", http.StatusFound)
		return
	}

	fileId := r.FormValue("fileId")

	err := gv.MyDB.DeleteAnswer(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/answers", http.StatusFound)
}

//=============================================================================
// Helper Functions
//=============================================================================
func renderTemplate(w http.ResponseWriter, templateName string, templateData *TemplateData) {
	lp := path.Join("templates", "layouts", "layout.html")
	fp := path.Join("templates", "answers", templateName+".html")

	tmpl, _ := template.ParseFiles(lp, fp)
	err := tmpl.ExecuteTemplate(w, "layout", templateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
