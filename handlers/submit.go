package handlers

import (
	"log"
	"net/http"
	"researchQuestionnaire/dao"
	"strconv"
)

// SubmitHandler 提交问卷的处理函数
func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}
	log.Print(r.Form)

	params, err := ParseQuestionParams(r.Form.Get("questionnaire_params"))
	if err != nil {
		http.Error(w, "Invalid form submission: missing params: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 从表单中获取用户的回答
	purchaseIntention, err := strconv.Atoi(r.Form.Get("purchase_intention"))
	if err != nil {
		http.Error(w, "Invalid form submission: missing purchase_intention", http.StatusBadRequest)
		return
	}

	err = dao.AnswerQuestion(dao.Db, params.QuestionnaireID, params.QuestionIndex, &dao.QuestionAnswer{
		PurchaseIntention: purchaseIntention,
	}, r.RemoteAddr)
	if err != nil {
		http.Error(w, "Failed to record your response", http.StatusInternalServerError)
		return
	}

	log.Printf("A response has been recorded. QuestionnaireID = %d, Question Index = %d, Response: %d\n", params.QuestionnaireID, params.QuestionIndex, purchaseIntention)
	_, _ = w.Write([]byte("OK"))
}
