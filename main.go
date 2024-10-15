package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"researchQuestionnaire/config"
	"researchQuestionnaire/dao"
	"researchQuestionnaire/handlers"
)

// 渲染问卷页面的模板

func main() {
	flags := os.Args[1:]
	if len(flags) > 0 {
		// some utilities
		switch flags[0] {
		case "migrate_db":
			dao.MigrateDB()
		case "gen_questions":
			urls, err := handlers.GenQuestionnaires(config.Conf.Host)
			if err != nil {
				log.Fatalf("Error when gen questions: %v", err)
			}
			fmt.Println("Generated questionnaires:")
			for _, url := range urls {
				fmt.Println(url)
			}
		default:
			log.Fatalf("Unknown command: %s", flags[0])
		}
		return
	}
	// 静态文件处理
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// 首页路由
	http.HandleFunc("/question", handlers.QuestionHandler)

	http.HandleFunc("/results", handlers.ResultsHandler)
	http.HandleFunc("/update-comment", handlers.UpdateCommentHandler)

	// 处理提交回答的路由
	http.HandleFunc("/submit", handlers.SubmitHandler)

	log.Printf("Server will start at %s\n", config.Conf.Host)
	log.Printf("Results page will start at %s\n", config.Conf.Host+"/results?pwd=123456")
	log.Printf("All questionnaires page will start at %s\n", config.Conf.Host+"/results?pwd=123456&show_all=1")

	// 启动服务器
	err := http.ListenAndServe(config.Conf.Serve, nil)
	if err != nil {
		log.Fatalf("Error when ListenAndServe: %v", err)
		return
	}
}
