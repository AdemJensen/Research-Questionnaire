package nicknames

import (
	"math/rand/v2"
	"researchQuestionnaire/handlers"
	"strconv"
	"testing"
)

func TestGenRandomNicknames(t *testing.T) {
	filepath := "default.txt"
	cnt := 100

	var nicknames []string
	for i := 0; i < cnt; i++ {
		// generate default nicknames
		name := "User" + strconv.Itoa(100000+rand.IntN(1000000))
		nicknames = append(nicknames, name)
	}
	err := handlers.WriteListToFile(nicknames, filepath)
	if err != nil {
		t.Errorf("Error writing nicknames to file: %s", err)
		return
	}
}
