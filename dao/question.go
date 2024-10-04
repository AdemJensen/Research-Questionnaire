package dao

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Question struct {
	ID              int64            `gorm:"column:id;primaryKey"`
	QuestionnaireID int64            `gorm:"column:questionnaire_id;index"`
	Index           int              `gorm:"column:index;index"` // starts from 0
	Content         *QuestionContent `gorm:"column:content;type:json"`
	Answer          *QuestionAnswer  `gorm:"column:answer;type:json"`
	OperatedBy      string           `gorm:"column:operated_by;type:text"`
	CreatedAt       time.Time        `gorm:"column:ctime"` // 创建时间
	UpdatedAt       time.Time        `gorm:"column:mtime"` // 更新时间
}

func (q *Question) TableName() string {
	return "questions"
}

type QuestionContent struct {
	ImageUrl        string    `json:"image_url"`
	Reviews         []*Review `json:"reviews"`
	ReviewValence   string    `json:"review_valence"`   // negative, neutral, positive
	ReviewDepth     string    `json:"review_depth"`     // long, short
	ReviewerProfile string    `json:"reviewer_profile"` // default, specified
}

func (qc *QuestionContent) Value() (driver.Value, error) {
	// Marshal the QuestionContent to JSON
	if qc == nil {
		return nil, nil
	}
	return json.Marshal(qc)
}

func (qc *QuestionContent) Scan(value interface{}) error {
	// Scan the JSON value into the QuestionContent struct
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &qc)
}

type Review struct {
	AvatarUri string `json:"avatar_uri"`
	Nickname  string `json:"nickname"`
	Text      string `json:"text"`
}

type QuestionAnswer struct {
	PurchaseIntention int `json:"purchase_intention"` // 0 = not answered, 1~6 = 1~6 stars
}

func (qa *QuestionAnswer) Value() (driver.Value, error) {
	// Marshal the QuestionContent to JSON
	if qa == nil {
		return nil, nil
	}
	return json.Marshal(qa)
}

func (qa *QuestionAnswer) Scan(value interface{}) error {
	// Scan the JSON value into the QuestionContent struct
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &qa)
}

func GetQuestionsByQuestionnaireID(tx *gorm.DB, questionnaireID int64) ([]*Question, error) {
	var questions []*Question
	if err := tx.Where("questionnaire_id = ?", questionnaireID).Order("`index`").Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

func CreateQuestion(tx *gorm.DB, question *Question) error {
	return tx.Create(question).Error
}

func AnswerQuestion(tx *gorm.DB, questionnaireID int64, index int, answer *QuestionAnswer, operatedBy string) error {
	return tx.Model(&Question{}).
		Where("questionnaire_id = ? AND `index` = ?", questionnaireID, index).
		Updates(map[string]interface{}{
			"answer":      answer,
			"operated_by": operatedBy,
		}).Error
}
