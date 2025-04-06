package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AleksandrTveritin/go_final_project/pkg/api"
	"github.com/AleksandrTveritin/go_final_project/pkg/config"
)

// Run - функция для запуска сервера.
// Она выполняет следующие шаги:
// 1. Инициализация API с использованием функции api.Init.
// 2. Настройка обработчика файлов с использованием функции http.FileServer.
// 3. Получение порта из конфигурации приложения.
// 4. Запуск сервера с использованием функции http.ListenAndServe. Если запуск сервера не удался, то возвращается ошибка.
// 5. Возвращение ошибки, если она возникла при запуске сервера.
func Run() error {
	api.Init()

	http.Handle("/", http.FileServer(http.Dir("web")))

	port := config.AppConfig.Port
	log.Printf("Сервер запущен на порту: %d", port)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
