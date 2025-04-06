package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/AleksandrTveritin/go_final_project/pkg/db"
)

// AddTaskHandler - обработчик HTTP-запросов для добавления новой задачи.
// Он выполняет следующие шаги:
// 1. Проверка метода запроса. Если метод не POST, то возвращается ошибка.
// 2. Декодирование тела запроса в структуру TaskRequest. Если возникает ошибка, то возвращается ошибка.
// 3. Проверка обязательного поля Title. Если поле не заполнено, то возвращается ошибка.
// 4. Обработка даты. Если дата не указана или указана как "today", то используется сегодняшняя дата.
// 5. Проверка формата даты. Если дата имеет неверный формат, то возвращается ошибка.
// 6. Если дата в прошлом, то используется правило повторения для вычисления следующей даты.
// 7. Создание задачи в БД. Если возникает ошибка, то возвращается ошибка.
// 8. Возвращение идентификатора созданной задачи в формате JSON.
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req Task
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, TaskAddResponse{
				Error: "Invalid JSON format",
			})
			return
		}

		if req.Title == "" {
			writeJSON(w, http.StatusBadRequest, TaskAddResponse{
				Error: "Task title is required",
			})
			return
		}

		now := time.Now()
		today := now.Format(DateFormat)

		switch req.Date {
		case "today", "":
			req.Date = today
		}

		date, err := time.Parse(DateFormat, req.Date)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, TaskAddResponse{
				Error: "Invalid date format, expected YYYYMMDD",
			})
			return
		}

		if date.Before(now) {
			switch {
			case req.Repeat == "":
				req.Date = today
			case req.Repeat == "d 1":
				req.Date = today
			default:
				next, err := NextDate(now, req.Date, req.Repeat)
				if err != nil {
					writeJSON(w, http.StatusBadRequest, TaskAddResponse{
						Error: err.Error(),
					})
					return
				}
				req.Date = next
			}
		}

		task := &db.Task{
			Date:    req.Date,
			Title:   req.Title,
			Comment: req.Comment,
			Repeat:  req.Repeat,
		}

		id, err := db.AddTask(task)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, TaskAddResponse{
				Error: "Failed to create task",
			})
			return
		}

		writeJSON(w, http.StatusOK, TaskAddResponse{
			ID: strconv.FormatInt(id, Base10),
		})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, TaskAddResponse{
			Error: "Method not allowed",
		})
	}
}
