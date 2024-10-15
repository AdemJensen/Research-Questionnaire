package handlers

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"researchQuestionnaire/config"
	"researchQuestionnaire/dao"
	"strings"
)

//go:embed templates/question.html
var questionTmplString string

//go:embed templates/scenario.html
var scenarioTmplString string

//go:embed templates/finish.html
var finishTmplString string

var questionTmpl = template.Must(template.New("question").Funcs(funcMap).Parse(questionTmplString))
var scenarioTmpl = template.Must(template.New("scenario").Funcs(funcMap).Parse(scenarioTmplString))
var finishTmpl = template.Must(template.New("finish").Funcs(funcMap).Parse(finishTmplString))

type QuestionnaireRenderContext struct {
	UrlBase                 *QuestionUrlBase
	Question                *dao.Question
	Questionnaire           *dao.Questionnaire
	QuestionAnswerStatus    []int // status of each question in the questionnaire
	IsQuestionnaireAnswered bool  // If any question in the questionnaire has been answered
	IsQuestionnaireDone     bool  // If all question in the questionnaire has been answered
	NotAnsweredQuestions    string
	TotalQuestions          int
	ParamsStr               string
	Options                 []int
	PreviousAnswer          int
	InitialReviewsToShow    int
}

var funcMap = template.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
	"len": func(a []int) int {
		return len(a)
	},
}

func GetQuestionnaireRenderContext(params *Parameters) (*QuestionnaireRenderContext, error) {
	questionnaire, err := dao.GetQuestionnaireByID(params.QuestionnaireID)
	if err != nil {
		return nil, fmt.Errorf("failed to get questionnaire: %w", err)
	}

	questions, err := dao.GetQuestionsByQuestionnaireID(dao.Db, params.QuestionnaireID)
	if err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}

	if params.QuestionIndex < 0 || params.QuestionIndex >= len(questions) {
		return nil, fmt.Errorf("invalid question index: %d", params.QuestionIndex)
	}

	var (
		questionAnswerStatus    = make([]int, len(questions))
		isQuestionnaireAnswered = false
		IsQuestionnaireDone     = true
		prevAnswer              = 0
		notDoneQuestions        []string
	)
	for i, question := range questions {
		if question.Answer == nil {
			question.Answer = &dao.QuestionAnswer{}
		}
		questionAnswerStatus[i] = question.Answer.PurchaseIntention
		isQuestionnaireAnswered = isQuestionnaireAnswered || questionAnswerStatus[i] > 0
		IsQuestionnaireDone = IsQuestionnaireDone && questionAnswerStatus[i] > 0
		if question.Index == params.QuestionIndex {
			prevAnswer = question.Answer.PurchaseIntention
		}
		if question.Answer.PurchaseIntention == 0 {
			notDoneQuestions = append(notDoneQuestions, fmt.Sprintf("%d", question.Index+1))
		}
	}

	return &QuestionnaireRenderContext{
		UrlBase:                 NewQuestionUrlBase(len(questions), params),
		Question:                questions[params.QuestionIndex],
		Questionnaire:           questionnaire,
		QuestionAnswerStatus:    questionAnswerStatus,
		IsQuestionnaireAnswered: isQuestionnaireAnswered,
		IsQuestionnaireDone:     IsQuestionnaireDone,
		NotAnsweredQuestions:    strings.Join(notDoneQuestions, ", "),
		TotalQuestions:          len(questions),
		ParamsStr:               GenerateQuestionPageParam(params.QuestionnaireID, params.QuestionIndex),
		Options:                 []int{1, 2, 3, 4, 5, 6},
		PreviousAnswer:          prevAnswer,
		InitialReviewsToShow:    config.Conf.NInitialReviews,
	}, nil
}

// QuestionHandler Display scenario and question pages
func QuestionHandler(w http.ResponseWriter, r *http.Request) {
	params, err := GetQuestionParamsFromQuery(r)
	if err != nil {
		http.Error(w, "Invalid questionnaire", http.StatusBadRequest)
		return
	}

	renderCtx, err := GetQuestionnaireRenderContext(params)
	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	switch params.PageType {
	case "s":
		err = scenarioTmpl.Execute(w, renderCtx)
	case "f":
		err = finishTmpl.Execute(w, renderCtx)
	case "q":
		err = questionTmpl.Execute(w, renderCtx)
	}
	if err != nil {
		http.Error(w, "Internal server error: Failed when executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
