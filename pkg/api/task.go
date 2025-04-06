package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AleksandrTveritin/go_final_project/pkg/db"
)

// getTaskHandler - обработчик HTTP-запросов для получения задачи по ее идентификатору.
// Он выполняет следующие шаги:
// 1. Получение параметра id из URL запроса.
// 2. Проверка наличия параметра id. Если параметр отсутствует, то возвращается ошибка.
// 3. Преобразование параметра id в целочисленное значение. Если преобразование не удалось, то возвращается ошибка.
// 4. Получение задачи из базы данных с использованием функции db.GetTask. Если задача не найдена, то возвращается ошибка.
// 5. Возвращение задачи в формате JSON.
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "Task ID is required",
		})
		return
	}

	id, err := strconv.ParseInt(idStr, Base10, BitSize64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "Invalid task ID",
		})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, ErrorResponse{
			Error: "Task not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, Task{
		ID:      idStr,
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	})
}

// updateTaskHandler - обработчик HTTP-запросов для обновления задачи.
// Он выполняет следующие шаги:
// 1. Декодирование тела запроса в структуру TaskRequest. Если декодирование не удалось, то возвращается ошибка.
// 2. Проверка обязательных полей: ID и Title. Если они отсутствуют, то возвращается ошибка.
// 3. Преобразование поля ID в целочисленное значение. Если преобразование не удалось, то возвращается ошибка.
// 4. Проверка правила повторения. Если оно указано, то проверяется его корректность с использованием функции NextDate. Если правило некорректно, то возвращается ошибка.
// 5. Обработка даты. Если дата не указана, то используется текущая дата. Если дата указана как "today", то также используется текущая дата. В противном случае, проверяется корректность формата даты. Если формат некорректен, то возвращается ошибка.
// 6. Обновление задачи в базе данных с использованием функции db.UpdateTask. Если задача не найдена, то возвращается ошибка.
// 7. Возвращение ответа в формате JSON.
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req Task
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "Invalid JSON format",
		})
		return
	}

	if req.ID == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "Task ID is required",
		})
		return
	}

	id, err := strconv.ParseInt(req.ID, Base10, BitSize64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "Invalid task ID",
		})
		return
	}

	if req.Title == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "Task title is required",
		})
		return
	}

	if req.Repeat != "" {
		if _, err := NextDate(time.Now(), time.Now().Format(DateFormat), req.Repeat); err != nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{
				Error: "Invalid repeat format",
			})
			return
		}
	}

	now := time.Now()
	today := now.Format(DateFormat)

	switch {
	case req.Date == "":
		req.Date = today
	case req.Date == "today":
		req.Date = today
	default:
		if _, err := time.Parse(DateFormat, req.Date); err != nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{
				Error: "Invalid date format, expected YYYYMMDD",
			})
			return
		}
	}

	task := &db.Task{
		ID:      id,
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	if err := db.UpdateTask(task); err != nil {
		writeJSON(w, http.StatusNotFound, ErrorResponse{
			Error: "Task not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{})
}

// taskDoneHandler - обработчик HTTP-запросов для отметки задачи как выполненной.
// Он выполняет следующие шаги:
// 1. Проверка метода запроса. Если метод не POST, то возвращается ошибка.
// 2. Получение параметра id из URL запроса.
// 3. Преобразование параметра id в целочисленное значение. Если преобразование не удалось или id меньше или равно нулю, то возвращается ошибка.
// 4. Получение задачи из базы данных с использованием функции db.GetTask. Если задача не найдена, то возвращается ошибка.
// 5. Если правило повторения не указано, то задача удаляется из базы данных с использованием функции db.DeleteTask. Если удаление не удалось, то возвращается ошибка.
// 6. Если правило повторения указано, то вычисляется следующая дата задачи с использованием функции NextDate. Если вычисление не удалось, то возвращается ошибка.
// 7. Обновление даты задачи в базе данных с использованием функции db.UpdateDate. Если обновление не удалось, то возвращается ошибка.
// 8. Возвращение ответа в формате JSON.
func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:

	default:
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "Only POST method allowed"})
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid task ID"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "Task not found"})
		return
	}

	switch task.Repeat {
	case "":
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete task"})
			return
		}

	default:
		baseDate, err := time.Parse(DateFormat, task.Date)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Invalid task date format"})
			return
		}

		nextDate, err := NextDate(baseDate, task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid repeat rule: " + err.Error()})
			return
		}

		if err := db.UpdateDate(id, nextDate); err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to update task date: " + err.Error()})
			return
		}
	}

	writeJSON(w, http.StatusOK, struct{}{})
}

// deleteTaskHandler - обработчик HTTP-запросов для удаления задачи.
// Он выполняет следующие шаги:
// 1. Проверка метода запроса. Если метод не DELETE, то возвращается ошибка.
// 2. Получение параметра id из URL запроса.
// 3. Проверка наличия параметра id. Если параметр отсутствует, то возвращается ошибка.
// 4. Преобразование параметра id в целочисленное значение. Если преобразование не удалось или id меньше или равно нулю, то возвращается ошибка.
// 5. Удаление задачи из базы данных с использованием функции db.DeleteTask. Если задача не найдена, то возвращается ошибка.
// 6. Возвращение ответа в формате JSON.
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
	default:
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
			Error: "Method not allowed",
		})
		return
	}

	idStr := r.URL.Query().Get("id")
	switch {
	case idStr == "":
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "Task ID is required",
		})
		return
	}

	id, err := strconv.ParseInt(idStr, Base10, BitSize64)
	switch {
	case err != nil, id <= 0:
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "Invalid task ID format",
		})
		return
	}

	switch err := db.DeleteTask(id); {
	case err == nil:
		writeJSON(w, http.StatusOK, struct{}{})
	case strings.Contains(err.Error(), "not found"):
		writeJSON(w, http.StatusNotFound, ErrorResponse{
			Error: err.Error(),
		})
	default:
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: "Could not delete task",
		})
	}
}
