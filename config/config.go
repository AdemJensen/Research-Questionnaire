package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	DbConUrl        string `json:"db_con_url"`
	Host            string `json:"host"`
	Serve           string `json:"serve"`
	NQuestionnaires int    `json:"n_questionnaires"`
	NQuestions      int    `json:"n_questions"`
	NReviews        int    `json:"n_reviews"`
	NInitialReviews int    `json:"n_initial_reviews"` // We added a "show more" button to avoid showing all reviews at once
	NOptions        int    `json:"n_options"`
}

var Conf *Config

func init() {
	confPath := os.Getenv("CONF_PATH")
	if confPath == "" {
		confPath = "config.json"
	}
	conf, err := ReadConfig(confPath)
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %w", err))
	}
	Conf = conf
}

func ReadConfig(filepath string) (*Config, error) {
	// read json content from file
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// 读取文件内容
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 反序列化 JSON 内容
	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &config, nil
}
