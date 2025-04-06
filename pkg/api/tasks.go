package api

import (
	"net/http"
	"strconv"

	"github.com/AleksandrTveritin/go_final_project/pkg/db"
)

type TasksResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}

// TasksHandler - обработчик HTTP-запросов для получения списка задач.
// Он выполняет следующие шаги:
// 1. Проверка метода запроса. Если метод не GET, то возвращается ошибка.
// 2. Получение списка задач из базы данных с использованием функции db.GetTasks. Если получение списка задач не удалось, то возвращается ошибка.
// 3. Преобразование списка задач в формат, подходящий для ответа. Для каждой задачи создается структура TaskResponse, которая содержит информацию о задаче.
// 4. Возвращение списка задач в формате JSON.
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
			Error: "Method not allowed",
		})
		return
	}

	tasks, err := db.GetTasks(DefaultTaskLimit, "")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get tasks",
		})
		return
	}

	respTasks := make([]TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		respTasks = append(respTasks, TaskResponse{
			ID:      strconv.FormatInt(task.ID, Base10),
			Date:    task.Date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		})
	}

	writeJSON(w, http.StatusOK, TasksResponse{
		Tasks: respTasks,
	})
}
