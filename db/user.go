package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
)

type User struct {
	FileId   string `json:"fileid"`
	Name     string `json:"name"`
	login    string `json:"login"`
	Password string `json:"password"`
}

/*****************************************************************************/
// Public User Methods
/*****************************************************************************/

/*---------- FindUsers ----------*/
func (db *DB) FindUsers() ([]*User, error) {
	var users []*User

	db.usersLock.RLock()
	defer db.usersLock.RUnlock()

	for _, fileId := range db.fileIdsInUsersDataDir() {
		user, err := db.loadUser(fileId)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

/*---------- FindUser ----------*/
func (db *DB) FindUser(fileId string) (*User, error) {
	db.usersLock.RLock()
	defer db.usersLock.RUnlock()

	return db.loadUser(fileId)
}

/*---------- SaveUser ----------*/
func (db *DB) SaveUser(u *User) (string, error) {
	db.usersLock.Lock()
	defer db.usersLock.Unlock()

	if u.FileId == "" {
		fileId, err := db.nextAvailableUserFileId()
		if err != nil {
			return "", err
		}

		u.FileId = fileId
	} else {
		// Is fileid valid?
		_, err := strconv.Atoi(u.FileId)
		if err != nil {
			return "", err
		}
	}

	marshalledAnswer, err := json.Marshal(u)

	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%v/%v.json", db.usersPath, u.FileId)

	err = ioutil.WriteFile(filename, marshalledAnswer, 0600)
	if err != nil {
		return "", err
	}

	return u.FileId, nil
}

/*---------- DeleteUser ----------*/
func (db *DB) DeleteUser(fileId string) error {
	_, err := strconv.Atoi(fileId)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%v/%v.json", db.usersPath, fileId)

	db.usersLock.Lock()
	defer db.usersLock.Unlock()

	err = os.Remove(filename)
	if err != nil {
		return err
	}

	return nil
}

//=============================================================================
// Private User Methods
//=============================================================================

/*---------- fileIdsInUsersDataDir ----------*/
func (db *DB) fileIdsInUsersDataDir() []string {
	var ids []string

	files, _ := ioutil.ReadDir(db.usersPath)
	for _, file := range files {
		if !file.IsDir() {
			if path.Ext(file.Name()) == ".json" {
				ids = append(ids, file.Name()[:len(file.Name())-5])
			}
		}
	}

	return ids
}

/*---------- nextAvailableUserFileId ----------*/
func (db *DB) nextAvailableUserFileId() (string, error) {
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

/*---------- loadUser ----------*/
func (db *DB) loadUser(fileId string) (*User, error) {
	var user *User

	filename := fmt.Sprintf("%v/%v.json", db.usersPath, fileId)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &user)

	return user, nil
}
