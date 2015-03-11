package db

import (
	"os"
	"sync"
)

type DB struct {
	path            string
	answersPath     string
	answersLock     *sync.RWMutex
	usersPath       string
	usersLock       *sync.RWMutex
	answerTagsIndex map[string][]string
}

func OpenDB(dbPath string) (*DB, error) {
	db := &DB{}

	db.path = dbPath
	db.answersPath = dbPath + "/answers"
	db.usersPath = dbPath + "/users"

	err := os.MkdirAll(db.answersPath, 0777)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(db.usersPath, 0777)
	if err != nil {
		return nil, err
	}

	db.answersLock = new(sync.RWMutex)
	db.usersLock = new(sync.RWMutex)

	db.answerTagsIndex = make(map[string][]string)

	err = db.initAnswerTagsIndex()

	return db, err
}

/*****************************************************************************/
// Public DB Methods
/*****************************************************************************/

/*---------- Close ----------*/
func (db *DB) Close() {
	db.answersLock.Lock()
	defer db.answersLock.Unlock()

	db.usersLock.Lock()
	defer db.usersLock.Unlock()
}

//=============================================================================
// Helper Functions
//=============================================================================

/*---------- stringInSlice ----------*/
func stringInSlice(s string, list []string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}
