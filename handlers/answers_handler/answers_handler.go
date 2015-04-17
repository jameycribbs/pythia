package answers_handler

import (
	"fmt"
	"github.com/jameycribbs/pythia/global_vars"
	"github.com/jameycribbs/pythia/models"
	"github.com/justinas/nosurf"
	"html/template"
	"net/http"
	"path"
	"strings"
	"time"
)

type IndexTemplateData struct {
	SearchTagsString  string
	Answers           []*models.Answer
	CurrentUser       *models.User
	DontShowLoginLink bool
	CurrentUserAdmin  bool
	CsrfToken         string
}

type TemplateData struct {
	Rec               *models.Answer
	CurrentUser       *models.User
	DontShowLoginLink bool
	CsrfToken         string
}

func Index(w http.ResponseWriter, r *http.Request, throwAway string, gv *global_vars.GlobalVars, currentUser *models.User) {
	var err error
	funcMap := template.FuncMap{
		"panelClass": func(i int) string {
			if i == 0 {
				return "collapse in"
			} else {
				return "collapse"
			}
		}}

	templateData := IndexTemplateData{}

	templateData.CurrentUser = currentUser

	if (currentUser != nil) && (currentUser.Level == "admin") {
		templateData.CurrentUserAdmin = true
	} else {
		templateData.CurrentUserAdmin = false
	}

	if r.FormValue("searchTags") != "" {
		templateData.SearchTagsString = r.FormValue("searchTags")

		ids, err := gv.MyDB.FindAllIdsForTags("answers", strings.Split(templateData.SearchTagsString, " "))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, id := range ids {
			answer := models.Answer{}

			err = gv.MyDB.Find("answers", &answer, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			templateData.Answers = append(templateData.Answers, &answer)
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

func View(w http.ResponseWriter, r *http.Request, fileId string, gv *global_vars.GlobalVars, currentUser *models.User) {
	if currentUser == nil {
		http.Redirect(w, r, "/answers", http.StatusFound)
		return
	}

	var rec models.Answer

	err := gv.MyDB.Find("answers", &rec, fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := TemplateData{CurrentUser: currentUser, Rec: &rec}

	renderTemplate(w, "view", &templateData)
}

func New(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars, currentUser *models.User) {
	if currentUser == nil {
		http.Redirect(w, r, "/answers", http.StatusFound)
		return
	}

	templateData := TemplateData{CurrentUser: currentUser, CsrfToken: nosurf.Token(r)}

	renderTemplate(w, "new", &templateData)
}

func Create(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars, currentUser *models.User) {
	if currentUser == nil {
		http.Redirect(w, r, "/answers", http.StatusFound)
		return
	}

	question := r.FormValue("question")
	answer := r.FormValue("answer")
	tags := r.FormValue("tags")

	rec := models.Answer{Question: question, Answer: answer, Tags: strings.Split(tags, " "), CreatedById: currentUser.FileId,
		CreatedAt: time.Now(), UpdatedById: currentUser.FileId, UpdatedAt: time.Now()}

	fileId, err := gv.MyDB.Create("answers", rec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, fmt.Sprintf("/answers/%v", fileId), http.StatusFound)
}

func Edit(w http.ResponseWriter, r *http.Request, fileId string, gv *global_vars.GlobalVars, currentUser *models.User) {
	if currentUser == nil {
		http.Redirect(w, r, "/answers", http.StatusFound)
		return
	}

	var rec models.Answer

	err := gv.MyDB.Find("answers", &rec, fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := TemplateData{CurrentUser: currentUser, Rec: &rec, CsrfToken: nosurf.Token(r)}
	renderTemplate(w, "edit", &templateData)
}

func Update(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars, currentUser *models.User) {
	if currentUser == nil {
		http.Redirect(w, r, "/answers", http.StatusFound)
		return
	}

	fileId := r.FormValue("fileId")
	question := r.FormValue("question")
	answer := r.FormValue("answer")
	tags := r.FormValue("tags")

	rec := models.Answer{FileId: fileId, Question: question, Answer: answer, Tags: strings.Split(tags, " "),
		UpdatedById: currentUser.FileId, UpdatedAt: time.Now()}

	err := gv.MyDB.Update("answers", rec, fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/answers/%v", fileId), http.StatusFound)
}

func Delete(w http.ResponseWriter, r *http.Request, fileId string, gv *global_vars.GlobalVars, currentUser *models.User) {
	if currentUser == nil {
		http.Redirect(w, r, "/answers", http.StatusFound)
		return
	}

	var rec models.Answer

	err := gv.MyDB.Find("answers", &rec, fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := TemplateData{CurrentUser: currentUser, Rec: &rec, CsrfToken: nosurf.Token(r)}
	renderTemplate(w, "delete", &templateData)
}

func Destroy(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars, currentUser *models.User) {
	if currentUser == nil {
		http.Redirect(w, r, "/answers", http.StatusFound)
		return
	}

	fileId := r.FormValue("fileId")

	err := gv.MyDB.Delete("answers", fileId)
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
