package handlers

import (
	"os"
	"path"
	"runtime"
	"testing"
)

func TestGenQuestionnaires(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	urls, err := GenQuestionnaires("http://127.0.0.1:8080")
	if err != nil {
		t.Errorf("Error when gen questions: %v", err)
	}
	for _, url := range urls {
		t.Log(url)
	}
}
