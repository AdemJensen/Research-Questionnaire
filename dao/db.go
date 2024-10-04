package dao

import (
	"fmt"
	"researchQuestionnaire/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func init() {
	dsn := config.Conf.DbConUrl
	dbCon, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}
	Db = dbCon
}

func MigrateDB() {
	var err error
	err = Db.AutoMigrate(&Question{})
	if err != nil {
		panic(fmt.Sprintf("failed to migrate Question: %v", err))
	}
	err = Db.AutoMigrate(&Questionnaire{})
	if err != nil {
		panic(fmt.Sprintf("failed to migrate Questionnaire: %v", err))
	}
}
