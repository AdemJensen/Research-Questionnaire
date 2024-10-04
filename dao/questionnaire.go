package dao

import (
	"time"

	"gorm.io/gorm"
)

type Questionnaire struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	CreatedAt time.Time `gorm:"column:ctime"`   // 创建时间
	UpdatedAt time.Time `gorm:"column:mtime"`   // 更新时间
	Comment   string    `gorm:"column:comment"` // 问卷备注
}

func (q *Questionnaire) TableName() string {
	return "questionnaires"
}

func GetAllQuestionnaires() ([]*Questionnaire, error) {
	var questionnaire []*Questionnaire
	if err := Db.Model(&Questionnaire{}).Scan(&questionnaire).Error; err != nil {
		return nil, err
	}
	return questionnaire, nil
}

func GetQuestionnaireByID(id int64) (*Questionnaire, error) {
	var questionnaire Questionnaire
	if err := Db.Where("id = ?", id).First(&questionnaire).Error; err != nil {
		return nil, err
	}
	return &questionnaire, nil
}

func CreateQuestionnaire(tx *gorm.DB, questionnaire *Questionnaire) error {
	return tx.Create(questionnaire).Error
}

func UpdateQuestionnaireComment(tx *gorm.DB, id int64, comment string) error {
	return tx.Model(&Questionnaire{}).Where("id = ?", id).Update("comment", comment).Error
}
