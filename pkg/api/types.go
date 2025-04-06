package api

type TaskRequest struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TaskResponse struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TaskAddResponse struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
