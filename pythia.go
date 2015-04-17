package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jameycribbs/ivy"
	"github.com/jameycribbs/pythia/global_vars"
	"github.com/jameycribbs/pythia/handlers/answers_handler"
	"github.com/jameycribbs/pythia/handlers/logins_handler"
	"github.com/jameycribbs/pythia/handlers/users_handler"
	"github.com/jameycribbs/pythia/models"
	"github.com/justinas/nosurf"
	"net/http"
	"os"
)

func main() {
	var port string

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error getting hostname:", err)
	}

	if hostname == "pythia" {
		port = ":80"
	} else {
		port = ":8080"
	}

	fieldsToIndex := make(map[string][]string)
	fieldsToIndex["answers"] = []string{"tags"}
	fieldsToIndex["users"] = []string{"login"}

	db, err := ivy.OpenDB("data", fieldsToIndex)
	if err != nil {
		fmt.Println("Database initialization failed:", err)
	}

	defer db.Close()

	store := sessions.NewCookieStore([]byte("pythia-is-awesome"))

	gv := global_vars.GlobalVars{MyDB: db, SessionStore: store}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	r := mux.NewRouter()
	r.HandleFunc("/", makeHandler(answers_handler.Index, &gv)).Methods("GET")

	r.HandleFunc("/answers", makeHandler(answers_handler.Index, &gv)).Methods("GET")
	r.HandleFunc("/answers/search", makeHandler(answers_handler.Index, &gv)).Methods("POST")
	r.HandleFunc("/answers/{id:[0-9]+}", makeHandler(answers_handler.View, &gv)).Methods("GET")
	r.HandleFunc("/answers/new", makeHandler(answers_handler.New, &gv)).Methods("GET")
	r.HandleFunc("/answers/create", makeHandler(answers_handler.Create, &gv)).Methods("POST")
	r.HandleFunc("/answers/{id:[0-9]+}/edit", makeHandler(answers_handler.Edit, &gv)).Methods("GET")
	r.HandleFunc("/answers/update", makeHandler(answers_handler.Update, &gv)).Methods("POST")
	r.HandleFunc("/answers/{id:[0-9]+}/delete", makeHandler(answers_handler.Delete, &gv)).Methods("GET")
	r.HandleFunc("/answers/destroy", makeHandler(answers_handler.Destroy, &gv)).Methods("POST")

	r.HandleFunc("/users", makeHandler(users_handler.Index, &gv)).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", makeHandler(users_handler.View, &gv)).Methods("GET")
	r.HandleFunc("/users/new", makeHandler(users_handler.New, &gv)).Methods("GET")
	r.HandleFunc("/users/create", makeHandler(users_handler.Create, &gv)).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}/edit", makeHandler(users_handler.Edit, &gv)).Methods("GET")
	r.HandleFunc("/users/update", makeHandler(users_handler.Update, &gv)).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}/delete", makeHandler(users_handler.Delete, &gv)).Methods("GET")
	r.HandleFunc("/users/destroy", makeHandler(users_handler.Destroy, &gv)).Methods("POST")

	r.HandleFunc("/logins/new", makeHandler(logins_handler.New, &gv)).Methods("GET")
	r.HandleFunc("/logins/create", makeHandler(logins_handler.Create, &gv)).Methods("POST")
	r.HandleFunc("/logout", makeHandler(logins_handler.Logout, &gv)).Methods("GET")

	http.Handle("/", r)

	csrfHandler := nosurf.New(http.DefaultServeMux)

	csrfHandler.SetFailureHandler(http.HandlerFunc(failHand))

	http.ListenAndServe(port, csrfHandler)
}

func failHand(w http.ResponseWriter, r *http.Request) {
	// will return the reason of the failure
	fmt.Fprintf(w, "%s\n", nosurf.Reason(r))
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string, *global_vars.GlobalVars, *models.User),
	gv *global_vars.GlobalVars) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		currentUser, err := getCurrentUser(r, gv)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		vars := mux.Vars(r)

		fn(w, r, vars["id"], gv, currentUser)
	}
}

func getCurrentUser(r *http.Request, gv *global_vars.GlobalVars) (*models.User, error) {
	var user models.User

	session, _ := gv.SessionStore.Get(r, "pythia")

	userId, ok := session.Values["user"]

	if !ok {
		return nil, nil
	}

	err := gv.MyDB.Find("users", &user, userId.(string))
	if err != nil {
		return nil, err
	}

	return &user, nil
}
