package api

import (
	"encoding/json"
	"net/http"
)

// Init - функция инициализации обработчиков HTTP-запросов.
// Она регистрирует обработчики для следующих эндпоинтов:
// - "/api/nextdate" - обработчик NextDateHandler
// - "/api/task" - обработчик taskHandler
// - "/api/tasks" - обработчик TasksHandler
// - "/api/task/done" - обработчик taskDoneHandler
func Init() {
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", TasksHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler)
}

// taskHandler - обработчик HTTP-запросов для работы с задачами.
// Он выполняет следующие шаги:
// 1. Проверка метода запроса.
// 2. В зависимости от метода запроса вызывается соответствующий обработчик:
//   - GET - getTaskHandler
//   - PUT - updateTaskHandler
//   - POST - AddTaskHandler
//   - DELETE - deleteTaskHandler
//
// 3. Если метод запроса не поддерживается, то возвращается ошибка.
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodPost:
		AddTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
			Error: "Method not allowed",
		})
	}
}

// writeJSON - функция для записи данных в формате JSON в HTTP-ответ.
// Она выполняет следующие шаги:
// 1. Установка заголовка Content-Type в "application/json; charset=UTF-8".
// 2. Установка статуса ответа.
// 3. Кодирование данных в формат JSON и запись в тело ответа.
// 4. Если возникает ошибка при кодировании данных, то возвращается ошибка.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
