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
	"time"
)

type Answer struct {
	FileId      string    `json:"fileid"`
	Question    string    `json:"question"`
	Answer      string    `json:"answer"`
	CreatedById string    `json:"createdbyid"`
	UpdatedById string    `json:"updatedbyid"`
	CreatedAt   time.Time `json:"createdat"`
	UpdatedAt   time.Time `json:"updatedat"`
	Tags        string    `json:"tags"`
	CreatedBy   string
	UpdatedBy   string
}

/*****************************************************************************/
// Public Answer Methods
/*****************************************************************************/

/*---------- FindAnswers ----------*/
func (db *DB) FindAnswers(searchString string) ([]*Answer, error) {
	var answers []*Answer
	var searchTags []string

	db.answersLock.RLock()
	defer db.answersLock.RUnlock()

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
	db.answersLock.RLock()
	defer db.answersLock.RUnlock()

	return db.loadAnswer(fileId)
}

/*---------- SaveAnswer ----------*/
func (db *DB) SaveAnswer(a *Answer) (string, error) {
	db.answersLock.Lock()
	defer db.answersLock.Unlock()

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

	db.answersLock.Lock()
	defer db.answersLock.Unlock()

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
// Private Answer Methods
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

	user, err := db.FindUser(answer.CreatedById)
	if err != nil {
		return nil, err
	}

	answer.CreatedBy = user.Name

	user, err = db.FindUser(answer.UpdatedById)
	if err != nil {
		return nil, err
	}

	answer.UpdatedBy = user.Name

	return answer, nil
}
