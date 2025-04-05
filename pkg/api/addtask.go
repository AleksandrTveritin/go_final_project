package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AleksandrTveritin/go_final_project/pkg/db"
)

const DateFormat = "20060102"

type TaskRequest struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TaskResponse struct {
	ID    int64  `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, TaskResponse{
			Error: "Method not allowed",
		})
		return
	}

	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, TaskResponse{
			Error: "Invalid JSON format",
		})
		return
	}

	if req.Title == "" {
		writeJSON(w, http.StatusBadRequest, TaskResponse{
			Error: "Task title is required",
		})
		return
	}

	now := time.Now()
	today := now.Format(DateFormat)

	// Если дата не указана — установить сегодняшнюю
	if req.Date == "" {
		req.Date = today
	}

	// Проверка формата даты
	reqDate, err := time.Parse(DateFormat, req.Date)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, TaskResponse{
			Error: "Invalid date format, expected YYYYMMDD",
		})
		return
	}

	// Если дата в прошлом — пересчитываем
	if reqDate.Before(now) {
		if req.Repeat == "" {
			req.Date = today
		} else {
			next, err := NextDate(now, req.Date, req.Repeat)
			if err != nil {
				writeJSON(w, http.StatusBadRequest, TaskResponse{
					Error: err.Error(),
				})
				return
			}
			req.Date = next
		}
	}

	// Если дата — сегодня: не пересчитываем, даже если есть repeat
	if reqDate.Equal(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())) {
		req.Date = today
	}

	task := &db.Task{
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	id, err := db.AddTask(task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, TaskResponse{
			Error: "Failed to create task",
		})
		return
	}

	writeJSON(w, http.StatusOK, TaskResponse{
		ID: id,
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
