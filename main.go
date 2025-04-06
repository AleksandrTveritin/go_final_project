package main

import (
	"log"

	"github.com/AleksandrTveritin/go_final_project/pkg/config"
	"github.com/AleksandrTveritin/go_final_project/pkg/db"
	"github.com/AleksandrTveritin/go_final_project/pkg/server"
)

func main() {
	// 1. Загрузка конфигурации
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	// 2. Инициализация БД
	if err := db.Init(config.AppConfig.DBFile); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 3. Запуск сервера
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
