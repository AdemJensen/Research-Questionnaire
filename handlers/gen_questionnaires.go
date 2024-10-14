package handlers

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"researchQuestionnaire/config"
	"researchQuestionnaire/dao"
	"researchQuestionnaire/utils"
	"strings"

	"gorm.io/gorm"
)

var (
	ReviewerProfileTypes = []string{"default", "specified"}
	ReviewValenceTypes   = []string{"negative", "neutral", "positive"}
	ReviewDepthTypes     = []string{"long", "short"}
)

func GenQuestionnaires(urlRoot string) ([]string, error) {
	// load all sources
	var (
		reviewerAvatarLists   = make(map[string][]string)
		reviewerNicknameLists = make(map[string][]string)
		reviewTextLists       = make(map[string]map[string][]string) // valence -> depth -> texts
	)
	// load reviewer avatars
	avatarDirRootPath := "static/questionnaire_contents/avatars"
	for _, profileType := range ReviewerProfileTypes {
		avatarDirPath := fmt.Sprintf("%s/%s", avatarDirRootPath, profileType)
		avatarFilenames, err := ReadFileNamesFromDir(avatarDirPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read avatar dir: %w", err)
		}
		for i := range avatarFilenames {
			avatarFilenames[i] = fmt.Sprintf("%s/%s", avatarDirPath, avatarFilenames[i])
		}
		reviewerAvatarLists[profileType] = avatarFilenames
	}
	log.Printf("loaded reviewer avatars: %v", utils.Json(reviewerAvatarLists))

	// load reviewer nicknames
	nicknameFileRootPath := "static/questionnaire_contents/nicknames"
	for _, profileType := range ReviewerProfileTypes {
		nicknameFilePath := fmt.Sprintf("%s/%s.txt", nicknameFileRootPath, profileType)
		nicknameFilenames, err := ReadListFromFile(nicknameFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read nickname file: %w", err)
		}
		reviewerNicknameLists[profileType] = nicknameFilenames
	}
	log.Printf("loaded reviewer nicknames: %v", utils.Json(reviewerNicknameLists))

	// load review texts
	reviewTextRootPath := "static/questionnaire_contents/reviews"
	for _, valenceType := range ReviewValenceTypes {
		reviewTextLists[valenceType] = make(map[string][]string)
		for _, depthType := range ReviewDepthTypes {
			reviewTextFilePath := fmt.Sprintf("%s/%s_%s.txt", reviewTextRootPath, valenceType, depthType)
			reviewTextFilenames, err := ReadListFromFile(reviewTextFilePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read review text file: %w", err)
			}
			reviewTextLists[valenceType][depthType] = reviewTextFilenames
		}
	}
	log.Printf("loaded review texts: %v", utils.Json(reviewTextLists))

	// generate questionnaires
	var (
		questionnaires []*dao.Questionnaire
		urls           []string
	)
	for i := 0; i < config.Conf.NQuestionnaires; i++ {
		var questions []*dao.Question
		for len(questions) < config.Conf.NQuestions {
			// every iteration, generate all possible questions, and shuffle or add more, to make types of questions more diverse
			var readyQuestions []*dao.Question
			for _, reviewerProfileType := range ReviewerProfileTypes {
				for _, reviewValenceType := range ReviewValenceTypes {
					for _, reviewDepthType := range ReviewDepthTypes {
						question := &dao.Question{
							QuestionnaireID: 0,
							Index:           0,
							Content: &dao.QuestionContent{
								ImageUrl:        "/static/product.png",
								Reviews:         nil,
								ReviewValence:   reviewValenceType,
								ReviewDepth:     reviewDepthType,
								ReviewerProfile: reviewerProfileType,
							},
							Answer:     &dao.QuestionAnswer{},
							OperatedBy: "gen_questionnaires",
						}
						// generate question by type
						avatarUris := utils.RandomPick(reviewerAvatarLists[reviewerProfileType], config.Conf.NReviews)
						nicknames := utils.RandomPick(reviewerNicknameLists[reviewerProfileType], config.Conf.NReviews)
						texts := utils.RandomPick(reviewTextLists[reviewValenceType][reviewDepthType], config.Conf.NReviews)
						for j := 0; j < config.Conf.NReviews; j++ {
							question.Content.Reviews = append(question.Content.Reviews, &dao.Review{
								AvatarUri: avatarUris[j],
								Nickname:  nicknames[j],
								Text:      texts[j],
							})
						}
						readyQuestions = append(readyQuestions, question)
					}
				}
			}
			rand.Shuffle(len(readyQuestions), func(i, j int) {
				readyQuestions[i], readyQuestions[j] = readyQuestions[j], readyQuestions[i]
			})
			questions = append(questions, readyQuestions...)
		}
		if len(questions) > config.Conf.NQuestions {
			questions = questions[:config.Conf.NQuestions]
		}

		err := dao.Db.Transaction(func(tx *gorm.DB) error {
			newQuestionnaire := &dao.Questionnaire{}
			err := dao.CreateQuestionnaire(tx, newQuestionnaire)
			if err != nil {
				return fmt.Errorf("failed to create questionnaire: %w", err)
			}
			questionnaires = append(questionnaires, newQuestionnaire)
			for j, question := range questions {
				question.QuestionnaireID = newQuestionnaire.ID
				question.Index = j
				err := dao.CreateQuestion(tx, question)
				if err != nil {
					return fmt.Errorf("failed to create question %d: %w", j, err)
				}
			}
			urls = append(urls, fmt.Sprintf("%s/question?p=%s", urlRoot, GenerateScenarioPageParam(newQuestionnaire.ID)))
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to execute DB Transaction: %w", err)
		}
	}
	return urls, nil
}

func WriteListToFile(list []string, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	for _, item := range list {
		_, err := file.WriteString("###\n" + item + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	return nil
}

func ReadListFromFile(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	var list []string
	var item string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "###" {
			item = strings.TrimSpace(item)
			if item != "" {
				list = append(list, item)
				item = ""
			}
		} else {
			item += line + "\n"
		}
	}
	list = append(list, item)
	return list, nil
}

func ReadFileNamesFromDir(dirPath string) ([]string, error) {
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %w", err)
	}

	var filenames []string
	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}
		filenames = append(filenames, entry.Name())
	}

	return filenames, nil
}
