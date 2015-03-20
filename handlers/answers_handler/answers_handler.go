package answers_handler

import (
	"fmt"
	"github.com/jameycribbs/pythia/db"
	"github.com/jameycribbs/pythia/global_vars"
	"html/template"
	"net/http"
	"path"
)

type IndexTemplateData struct {
	SearchTagsString  string
	Answers           []*db.Answer
	CurrentUser       *db.User
	DontShowLoginLink bool
}

type TemplateData struct {
	Rec               *db.Answer
	CurrentUser       *db.User
	DontShowLoginLink bool
}

func Index(w http.ResponseWriter, r *http.Request, throwAway string, gv *global_vars.GlobalVars, currentUser *db.User) {
	var err error

	templateData := IndexTemplateData{}

	templateData.CurrentUser = currentUser

	templateData.SearchTagsString = r.FormValue("searchTags")

	templateData.Answers, err = gv.MyDB.FindAnswers(templateData.SearchTagsString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lp := path.Join("templates", "layouts", "layout.html")
	fp := path.Join("templates", "answers", "index.html")

	tmpl, _ := template.ParseFiles(lp, fp)
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
	templateData := TemplateData{CurrentUser: currentUser}

	renderTemplate(w, "new", &templateData)
}

func Create(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars, currentUser *db.User) {
	fileId, err := saveFormDataToDb(gv.MyDB, "", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/answers/%v", fileId), http.StatusFound)
}

func Edit(w http.ResponseWriter, r *http.Request, fileId string, gv *global_vars.GlobalVars, currentUser *db.User) {
	rec, err := gv.MyDB.FindAnswer(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := TemplateData{CurrentUser: currentUser, Rec: rec}
	renderTemplate(w, "edit", &templateData)
}

func Update(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars, currentUser *db.User) {
	fileId := r.FormValue("fileId")

	_, err := saveFormDataToDb(gv.MyDB, fileId, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/answers/%v", fileId), http.StatusFound)
}

func Delete(w http.ResponseWriter, r *http.Request, fileId string, gv *global_vars.GlobalVars, currentUser *db.User) {
	rec, err := gv.MyDB.FindAnswer(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := TemplateData{CurrentUser: currentUser, Rec: rec}
	renderTemplate(w, "delete", &templateData)
}

func Destroy(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars, currentUser *db.User) {
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
