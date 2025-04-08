package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AleksandrTveritin/go_final_project/pkg/config"
	"github.com/AleksandrTveritin/go_final_project/pkg/db"
)

// AddTaskHandler - обработчик HTTP-запросов для добавления новой задачи.
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, TaskAddResponse{Error: "Method not allowed"})
		return
	}

	var req Task
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, TaskAddResponse{Error: "Invalid JSON format"})
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		writeJSON(w, http.StatusBadRequest, TaskAddResponse{Error: "Task title is required"})
		return
	}

	switch req.Date {
	case "", "today":
		req.Date = time.Now().Format(config.DateFormat)
	default:
		if _, err := time.Parse(config.DateFormat, req.Date); err != nil {
			writeJSON(w, http.StatusBadRequest, TaskAddResponse{Error: "Invalid date format, expected YYYYMMDD"})
			return
		}
	}

	taskDate, _ := time.Parse(config.DateFormat, req.Date)
	now := time.Now()
	taskDate = time.Date(taskDate.Year(), taskDate.Month(), taskDate.Day(), 0, 0, 0, 0, time.UTC)
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	if taskDate.Before(nowDate) {
		if req.Repeat == "" {
			req.Date = now.Format(config.DateFormat)
		} else {
			nextDate, err := NextDate(now, req.Date, req.Repeat)
			if err != nil {
				writeJSON(w, http.StatusBadRequest, TaskAddResponse{Error: err.Error()})
				return
			}
			req.Date = nextDate
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
		writeJSON(w, http.StatusInternalServerError, TaskAddResponse{Error: "Failed to create task"})
		return
	}

	writeJSON(w, http.StatusOK, TaskAddResponse{ID: strconv.FormatInt(id, config.Base10)})
}
