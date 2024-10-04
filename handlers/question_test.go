package handlers

import "testing"

func TestGenerateQuestionPageParam(t *testing.T) {
	baseUrl := "http://localhost:8080"
	var questionnaireID int64 = 1
	questionID := 1
	params := GenerateQuestionPageParam(questionnaireID, questionID)
	t.Log(baseUrl + "/question?p=" + params)
}

func TestGenerateScenarioPageParam(t *testing.T) {
	baseUrl := "http://localhost:8080"
	var questionnaireID int64 = 1
	params := GenerateScenarioPageParam(questionnaireID)
	t.Log(baseUrl + "/question?p=" + params)
}
