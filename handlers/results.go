package handlers

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"researchQuestionnaire/config"
	"researchQuestionnaire/dao"
	"slices"
	"strings"
	"time"
)

//go:embed templates/results.html
var resultsTmplString string

var resultsTmpl = template.Must(template.New("results").Funcs(funcMap).Parse(resultsTmplString))

func ResultsHandler(w http.ResponseWriter, r *http.Request) {
	pwd := r.URL.Query().Get("pwd")
	if pwd != "123456" {
		http.NotFound(w, r)
		return
	}
	doSort := r.URL.Query().Get("sort") == "1"
	// 从数据库中获取所有问卷答案
	var answers []*QuestionnaireAnswer
	allQuestionnaires, err := dao.GetAllQuestionnaires()
	for _, questionnaire := range allQuestionnaires {
		ans := &QuestionnaireAnswer{
			ID:           questionnaire.ID,
			Url:          "/question?p=" + GenerateScenarioPageParam(questionnaire.ID),
			UrlDisplay:   config.Conf.Host + "/question?p=" + GenerateScenarioPageParam(questionnaire.ID),
			CreateTime:   questionnaire.CreatedAt.Format(time.DateTime),
			LastFillTime: "",
			Comment:      questionnaire.Comment,
			Answers:      nil,
			TypeStr:      nil,
		}
		// 获取问卷的所有问题
		questions, err := dao.GetQuestionsByQuestionnaireID(dao.Db, questionnaire.ID)
		if err != nil {
			log.Printf("Failed to get questions for questionnaire %d: %v\n", questionnaire.ID, err)
			continue
		}
		// sort questions by type
		if doSort {
			slices.SortFunc(questions, func(i, j *dao.Question) int {
				if res := strings.Compare(i.Content.ReviewerProfile, j.Content.ReviewerProfile); res != 0 {
					return res
				}
				if res := strings.Compare(i.Content.ReviewValence, j.Content.ReviewValence); res != 0 {
					return res
				}
				return strings.Compare(i.Content.ReviewDepth, j.Content.ReviewDepth)
			})
		}
		// 为每个问题统计答案
		var (
			lastFillTime = questionnaire.CreatedAt
			filled       = false
		)
		for _, question := range questions {
			if question.Answer == nil {
				question.Answer = &dao.QuestionAnswer{}
			}
			ans.Answers = append(ans.Answers, question.Answer.PurchaseIntention)
			ans.TypeStr = append(ans.TypeStr, fmt.Sprintf("%s, %s, %s", question.Content.ReviewerProfile, question.Content.ReviewValence, question.Content.ReviewDepth))
			if question.Answer.PurchaseIntention > 0 && question.UpdatedAt.After(lastFillTime) {
				lastFillTime = question.UpdatedAt
				filled = true
			}
		}
		if filled {
			ans.LastFillTime = lastFillTime.Format(time.DateTime)
		}
		answers = append(answers, ans)
	}

	err = resultsTmpl.Execute(w, &ResultsContext{
		DisplayAll: r.URL.Query().Get("show_all") == "1",
		Sort:       doSort,
		Content:    answers,
	})
	if err != nil {
		log.Fatal("Failed to execute template:", err)
	}
}

type ResultsContext struct {
	DisplayAll bool
	Sort       bool
	Content    []*QuestionnaireAnswer
}

type QuestionnaireAnswer struct {
	ID           int64
	Url          string
	UrlDisplay   string
	CreateTime   string
	LastFillTime string
	Comment      string
	Answers      []int
	TypeStr      []string
}
