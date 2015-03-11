package db

import (
	"fmt"
	"os"
	"testing"
)

var testDir = "/tmp/pythia"
var testDB *DB

func TestSave(t *testing.T) {
	setup()

	a := &Answer{Question: "Test Question", Answer: "Test Answer", Tags: "test testing"}

	fileId, err := testDB.SaveAnswer(a)
	if err != nil {
		t.Error("Error saving:", err)
	}

	if fileId != "1" {
		t.Error("Expected '1', got ", fileId)
	}

	teardown()
}

func TestFind(t *testing.T) {
	setup()

	a := &Answer{Question: "Test Question", Answer: "Test Answer", Tags: "test testing"}

	_, err := testDB.SaveAnswer(a)
	if err != nil {
		t.Error("Error saving:", err)
	}

	b, err := testDB.FindAnswer(a.FileId)
	if err != nil {
		fmt.Println("Error saving:", err)
	}

	if b.FileId != "1" {
		t.Error("Expected '1', got ", b.FileId)
	}

	if b.Question != "Test Question" {
		t.Error("Expected 'Test Question', got ", b.Question)
	}

	if b.Answer != "Test Answer" {
		t.Error("Expected 'Test Answer', got ", b.Answer)
	}

	if b.Tags != "test testing" {
		t.Error("Expected 'test testing', got ", b.Tags)
	}

	teardown()
}

// Helper Functions
func setup() {
	var err error

	teardown()

	testDB, err = OpenDB(testDir)
	if err != nil {
		fmt.Println("Database initialization failed:", err)
	}
}

func teardown() {
	if _, err := os.Stat(testDir); err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("Error checking to see if tmp dir exists:", err)
		}
	} else {
		err := os.RemoveAll(testDir)
		if err != nil {
			fmt.Println("Could not delete tmp dir:", err)
		}
	}
}
