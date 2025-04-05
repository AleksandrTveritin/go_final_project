package api

import "net/http"

func Init() {
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", taskHandler)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AddTaskHandler(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, TaskResponse{
			Error: "Method not allowed",
		})
	}
}
