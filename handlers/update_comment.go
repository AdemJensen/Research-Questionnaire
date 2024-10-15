package handlers

import (
	"encoding/json"
	"net/http"
	"researchQuestionnaire/dao"
	"strconv"
)

type UpdateCommentRequest struct {
	ID      string `json:"id"`
	Comment string `json:"comment"`
}

func UpdateCommentHandler(w http.ResponseWriter, r *http.Request) {
	var req UpdateCommentRequest

	// Parse the request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// parse id
	if req.ID == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var id int64
	id, err = strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Logic to update the comment in the database
	// You should replace this with your actual database update code
	err = dao.UpdateQuestionnaireComment(dao.Db, id, req.Comment)
	if err != nil {
		http.Error(w, "Failed to update comment", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
