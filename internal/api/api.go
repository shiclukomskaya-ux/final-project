package api

import "net/http"

func Init() {
	http.HandleFunc("/api/nextdate", NextDayHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", TasksHandler)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	// case http.MethodGet:

	// case http.MethodDelete:

	case http.MethodPost:
		addTaskHandler(w, r)
	}
}
