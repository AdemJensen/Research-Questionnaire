package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type QuestionUrlBase struct {
	ProductImage     string
	FirstQuestion    string
	NextQuestion     string
	PreviousQuestion string
	Finish           string
	Submit           string
	Questions        []string
}

func NewQuestionUrlBase(totalQuestionNum int, params *Parameters) *QuestionUrlBase {
	questionUrlBase := "/question?p="
	res := &QuestionUrlBase{
		ProductImage:     "/static/product.png",
		FirstQuestion:    questionUrlBase + GenerateQuestionPageParam(params.QuestionnaireID, 0),
		NextQuestion:     "",
		PreviousQuestion: "",
		Finish:           "",
		Submit:           "/submit",
		Questions:        make([]string, totalQuestionNum),
	}
	if params.QuestionIndex > 0 {
		res.PreviousQuestion = questionUrlBase + GenerateQuestionPageParam(params.QuestionnaireID, params.QuestionIndex-1)
	}
	if params.QuestionIndex < totalQuestionNum-1 {
		res.NextQuestion = questionUrlBase + GenerateQuestionPageParam(params.QuestionnaireID, params.QuestionIndex+1)
	} else {
		res.Finish = questionUrlBase + GenerateFinishPageParam(params.QuestionnaireID)
	}
	for i := range res.Questions {
		res.Questions[i] = questionUrlBase + GenerateQuestionPageParam(params.QuestionnaireID, i)
	}

	return res
}

type Parameters struct {
	PageType        string `json:"pt"`  // scene, finish, question
	QuestionnaireID int64  `json:"id"`  // 问卷ID
	QuestionIndex   int    `json:"idx"` // 问题索引
}

func GenerateScenarioPageParam(questionnaireID int64) string {
	params := Parameters{
		PageType:        "scene",
		QuestionnaireID: questionnaireID,
	}
	v, _ := json.Marshal(params)
	return base64.URLEncoding.EncodeToString(v)
}

func GenerateFinishPageParam(questionnaireID int64) string {
	params := Parameters{
		PageType:        "finish",
		QuestionnaireID: questionnaireID,
	}
	v, _ := json.Marshal(params)
	return base64.URLEncoding.EncodeToString(v)
}

func GenerateQuestionPageParam(questionnaireID int64, questionIndex int) string {
	params := Parameters{
		PageType:        "question",
		QuestionnaireID: questionnaireID,
		QuestionIndex:   questionIndex,
	}
	v, _ := json.Marshal(params)
	return base64.URLEncoding.EncodeToString(v)
}

func GetQuestionParamsFromQuery(r *http.Request) (*Parameters, error) {
	return ParseQuestionParams(r.URL.Query().Get("p"))
}

func ParseQuestionParams(str string) (*Parameters, error) {
	params, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}
	var p Parameters
	err = json.Unmarshal([]byte(strings.TrimSpace(string(params))), &p)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return &p, nil
}
