package global_vars

import (
	"github.com/gorilla/sessions"
	"github.com/jameycribbs/pythia/db"
)

type GlobalVars struct {
	MyDB         *db.DB
	SessionStore *sessions.CookieStore
}
