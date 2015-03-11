package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type DB struct {
	path            string
	answersPath     string
	schemaLock      *sync.RWMutex
	answerTagsIndex map[string][]string
}

type Answer struct {
	FileId      string    `json:"fileid"`
	Question    string    `json:"question"`
	Answer      string    `json:"answer"`
	CreatedById string    `json:"createdbyid"`
	UpdatedById string    `json:"updatedbyid"`
	CreatedAt   time.Time `json:"createdat"`
	UpdatedAt   time.Time `json:"updatedat"`
	Tags        string    `json:"tags"`
}

func OpenDB(dbPath string) (*DB, error) {
	db := &DB{}

	db.path = dbPath
	db.answersPath = dbPath + "/answers"

	err := os.MkdirAll(db.answersPath, 0777)
	if err != nil {
		return nil, err
	}

	db.schemaLock = new(sync.RWMutex)
	db.answerTagsIndex = make(map[string][]string)

	err = db.initAnswerTagsIndex()

	return db, err
}

/*****************************************************************************/
// Public Methods
/*****************************************************************************/

/*---------- Close ----------*/
func (db *DB) Close() {
	db.schemaLock.Lock()
	defer db.schemaLock.Unlock()
}

/*---------- FindAnswers ----------*/
func (db *DB) FindAnswers(searchString string) ([]*Answer, error) {
	var answers []*Answer
	var searchTags []string

	db.schemaLock.RLock()
	defer db.schemaLock.RUnlock()

	if len(searchString) == 0 {
		for _, fileId := range db.fileIdsInAnswersDataDir() {
			answer, err := db.loadAnswer(fileId)
			if err != nil {
				return nil, err
			}

			answers = append(answers, answer)
		}
	} else {
		searchTags = strings.Split(searchString, " ")

		for fileId, tags := range db.answerTagsIndex {
			found := true

			for _, searchTag := range searchTags {
				if !stringInSlice(searchTag, tags) {
					found = false
					break
				}
			}

			if found {
				answer, err := db.loadAnswer(fileId)
				if err != nil {
					return nil, err
				}
				answers = append(answers, answer)
			}
		}
	}

	return answers, nil
}

/*---------- FindAnswer ----------*/
func (db *DB) FindAnswer(fileId string) (*Answer, error) {
	db.schemaLock.RLock()
	defer db.schemaLock.RUnlock()

	return db.loadAnswer(fileId)
}

/*---------- SaveAnswer ----------*/
func (db *DB) SaveAnswer(a *Answer) (string, error) {
	db.schemaLock.Lock()
	defer db.schemaLock.Unlock()

	if a.FileId == "" {
		fileId, err := db.nextAvailableAnswerFileId()
		if err != nil {
			return "", err
		}

		a.FileId = fileId
		a.CreatedById = "1"
		a.CreatedAt = time.Now()
	} else {
		// Is fileid valid?
		_, err := strconv.Atoi(a.FileId)
		if err != nil {
			return "", err
		}
	}

	a.UpdatedById = "1"
	a.UpdatedAt = time.Now()

	marshalledAnswer, err := json.Marshal(a)

	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%v/%v.json", db.answersPath, a.FileId)

	err = ioutil.WriteFile(filename, marshalledAnswer, 0600)
	if err != nil {
		return "", err
	}

	err = db.initAnswerTagsIndex()
	if err != nil {
		return "", err
	}

	return a.FileId, nil
}

/*---------- DeleteAnswer ----------*/
func (db *DB) DeleteAnswer(fileId string) error {
	_, err := strconv.Atoi(fileId)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%v/%v.json", db.answersPath, fileId)

	db.schemaLock.Lock()
	defer db.schemaLock.Unlock()

	err = os.Remove(filename)
	if err != nil {
		return err
	}

	err = db.initAnswerTagsIndex()
	if err != nil {
		return err
	}

	return nil
}

//=============================================================================
// Private Methods
//=============================================================================

/*---------- initAnswerTagsIndex ----------*/
func (db *DB) initAnswerTagsIndex() error {
	for k := range db.answerTagsIndex {
		delete(db.answerTagsIndex, k)
	}

	for _, f := range db.fileIdsInAnswersDataDir() {
		a, err := db.loadAnswer(f)
		if err != nil {
			return err
		}

		tags := strings.Split(a.Tags, " ")

		db.answerTagsIndex[a.FileId] = tags
	}

	return nil
}

/*---------- fileIdsInAnswersDataDir ----------*/
func (db *DB) fileIdsInAnswersDataDir() []string {
	var ids []string

	files, _ := ioutil.ReadDir(db.answersPath)
	for _, file := range files {
		if !file.IsDir() {
			if path.Ext(file.Name()) == ".json" {
				ids = append(ids, file.Name()[:len(file.Name())-5])
			}
		}
	}

	return ids
}

/*---------- nextAvailableAnswerFileId ----------*/
func (db *DB) nextAvailableAnswerFileId() (string, error) {
	var fileIds []int
	var nextFileId string

	for _, f := range db.fileIdsInAnswersDataDir() {
		fileId, err := strconv.Atoi(f)
		if err != nil {
			return "", err
		}

		fileIds = append(fileIds, fileId)
	}

	if len(fileIds) == 0 {
		nextFileId = "1"
	} else {
		sort.Ints(fileIds)
		lastFileId := fileIds[len(fileIds)-1]

		nextFileId = strconv.Itoa(lastFileId + 1)
	}

	return nextFileId, nil
}

/*---------- loadAnswer ----------*/
func (db *DB) loadAnswer(fileId string) (*Answer, error) {
	var answer *Answer

	filename := fmt.Sprintf("%v/%v.json", db.answersPath, fileId)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &answer)

	return answer, nil
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
