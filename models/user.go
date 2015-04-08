package models

import (
	"github.com/jameycribbs/ivy"
)

type User struct {
	FileId   string `json:"-"y`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password []byte `json:"password"`
	Level    string `json:"level"`
}

func (user *User) AfterFind(db *ivy.DB, fileId string) {
	*user = User(*user)

	user.FileId = fileId
}
