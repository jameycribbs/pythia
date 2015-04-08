package global_vars

import (
	"github.com/gorilla/sessions"
	"github.com/jameycribbs/ivy"
)

type GlobalVars struct {
	MyDB         *ivy.DB
	SessionStore *sessions.CookieStore
}
