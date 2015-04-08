package models

import (
	"fmt"
	"github.com/jameycribbs/ivy"
	"time"
)

type Answer struct {
	FileId      string    `json:"-"y`
	Question    string    `json:"question"`
	Answer      string    `json:"answer"`
	CreatedById string    `json:"createdbyid"`
	UpdatedById string    `json:"updatedbyid"`
	CreatedAt   time.Time `json:"createdat"`
	UpdatedAt   time.Time `json:"updatedat"`
	Tags        []string  `json:"tags"`
	CreatedBy   string    `json:"-"`
	UpdatedBy   string    `json:"-"`
}

func (answer *Answer) AfterFind(db *ivy.DB, fileId string) {
	var createUser User
	var updateUser User

	*answer = Answer(*answer)

	answer.FileId = fileId

	err := db.Find("users", &createUser, answer.CreatedById)
	if err != nil {
		fmt.Println("Could not find creator:", err)
	}

	answer.CreatedBy = createUser.Name

	err = db.Find("users", &updateUser, answer.UpdatedById)
	if err != nil {
		fmt.Println("Could not find updater:", err)
	}

	answer.UpdatedBy = updateUser.Name
}
