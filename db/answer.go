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
	CreatedBy   string    `json:"createdby,omitempty"`
	UpdatedBy   string    `json:"updatedby,omitempty"`
}

/*****************************************************************************/
// Public Answer Methods
/*****************************************************************************/

/*---------- FindAnswers ----------*/
func (db *DB) FindAnswers(searchString string) ([]*Answer, error) {
	var answers []*Answer
	var possibleMatchingFileIdsMap map[string]int

	db.answersLock.RLock()
	defer db.answersLock.RUnlock()

	if len(searchString) != 0 {
		searchTags := strings.Split(searchString, " ")

		// Need a map to hold possible file ids for answers whose tags include at least one of the search tags.
		possibleMatchingFileIdsMap = make(map[string]int)

		// For each one of the search tags...
		for _, tag := range searchTags {
			// If the search tag is in the index...
			if fileIds, ok := db.answerTagsIndex[tag]; ok {
				// Loop through all the file ids that have that tag in the index...
				for _, fileId := range fileIds {
					// If we have already added that file id to the map of possible matching file ids, then just add 1 to the number of
					// occurrences of that file id.
					if numOfOccurences, ok := possibleMatchingFileIdsMap[fileId]; ok {
						possibleMatchingFileIdsMap[fileId] = numOfOccurences + 1
						// Otherwise, add the file id as a new key in the map of possible matching file ids and set the number of occurrences to 1.
					} else {
						possibleMatchingFileIdsMap[fileId] = 1
					}
				}
			}
		}

		// How many search tags were entered?  We will use this number when we loop through all of the possible matches to determine if the
		// possible match has a number of occurrences as the number of search tags.  If it does, that means that that possible match had
		// all of the tags that we are searching for.
		searchTagsLen := len(searchTags)

		// Now, we only want the possible matching file ids that have a number of occurrences equal to the number of search tags.  If the
		// number of occurrences is less, that means that that particular answer did not have all of the search tags in it's tag list.
		for fileId, numOfOccurrences := range possibleMatchingFileIdsMap {
			if numOfOccurrences == searchTagsLen {
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
func (db *DB) SaveAnswer(answer *Answer) (string, error) {
	db.answersLock.Lock()
	defer db.answersLock.Unlock()

	if answer.FileId == "" {
		fileId, err := db.nextAvailableAnswerFileId()
		if err != nil {
			return "", err
		}

		answer.FileId = fileId

	} else {
		// Is fileid valid?
		_, err := strconv.Atoi(answer.FileId)
		if err != nil {
			return "", err
		}

		originalAnswer, err := db.loadAnswer(answer.FileId)
		if err != nil {
			return "", err
		}

		// These fields never change after creation.
		answer.CreatedById = originalAnswer.CreatedById
		answer.CreatedAt = originalAnswer.CreatedAt
	}

	marshalledAnswer, err := json.Marshal(answer)

	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%v/%v.json", db.answersPath, answer.FileId)

	err = ioutil.WriteFile(filename, marshalledAnswer, 0600)
	if err != nil {
		return "", err
	}

	err = db.initAnswerTagsIndex()
	if err != nil {
		return "", err
	}

	db.initAvailableAnswerTags()

	return answer.FileId, nil
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

	db.initAvailableAnswerTags()

	return nil
}

//=============================================================================
// Private Answer Methods
//=============================================================================

/*---------- initAvailableAnswerTags ----------*/
func (db *DB) initAvailableAnswerTags() {
	db.AvailableAnswerTags = make([]string, len(db.answerTagsIndex))

	i := 0
	for k := range db.answerTagsIndex {
		db.AvailableAnswerTags[i] = k
		i += 1
	}

	sort.Strings(db.AvailableAnswerTags)
}

/*---------- initAnswerTagsIndex ----------*/
func (db *DB) initAnswerTagsIndex() error {
	// Delete all the entries in the index.
	for k := range db.answerTagsIndex {
		delete(db.answerTagsIndex, k)
	}

	// For every file in the data dir...
	for _, f := range db.fileIdsInAnswersDataDir() {
		// Load the file into an answer struct.
		a, err := db.loadAnswer(f)
		if err != nil {
			return err
		}

		// For every tag in the answer...
		for _, tag := range strings.Split(a.Tags, " ") {
			// If the tag already exists as a key in the index...
			if fileIds, ok := db.answerTagsIndex[tag]; ok {
				// Add the file id to the list of ids for that tag, if it is not already in the list.
				if !stringInSlice(a.FileId, fileIds) {
					db.answerTagsIndex[tag] = append(fileIds, a.FileId)
				}
			} else {
				// Otherwise, add the tag with associated new file id to the index.
				db.answerTagsIndex[tag] = []string{a.FileId}
			}
		}
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
